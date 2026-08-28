# Implementation plan

Current slice (overwrite confirm):

1. Keep `PUT /attendance` as an upsert. Do not add a migration; `UNIQUE (meeting_date, kind)` is already in `000001_init`. Do not return 409.
2. Narrow `putAttendance` to an `attendanceWriter` interface so httptest can use a memory store (`backend/internal/http/attendance.go`).
3. Add `backend/internal/http/attendance_test.go`: two `PUT`s for the same date and kind keep one id and one row; invalid date or kind return 400.
4. In `frontend/src/features/attendance.tsx`, if the loaded month list already has that date and kind, open the shadcn Base UI `AlertDialog`. A new date and kind saves immediately. Spanish copy: “¿Actualizar esta reunión?”; Cancelar / Actualizar.
5. Verify: `cd backend && go test ./...`. In the browser, save a new meeting, save the same date and kind again, Cancel vs Actualizar.

Later (calendar + Memorial; not this PR):

1. Add `congregation_meeting_schedules`: `congregation_id`, midweek weekday + time, weekend weekday + time, `effective_from`, `effective_to` (nullable = current). No overlapping ranges per congregation. Seed the current (or implicit) congregation with one open version.
2. Add `congregation_meeting_exceptions`: `congregation_id`, `type` (`memorial` | `special_visit` | `assembly`), `from_date`, `to_date`, which regular kinds it **cancels** (`midweek`, `weekend`, both, or none), optional short label (not pastoral notes). No overlapping **same type** on the same date per congregation unless the design later allows it; Memorial is one range per year in practice.
3. Widen `meeting_attendance.kind` to `midweek` | `weekend` | `memorial` | `special_visit`. Replace `UNIQUE (meeting_date, kind)` with `UNIQUE (congregation_id, meeting_date, kind)` when `congregation_id` exists (or keep the single-install unique until accounts, then widen — never leave a global unique after tenancy).
4. `PUT /attendance`: resolve schedule version + exceptions for `(congregation, date)`. 422 if a regular kind’s weekday does not match the version, or if an exception cancelled that kind. 422 if `memorial` / `special_visit` has no matching exception that day. 409 on unique violation. Stable JSON errors.
5. Congregation settings UI (Spanish): current timetable + “new version”; read-only past versions; exception list (add Memorial date, special-visit range, assembly range). Attendance form: date + kind constrained by that calendar; show why a date is blocked (asamblea, etc.).
6. Dashboard: midweek/weekend averages ignore memorial and special-visit rows and skip cancelled dates (do not treat a missing cancelled weekend as zero). Surface the latest Memorial total for the service year.
7. Handler tests: two congregations may share a date; same congregation + date + kind is 409; cancelled assembly weekend is 422 for `weekend`; Memorial allowed only on the exception date; old schedule version still accepts its weekday. Browser: enter Memorial, a visit, and an assembly weekend.
