# Utility analysis

- **Score:** 2
- **Scale:** 1 = low … 5 = high

## Justification

The secretary never sees this. Operators and agents hit it whenever they install frontend deps or when CI builds the image. Today local Node 24.11 / npm 11.6 rewrites `package-lock.json` while CI stays on Node 22; skipping the bump leaves that split and keeps the build image on Maintenance LTS until April 2027. That is occasional toolchain hygiene, not a weekly congregation task.
