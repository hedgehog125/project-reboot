# TODO

-   Store logs to database
-   -   Should have categories (e.g login) and types (e.g failed login)
-   -   Logs associated with your user that have the user facing attribute set to true are sent in the regular email
-   -   Successful login attempts also directly send a message, a job then runs every day to send reminders
-   Email messenger
-   Retry sending messages
-   Notify using the other methods when a message fails to send
-   -   Get auth code endpoint should require at least 1 messenger to succeed
-   Repeat password in sign up form
-   Account locking until a specified date for if you know you won't have access to your devices for a while
-   SMS messenger

# To research

Can I wake up a sleeping railway app by just having a separate cron service send an HTTP request over the internal network?

Maybe have the server save the time periodically and on shutdown? Then when it starts it runs through the cron jobs it missed? It probably shouldn't run the same jobs multiple times though
