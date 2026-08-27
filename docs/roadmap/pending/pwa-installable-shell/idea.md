# Installable PWA shell

- **Slug:** pwa-installable-shell
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

Make `https://khhub.app` installable on a phone: home-screen icon, standalone display, no browser chrome. The secretary still talks to `api.khhub.app` with the same cookie session. Static assets may be cached; congregation data is never stored in the service worker.

## Out of scope

- Offline use of reports, publishers, or any API payload
- Push notifications
- Play Store / App Store / Capacitor
- Mobile layout pass (see `mobile-layout-pass`)
