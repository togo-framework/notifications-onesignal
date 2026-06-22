# notifications-onesignal

[OneSignal](https://onesignal.com) push channel for togo
[notifications](https://github.com/togo-framework/notifications). A Notification that
implements `ToPush` is delivered to the recipient's push tokens.

```bash
togo install togo-framework/notifications
togo install togo-framework/notifications-onesignal
```

Env: `ONESIGNAL_APP_ID`, `ONESIGNAL_API_KEY`. Add `"onesignal"` to a notification's `Via`.
