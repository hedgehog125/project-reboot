# TODO

-   Custom error interface for messengers so the caller knows which failed?
-   How should messengers be concurrent? Maybe the different types run concurrently but the messages for each are sent sequentially?
-   Rename loginAttempt table to loginSessions and only store successful logins. Expired sessions are deleted since it'll have already been logged in the log table
-   Email messenger
-   Retry sending messages
-   Store logs to database
-   -   Should have categories (e.g login) and types (e.g failed login)
-   -   Logs associated with your user that have the user facing attribute set to true are sent in the regular email
-   -   Successful login attempts also directly send a message, a job then runs every day to send reminders
-   Revamp login attempts table to store more different types of login attempts. Maybe most should be moved to the general log table?
-   Notify using the other methods when a message fails to send
-   Account locking until a specified date for if you know you won't have access to your devices for a while
-   SMS messager

# To research

Can I wake up a sleeping railway app by just having a separate cron service send an HTTP request over the internal network?

Maybe have the server save the time periodically and on shutdown? Then when it starts it runs through the cron jobs it missed? It probably shouldn't run the same jobs multiple times though
