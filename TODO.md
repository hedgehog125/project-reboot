# TODO

-   Store timestamp in session so it can be double checked to ensure it's still valid by comparing to auth_timestamps_valid_from in users table
-   Account locking until a specified date for if you know you won't have access to your devices for a while
-   Store logs to database
-   -   Should have categories (e.g login) and types (e.g failed login)
-   -   Logs associated with your user that have the user facing attribute set to true are sent in the regular email
-   -   Successful login attempts also directly send a message, a job then runs every day to send reminders
-   Email messenger
-   Retry sending messages
-   Notify using the other methods when a message fails to send
-   -   Get auth code endpoint should require at least 1 messenger to succeed
-   -   If all messengers fail, send the message to env.ADMIN_USERNAME
-   Repeat password in sign up form
-   Use transactions
-   SMS messenger
-   Is the benchmark properly thread-safe? Can guessChan be received in multiple places like that? Maybe should send a done signal down nextPasswordChan to the workers?

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
