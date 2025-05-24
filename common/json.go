package common

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrInvalidJSONString = errors.New("couldn't parse JSON string")

// Adapted from https://github.com/vtopc/epoch/blob/master/str_seconds.go
// TODO: Sounds like json.Marshal already supports this? But maybe Gin doesn't?
type ISOTimeString struct {
	time.Time
}

func (isoTime *ISOTimeString) MarshalJSON() ([]byte, error) {
	return json.Marshal(isoTime.Time.Format(time.RFC3339))
}

func (isoTime *ISOTimeString) UnmarshalJSON(data []byte) error {
	var asString string
	err := json.Unmarshal(data, &asString)
	if err != nil {
		return ErrInvalidJSONString
	}

	until, err := time.Parse(time.RFC3339, asString)
	if err != nil {
		return err
	}
	isoTime.Time = until
	return nil
}
