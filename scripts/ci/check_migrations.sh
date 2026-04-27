#!/usr/bin/env bash

set -euo pipefail

BASE_SHA="${1:-}"
HEAD_SHA="${2:-}"

if [[ -z "$BASE_SHA" || -z "$HEAD_SHA" ]]; then
  echo "Usage: $0 <base_sha> <head_sha>"
  exit 2
fi

if ! git cat-file -e "$BASE_SHA^{commit}" 2>/dev/null; then
  echo "Base commit not found locally: $BASE_SHA"
  exit 2
fi

if ! git cat-file -e "$HEAD_SHA^{commit}" 2>/dev/null; then
  echo "Head commit not found locally: $HEAD_SHA"
  exit 2
fi

mapfile -t migration_changes < <(git diff --name-status --find-renames "$BASE_SHA" "$HEAD_SHA" -- migrations)

if [[ ${#migration_changes[@]} -eq 0 ]]; then
  echo "No migration changes detected."
  exit 0
fi

echo "Detected migration changes:"
printf '%s\n' "${migration_changes[@]}"

violations=()
new_files=()

for line in "${migration_changes[@]}"; do
  status="${line%%$'\t'*}"

  if [[ "$status" == A ]]; then
    path="${line#*$'\t'}"
    new_files+=("$path")
    continue
  fi

  if [[ "$status" == R* ]]; then
    old_path="$(echo "$line" | awk -F '\t' '{print $2}')"
    new_path="$(echo "$line" | awk -F '\t' '{print $3}')"
    violations+=("Rename is not allowed for migration files: $old_path -> $new_path")
    continue
  fi

  if [[ "$status" == D ]]; then
    path="${line#*$'\t'}"
    violations+=("Delete is not allowed for migration files: $path")
    continue
  fi

  if [[ "$status" == M ]]; then
    path="${line#*$'\t'}"
    violations+=("Editing existing migration is not allowed: $path")
    continue
  fi

  violations+=("Unsupported migration change detected ($status): $line")
done

name_regex='^migrations/[0-9]{14}_[a-z0-9][a-z0-9_-]*\.sql$'
for path in "${new_files[@]}"; do
  if [[ ! "$path" =~ $name_regex ]]; then
    violations+=("Invalid migration filename: $path (expected migrations/YYYYMMDDHHMMSS_description.sql)")
  fi
done

if [[ ${#violations[@]} -gt 0 ]]; then
  echo
  echo "Migration policy violations found:"
  for violation in "${violations[@]}"; do
    printf ' - %s\n' "$violation"
    echo "::error::$violation"
  done
  exit 1
fi

echo "Migration policy check passed."
