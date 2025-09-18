#!/bin/sh
set -eu

# Require endpoint
if [ -z "${VITE_API_ENDPOINT:-}" ]; then
  echo "No API endpoint given. Please define VITE_API_ENDPOINT in the environment."
  echo "Exiting."
  exit 1
fi

PLACEHOLDER="%VITE_API_ENDPOINT%"
# Set to wherever your built files live (e.g., /app, /app/dist). You used /app.
ROOT="${REPLACE_ROOT:-/app}"

FOUND=0
CHANGED=0
SKIPPED=0

# Only touch files that actually contain the placeholder
IFS='
'
for f in $(grep -rl --null -- "$PLACEHOLDER" "$ROOT" 2>/dev/null | tr '\0' '\n'); do
  FOUND=$((FOUND+1))
  dir="$(dirname "$f")"

  # Need both dir and file writable to overwrite safely
  if [ -w "$dir" ] && [ -w "$f" ]; then
    tmp="$(mktemp)"
    # No -i: write to tmp, then overwrite file (avoids cross-device rename/RO issues)
    if sed "s|$PLACEHOLDER|${VITE_API_ENDPOINT}|g" "$f" > "$tmp"; then
      if cat "$tmp" > "$f"; then
        CHANGED=$((CHANGED+1))
      else
        echo "warn: write failed: $f; skipping."
        SKIPPED=$((SKIPPED+1))
      fi
    else
      echo "warn: sed failed: $f; skipping."
      SKIPPED=$((SKIPPED+1))
    fi
    rm -f "$tmp" || true
  else
    echo "info: $f not writable; skipping."
    SKIPPED=$((SKIPPED+1))
  fi
done

echo "summary: found=$FOUND, changed=$CHANGED, skipped=$SKIPPED"

# Hand off to your app's CMD/args
if [ "$#" -gt 0 ]; then
  exec "$@"
fi
