# TODO

-   Delete sessions on self lock, user update or admin lock
-   Require at least 2 login alert messages for n messengers to have been successfully sent before authorising download
-   -   n = max(ceil(configured_messengers / 2), 1)
-   -   If 1 is configured, can only require one. Configuring 2 allows one to fail so it's a bit more resilient. 3 still means only 1 can fail, so you get a good balance. And then after 4, neither way is likely to be an issue
-   Use hash-wasm on the frontend when backend rate limits the hashing. It's single threaded but multithreaded WASM seems to be patchy at the moment, even in languages with good support like Rust
-   -   Switch to SvelteKit at the same time
-   -   Backend should limit number of concurrent hash requests to avoid using too much RAM
-   Add limits on self-locking so a hacker can't lock you out forever
-   -   Attempting to get an authorisation code when locked should send the unlock date
-   -   Admins should be able to reset it so if there's an unauthorised login, the user can block with a self lock, the admin can reset them and then they can block again without waiting
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
-   Admin endpoints for troubleshooting:
-   -   Dump database as sqlite file
-   -   Cancel failed job
-   -   Retry failed job
-   -   Update job body
-   Send regular messages to users and the admin
-   -   Should have a clear message if nothing has happened, otherwise it displays totals for each type of message (e.g failed login) and all of the logs in chronological order
-   -   Is it worth having general categories in logs (e.g login) like errors do?
-   -   Occasionally have to click a link in it to verify that messenger is still working
-   -   -   Should that link only be there when necessary?
-   Audit use of time.sleep. Prefer time.After in a select so context cancellations can be respected
-   Recover panics in all of the service implementations and trigger a shutdown. They should recover once if it's a service like the database but otherwise remain shut down
-   Prevent timing attacks from revealing if a user exists or not
-   -   Create test with real endpoint, users in the test database and real hashing to see if an attacker could tell more than 80% of the time with 1000 samples. I guess disable the rate limiting though?
-   -   The tests should have a singlethreaded and multithreaded variant to see if an increased server load reveals more information
-   -   Can probably mitigate by waiting until a response time has been reached before sending the response. Maybe it would start at 1 second but it if it's ever exceeded, the new target would be a whole number of seconds. e.g 1.5 seconds of real processing time would result in a 2 second response time.
-   -   -   How does this safely go down again? Going up isn't particularly safe either
-   -   Admin endpoints don't need this security, as long as they fail early if unauthorised
-   When the admin is locked, whether temporarily or permanently, errors should make the server enter some kind of lockdown state? Need to weigh up pros and cons

-   Create servicescommon so things can be split up better?
-   Job engine should support rate limiting for each API by each definition having an optional function to modify the database object.
    There could be a function to increase the due time based on the internal rate limit for the API. Probably not needed though
-   Refactor the logger
-   -   Mostly to improve the self logging
-   Is the benchmark properly thread-safe? Can guessChan be received in multiple places like that? Maybe should send a done signal down nextPasswordChan to the workers?
-   Bump priority of jobs as they get older

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
