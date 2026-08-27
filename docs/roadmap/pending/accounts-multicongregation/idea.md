# Accounts, publisher portal, and multi-congregation tenancy

- **Slug:** accounts-multicongregation
- **Status:** proposed
- **Merge date:**
- **App version:**

## Summary

One install hosts many congregations. Every user has exactly one congregation. Publishers (and optional students) can log in. Privileges on the publisher row stack; secretary and superadmin are separate flags. The secretary enables access (password, invite email, passkey QR). Superadmin has a home congregation plus a platform switcher. Replaces `roles-beyond-admin`. Design: `docs/superpowers/specs/2026-08-27-accounts-multicongregation-design.md`.

## Out of scope

- Public signup or self-service “create a congregation”.
- Membership in several congregations at once.
- Creating a second superadmin from the UI in the first slices.
- Stored shepherding note text (see dropped `shepherding-notes`).
- Special hall roles (audio/video, attendants, accounts, territories).
- Move UI in the first slice (default).
- Full meeting program, assignment workflow, groups schema, pioneer goals UI — those are other ideas; this one ships permissions, empty nav, and 403 rules.
