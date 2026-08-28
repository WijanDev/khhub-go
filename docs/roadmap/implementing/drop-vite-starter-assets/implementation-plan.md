# Implementation plan

1. Delete unused Vite leftovers: `frontend/src/assets/react.svg`, `frontend/src/assets/vite.svg`, `frontend/public/icons.svg`, and unused `App.tsx` / `App.css`.
2. Keep `frontend/public/favicon.svg` (linked from `index.html`) and `frontend/public/config.js`.
3. Grep for imports and run `npm run lint` and `npm run build`.
