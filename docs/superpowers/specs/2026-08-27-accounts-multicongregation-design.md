# Design: Accounts, publisher portal, and multi-congregation tenancy

**Date:** 2026-08-27  
**Status:** approved — 2026-08-27  
**Product:** khhub (not official JW software; no jw.org; UI in Spanish, docs in English)

## Problem

khhub is a single-admin, single-congregation tool. The secretary is the only login. Publishers cannot see the meeting program, accept assignments, or submit their own report. The product will host **several congregations on one install**, with a platform owner (superadmin) and per-congregation secretaries. Elders are uncomfortable with a developer reading pastoral notes, so **stored shepherding notes are out**.

## Goals

1. Many congregations in one Postgres. No public signup.
2. Every user has **exactly one** assigned congregation. Changing congregation is a **move**, not multi-membership.
3. Logins for publishers (and optional student logins). Privileges on the publisher row **stack**. **Secretary** is a separate flag. **Superadmin** is the platform owner.
4. Auth stays httpOnly cookie `khhub_session`. No JWT.
5. Password (operator-set), email invite (when SMTP exists), and WebAuthn passkeys via a short-lived QR/code.
6. Login: discoverable passkey (“Entrar con este aparato”) and email/username + password or passkey. No publisher name list on the login screen.
7. Any **403** in the SPA sends the user to the **public landing** and shows an **error toast**.
8. No table, field, upload, or route for shepherding **note text**. Shepherding **planning** (attention list) is allowed.

## Non-goals

- Public self-signup or self-service “create a congregation”.
- Membership in several congregations at once (circuit overseer multi-home). Superadmin uses a **platform switcher** instead.
- A second superadmin creatable from the UI in v1.
- Stored pastoral notes (the existing roadmap idea `shepherding-notes` is dropped).
- Special hall roles (audio/video, attendants, accounts, territories). Later, same stacking model.
- SMTP implementation details beyond “invite email is a first-class path once mail is configured”.
- Replacing S-21 or talking to jw.org.
- Building the full meeting program, assignment workflow, field-service groups, pioneer hour goals, month lock, Memorial, public talks, CSV import, login rate limit, audit log, and dark mode **inside this foundation**. Those stay separate ideas and plug into the permission model below.

## Decisions (locked)

| Topic | Choice |
| --- | --- |
| Tenancy | One install, many congregations, no public signup |
| User ↔ congregation | Exactly one. Move = reassign |
| Account model | `users` row linked to `publishers` (`publisher_id` nullable) |
| Privileges | Read from the publisher row; do not copy into `users` |
| Stacking | Additive + secretary flag + superadmin flag |
| Superadmin | Has a home congregation + platform congregation switcher |
| First superadmin | Today’s `ADMIN_EMAIL` user |
| Publisher v1 (no extra flags) | Program, own assignments accept/reject, own monthly report, directory + groups (see fields) |
| Directory / group fields | Given name, family name, household name (no address), group, privileges. No phone, email, reports, or stats |
| Elder / ministerial servant | Publisher portal + all assignments + attention list. No others’ reports. No notes |
| Regular / special pioneer | Publisher portal + **own** hour goal/progress |
| Auxiliary pioneer | No extra zone; hours stay on that month’s report |
| Student | Meeting program only |
| Unbaptized publisher | Same portal as baptized publisher |
| Self-edit | Any logged-in user may fix **their** first and last name |
| Last secretary | Cannot remove the last secretary of a congregation |
| 403 UX | Redirect to public landing + toast; do not revoke the session |

## Architecture

```
Browser (Vite SPA)
  cookie khhub_session
  api() → on 403: toast + navigate to landing
         on 401: navigate to landing (session missing)

API (Gin)
  session → user → congregation_id (required)
  superadmin may set session active_congregation_id
  authorize(user, publisher flags, secretary, superadmin, active congregation)

Postgres
  congregations (many)
  users.congregation_id + users.publisher_id + is_secretary + is_superadmin
  all congregation data scoped by congregation_id
```

Permissions are computed at request time from the linked publisher row plus flags on `users`. There is no parallel RBAC table.

## Data model

### Congregations

Replace the singleton `congregation` row (`id = 1`) with a normal `congregations` table (UUID PK): name, number, midweek day, weekend day, archived flag, timestamps.

Every household, publisher, field-service report, meeting attendance row, and later group/assignment row gets `congregation_id` (NOT NULL, indexed). Queries always filter by the **active** congregation.

### Users

Add to `users`:

- `congregation_id` NOT NULL (home congregation)
- `publisher_id` NULL UNIQUE (1:1 when set). Must belong to the same congregation.
- `username` NULL UNIQUE, optional short login
- `email` stays unique when set; allow users with passkey-only and no remembered email (email may be empty string or NULL — pick NULL and drop the current NOT NULL if needed)
- `password_hash` nullable (passkey-only users)
- `is_secretary` boolean
- `is_superadmin` boolean
- `login_enabled` boolean

Constraints:

- At least one of password or a passkey must exist before `login_enabled` is true, except during invite/QR pending.
- Superadmin still has a home `congregation_id`.
- Application rule: at least one secretary per non-archived congregation.

### Sessions

Keep token-hash sessions. Store `active_congregation_id` on the session (defaults to `users.congregation_id`). Only superadmin may change it. Non-superadmin: active must equal home; ignore or 403 if they try to switch.

### Passkeys and enrollment

- `webauthn_credentials` (user_id, credential id, public key, sign count, created_at).
- `enrollment_tokens` (hash, user_id, expires_at, used_at, kind: `passkey_qr` | `email_invite`). Single use. Short TTL (e.g. 30–60 minutes).

QR/code: secretary (or superadmin) mints a token; the person opens the link on **their** phone and finishes WebAuthn there. Biometrics never leave the device. Nobody enrolls a passkey “for” someone on the helper’s phone.

### Moves

Secretary of origin or destination, or superadmin, sets `users.congregation_id` (and the linked publisher/household as the move runbook will specify). History of reports stays with the old congregation unless a later idea says otherwise. v1: move the **person** (user + publisher row) to the new congregation; do **not** rewrite past reports into the new congregation (orphan or keep `publisher_id` history in the old congregation — implementation plan must pick one and test it). Preferred: keep historical reports in the **source** congregation (publisher row archived-or-left there) and create or reuse a publisher row in the destination. If that is too heavy for v1, block moves until that idea is designed; do not silently merge years of hours across congregations.

**v1 move (explicit):** user + current publisher row are reassigned to the destination congregation. Existing reports keep their old `congregation_id` and remain visible only in the source congregation’s secretary tools (publisher appears as moved-out / not in the destination directory until the new row exists). Simpler alternative if that split is painful: v1 ships **without** the move UI and only superadmin can reassign in SQL. The implementation plan must choose; default is **no move UI in the first slice**, only `congregation_id` on new rows.

### Shepherding

- **Never:** `shepherding_notes` table, remark fields meant as pastoral notes, file uploads of notes.
- **Yes (later idea, permission already reserved):** attention list (irregular / inactive), visit planned/done, no narrative.

## Authorization (v1)

`GET /auth/me` returns identity, home and active congregation, flags, stacked privileges, and a list of allowed nav keys so the SPA can hide routes.

| Actor | Can |
| --- | --- |
| Student | Meeting program only |
| Publisher (baptized or not) | Program; own assignments accept/reject; own monthly report; directory and groups with the agreed fields; edit own first/last name |
| Regular / special pioneer | Plus own hour goal/progress |
| Elder or ministerial servant | Plus all assignments in that congregation; attention list. Not others’ reports |
| Secretary | Today’s operator surfaces + enable logins (password, invite, QR) + grant/revoke secretary in **that** congregation + congregation settings. Not other congregations. Not grant superadmin |
| Superadmin | Secretary powers in the **active** congregation + platform switcher + create/edit/archive congregations + grant/revoke secretaries anywhere |

API: 401 if no session. 403 if authenticated but wrong congregation, missing privilege, or student hitting directory/report. Cross-congregation IDs never leak as 200 (403 or 404; pick **404** for foreign IDs so they do not confirm existence).

Directory and group JSON is a **narrow DTO**. No phone, email, address, baptism date, or report totals.

## UI and data flow

### Public landing

Today the only public screen is `/login`. That **is** the landing until a marketing page exists.

- **401** and **403:** navigate to `/login` (landing) and show a toast with the API `error` text (Spanish).
- **403 does not delete the cookie.**
- After a 403, the login/landing page must **not** auto-redirect back into the forbidden route (avoid a loop). If the session is still valid, show the toast and a single action that goes to the user’s **first allowed** app route (student/publisher → program; secretary → dashboard).
- **401:** session missing; stay on landing and sign in again.

Implement this in `frontend/src/lib/api.ts` (or a small wrapper) so every `fetch` behaves the same. Add a toast (shadcn/Base UI). Copy stays in Spanish.

### App shell

Nav from `/auth/me`:

- Student: Program
- Publisher: Program, Assignments, My report, Brothers, Groups
- Pioneer: plus Hours goal
- Elder / MS: plus Attention list (and full assignments, not only “mine”)
- Secretary / superadmin in active congregation: plus Home (dashboard), Publishers, Reports, Attendance, Congregation, Access (enable login, QR, invite)
- Superadmin only: congregation switcher in the header

Program and Assignments may be **empty states** until those ideas ship. Routes and 403 rules still exist so the portal does not lie about permissions.

Profile: edit own first and last name; change own password if they have one.

## Error handling

- JSON errors stay `{ "error": "..." }` until the consistent-error idea.
- Enrollment token expired or reused: **410**. SPA: toast + stay/return to landing or Access as appropriate (not a 403).
- Last-secretary removal: **409**.
- Invite email when SMTP is unset: **503** or hide the button if `/auth/me` says `inviteEmailEnabled: false`.

## Testing

Table-driven httptest:

- Cookie login still works; 401 without session.
- User in congregation A gets 404/403 for B’s rows.
- Student: program 200; directory, groups, report 403.
- Publisher: `PUT` own name 200; cannot change privileges or another publisher.
- Publisher cannot read others’ reports (403). Secretary can.
- Elder/MS: attention list 200; others’ reports 403.
- Cannot remove last secretary (409).
- Passkey token: second use and expiry → 410.
- Superadmin switches active congregation and then reads B; secretary has no switch endpoint (403).
- No shepherding-notes route or migration.

Frontend: nav matches `/auth/me`. Switcher only if `isSuperadmin`. A forced 403 (e.g. publisher opens `/reports`) lands on `/login` with a toast and no redirect loop.

Seed / migration: current congregation gets a UUID; existing publishers, households, reports, attendance attach to it; `ADMIN_EMAIL` user is superadmin and secretary of that congregation.

## Phasing (this spec is one product; several slices)

1. **Tenancy:** `congregations` + `congregation_id` on existing tables. One congregation in the UI. Tests for isolation (second congregation in fixtures only).
2. **Users:** new columns, `/auth/me` claims, secretary vs superadmin, self-edit name, gate existing routes.
3. **403/401 landing + toast** on the SPA.
4. **Access:** enable login, set password, QR enrollment, passkey login. Invite email behind SMTP flag.
5. **Publisher portal:** narrow directory; own report; empty Program/Assignments/Groups/Attention/Goals as needed.
6. **Superadmin:** create congregation, switcher, archive.

Do not start slice 5 screens that belong to a separate idea’s schema (groups, LMM) until that idea’s plan exists. Empty nav + 403 is enough.

## Related roadmap

- **Drop** `shepherding-notes` (storage). Keep a future **attention list / planning** idea.
- **Replace** `roles-beyond-admin` with this design (same outcome, not a second roles system).
- **Revise** the summary out-of-scope line: public signup and multi-congregation **SaaS signup** stay out; **one install, many congregations** is in.
- The fifteen product ideas approved in brainstorming (service-year card, groups, month lock, attention list, privilege history, Memorial, public talks, pioneer goals, self-report, CSV import, 2FA, login throttle, audit log, dark mode, missing-report email) remain separate folders after this spec is accepted. Self-report and 2FA/passkeys overlap this foundation — do not duplicate folders. Memorial attendance later moved into `attendance-duplicate-guard`; `memorial-attendance` is dropped.

## Open implementation choices (not product questions)

- Exact WebAuthn library on Go and the SPA.
- Toast component (Base UI / sonner-style).
- Whether v1 ships a move UI (default: **no**).
- NULL vs empty email for passkey-only users (prefer NULL + unique constraint that allows several NULLs).
