# Congregation meeting calendar and attendance uniqueness

- **Slug:** attendance-duplicate-guard
- **Status:** in-progress
- **Merge date:**
- **App version:**

## Summary

Each congregation owns a versioned weekly timetable (midweek weekday + time, weekend weekday + time) and a dated **exception calendar**. Attendance uniqueness is **per congregation**, never global. Regular meetings, the Memorial, special visits, and assembly (or similar) cancellations all go through that calendar so the secretary cannot double-count or enter a meeting that did not happen.

Absorbs `memorial-attendance`.

**Current slice:** `UNIQUE (meeting_date, kind)` already exists and `PUT /attendance` upserts. Saving again for a date and type already in the loaded month asks for confirmation before replacing counts. No 409 and no calendar in this slice.

## Rules

- Unique attendance key: `(congregation_id, meeting_date, kind)`.
- Regular `kind` values: `midweek`, `weekend`. The date must match the schedule **version in force** that day, unless an exception cancels or replaces that meeting.
- Exception types (per congregation, date or date range):
  - **Memorial** — own attendance row (`kind=memorial`), not a fake weekend. One Memorial row per congregation per date. Does not count in midweek/weekend month averages.
  - **Special visit** — own attendance row (`kind=special_visit`). May add a meeting and/or replace the regular midweek or weekend on those dates.
  - **Assembly** (circuit assembly, regional convention, or similar) — cancels the local midweek and/or weekend in the range. No local regular attendance row for a cancelled meeting. No assembly headcount in this slice (people attended elsewhere).
- Until `accounts-multicongregation` lands, the install is one implicit congregation.

## Out of scope

- Midweek assignments, public talks, or a full meeting program.
- Per-publisher check-in (Memorial or otherwise).
- jw.org reporting APIs or submitting S-88 automatically.
- Two standing midweek meetings every week (an extra meeting is a dated special-visit exception, not a second template slot).
- Editing attendance across service years as a bulk tool.
- Sharing one timetable across congregations.
