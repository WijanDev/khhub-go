# Login rate limit

- **Slug:** login-rate-limit
- **Status:** uat
- **Merge date:**
- **App version:**

## Summary

Temporary lockout after repeated failed logins (email/username + password). Passkey attempts should not be a cheap oracle either.

## Out of scope

- CAPTCHA or third-party bot services in the first slice.
- Rate-limiting every GET.
