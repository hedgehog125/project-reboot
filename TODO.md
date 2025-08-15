# TODO

-   Fix issues around failed jobs. Jobs aren't being marked as completed? Job deadlock
-   Jobs aren't being marked as complete? Maybe just 2FA actions?
-   Rename userRow etc to userOb (or dbUser?)
-   Write some tests to demo how [package] error categories get shifted about
-   Job system
-   -   Rework 2FA actions so there aren't definitions, they're just a full job type ID, body and a username
-   -   Add MostlyDatabase option, only one runs simultaneously
-   -   Continue looking through the job list for lower weight jobs if at max? Maybe only within the priority level?
-   -   Should send a signal to shut down the server instead of shutting itself down. Avoids issues with the shutdown channels
-   Only load env files in development
-   Bump priority of jobs as they get older
-   Replace cron system with a simple custom job scheduler
-   -   Log warning with how many scheduled runs were missed for each scheduled job on startup. Probably not worth adding an option to run them multiple times though
-   Store logs to database
-   -   Should have categories (e.g login) and types (e.g failed login)
-   -   Logs associated with your user that have the user facing attribute set to true are sent in the regular email
-   -   Successful login attempts also directly send a message, a job then runs every day to send reminders
-   -   Probably use slog
-   -   Errors should send a message to the admin user with a cooldown so as not to spam. Maybe just say there's an error, please investigate?
-   Require at least 2 login alert messages for n messengers to have been successfully sent before authorising download
-   -   n = max(ceil(configured_messengers / 2), 1)
-   -   If 1 is configured, can only require one. Configuring 2 allows one to fail so it's a bit more resilient. 3 still means only 1 can fail, so you get a good balance. And then after 4, neither way is likely to be an issue
-   Use hash-wasm on the frontend when backend rate limits the hashing. It's single threaded but multithreaded WASM seems to be patchy at the moment, even in languages with good support like Rust
-   -   Switch to SvelteKit at the same time
-   -   Backend should limit number of concurrent hash requests to avoid using too much RAM
-   Make util that allows returning a servercommon.Error from endpoint handlers
-   Use common.ErrWrapperDatabase as base for DB error wrappers
-   Account locking until a specified date for if you know you won't have access to your devices for a while
-   Automatically delete expired 2FA actions
-   Store timestamp in session so it can be double checked to ensure it's still valid by comparing to auth_timestamps_valid_from in users table
-   Email messenger
-   Signal messenger? Using that REST API in a separate container over the internal network only so no security required, hosting should be very cheap if it's serverless
-   Repeat password in sign up form
-   Does gin.ctx.Context include the timeout info from the middleware?
-   Move more logic out of endpoints
-   Rework endpoint system, maybe the endpoint functions could return an Endpoint struct with an array of handlers and some other things? Middleware should be defined there instead of in RegisterEndpoints
-   Review contexts. Possibly want to give them all a timeout, partly to make shutdowns more predictable
-   SMS messenger
-   Split endpoints into admin and normal?
-   Switch to a pure Go SQLite implementation, speed will be fine considering SQLite it already has the single writer system
-   Does log.Fatalf stop the shutdown logic running if the server crashes on startup?
-   Require both admin and users to click a link every 4 weeks (unless already locked) to confirm their contacts are working. If they don't click it, users will automatically lock and have to be unlocked by an admin. If the admin doesn't, all users will automatically lock

-   Is the benchmark properly thread-safe? Can guessChan be received in multiple places like that? Maybe should send a done signal down nextPasswordChan to the workers?

# To watch

-   Timeouts sometimes incorrectly send 500s?

# To research

Can I wake up a sleeping railway app by just having a separate cron service send an HTTP request over the internal network?

Maybe have the server save the time periodically and on shutdown? Then when it starts it runs through the cron jobs it missed? It probably shouldn't run the same jobs multiple times though

# Testing

-   Create mock messenger
-   -   Register it multiple times in place of the actual ones to ensure the contacts are being passed correctly?
-   Continue fixing linting errors once golang ci v2 is working properly in VSCode
-   Race condition fuzzer that spams a bunch of endpoints
-   -   Would be run with the -race flag
-   -   In particular, test that spamming get-authorization-code with the correct password then updating the password invalidates all of the codes generated using the old password
-   Endpoints
-   -   Do they cancel their work if a request times out? Can encryption/decryption run in the background?
