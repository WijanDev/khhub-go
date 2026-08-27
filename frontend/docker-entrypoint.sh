#!/bin/sh
set -eu

# Runtime API origin so the same web image can serve staging and production.
url="${KHHUB_API_URL:-https://api.khhub.app}"
case "$url" in
  http://*|https://*) ;;
  *)
    echo "KHHUB_API_URL must start with http:// or https://" >&2
    exit 1
    ;;
esac
escaped=$(printf '%s' "$url" | sed 's/\\/\\\\/g; s/"/\\"/g')
printf 'window.__KHHUB_API_URL__="%s";\n' "$escaped" > /public/config.js

exec /usr/local/bin/static-web-server
