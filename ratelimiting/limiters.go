package ratelimiting

import (
	"fmt"
	"sync"
	"time"

	"github.com/hedgehog125/project-reboot/common"
	"github.com/jonboulle/clockwork"
)

type Limiter struct {
	Clock  clockwork.Clock
	limits map[string]*limit
	mu     sync.RWMutex
}
type limit struct {
	globalCounter counter
	globalMax     int
	userCounters  map[string]counter
	userMax       int
	resetDuration time.Duration
	mu            sync.Mutex
}
type counter struct {
	value     int
	nextReset time.Time
}

func (counter counter) refresh(now time.Time, resetDuration time.Duration) counter {
	if !now.Before(counter.nextReset) {
		counter.value = 0
		counter.nextReset = now.Add(resetDuration)
	}
	return counter
}

func NewLimiter(clock clockwork.Clock) *Limiter {
	return &Limiter{
		Clock:  clock,
		limits: map[string]*limit{},
	}
}
func (limiter *Limiter) Register(eventName string, globalMax int, userMax int, resetDuration time.Duration) {
	if globalMax != -1 && userMax != -1 && globalMax < userMax {
		panic("globalMax cannot be less than userMax")
	}

	limiter.mu.Lock()
	defer limiter.mu.Unlock()
	_, ok := limiter.limits[eventName]
	if ok {
		panic(fmt.Sprintf("event name %v has already been registered", eventName))
	}
	limiter.limits[eventName] = &limit{
		globalCounter: counter{
			value:     0,
			nextReset: limiter.Clock.Now().Add(resetDuration),
		},
		globalMax:     globalMax,
		userMax:       userMax,
		resetDuration: resetDuration,
		userCounters:  map[string]counter{},
	}
}

type Session struct {
	Amount  int
	User    string
	limit   *limit
	limiter *Limiter
}

func (limiter *Limiter) RequestSession(eventName string, amount int, user string) (*Session, *common.Error) {
	limiter.mu.RLock()
	defer limiter.mu.RUnlock()
	limit, ok := limiter.limits[eventName]
	if !ok {
		return nil, ErrWrapperRequestSession.Wrap(ErrInvalidEventName)
	}
	limit.mu.Lock()
	defer limit.mu.Unlock()
	globalCounter := limit.globalCounter.refresh(limiter.Clock.Now(), limit.resetDuration)
	globalCounter.value += amount
	if limit.globalMax != -1 && globalCounter.value > limit.globalMax {
		return nil, ErrWrapperRequestSession.Wrap(ErrGlobalRateLimitExceeded)
	}
	userCounter := limit.userCounters[user].refresh(limiter.Clock.Now(), limit.resetDuration)
	userCounter.value += amount
	if limit.userMax != -1 && userCounter.value > limit.userMax {
		return nil, ErrWrapperRequestSession.Wrap(ErrUserRateLimitExceeded)
	}

	limit.globalCounter = globalCounter
	limit.userCounters[user] = userCounter
	return &Session{
		Amount:  amount,
		User:    user,
		limit:   limit,
		limiter: limiter,
	}, nil
}
func (session *Session) AdjustTo(amount int) *common.Error {
	if amount == session.Amount {
		return nil
	}
	session.limit.mu.Lock()
	defer session.limit.mu.Unlock()

	diff := amount - session.Amount
	globalCounter := session.limit.globalCounter.refresh(session.limiter.Clock.Now(), session.limit.resetDuration)
	globalCounter.value = max(globalCounter.value+diff, 0)
	if session.limit.globalMax != -1 && globalCounter.value > session.limit.globalMax {
		return ErrWrapperAdjustTo.Wrap(ErrGlobalRateLimitExceeded)
	}
	userCounter := session.limit.userCounters[session.User].refresh(session.limiter.Clock.Now(), session.limit.resetDuration)
	userCounter.value = max(userCounter.value+diff, 0)
	if session.limit.userMax != -1 && userCounter.value > session.limit.userMax {
		return ErrWrapperAdjustTo.Wrap(ErrUserRateLimitExceeded)
	}

	session.limit.globalCounter = globalCounter
	session.limit.userCounters[session.User] = userCounter
	session.Amount = amount
	return nil
}
func (session *Session) Cancel() {
	commErr := session.AdjustTo(0)
	if commErr != nil {
		panic(fmt.Sprintf(
			"ratelimiting.Session.Cancel: session.AdjustTo(0) returned an error. this should not happen! error:\n%v",
			commErr.Dump(),
		))
	}
}

// func (session *Session) RequestSubSession(resource string, amount int) (*Session, *common.Error)
