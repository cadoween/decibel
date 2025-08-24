#!/usr/bin/env bash
# stg-update-push.sh
# Cascade: update the given patch AND all patches after it in `stg series`,
# move their branches, and push to remote.

set -euo pipefail

REMOTE="origin"
BR_PREFIX=""

die() { echo "ERR: $*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || die "missing dependency: $1"; }
need git
need stg
need sed

# Normalize series lines: strip leading spaces, marker (+ - = >) and following spaces.
norm_series() {
  stg series | sed -E 's/^[[:space:]]*[\+\-\=\>][[:space:]]*//; s/^[[:space:]]+//'
}

# Load SERIES[] in the order printed by `stg series`
read_series_into_array() {
  mapfile -t SERIES < <(norm_series | sed '/^$/d')
  [ ${#SERIES[@]} -gt 0 ] || die "no patches in stack (run on an stgit-initialized branch)"
}

# Index of $1 in SERIES[], or -1
index_in_series() {
  local want="$1"
  for (( i=0; i<${#SERIES[@]}; i++ )); do
    if [ "${SERIES[$i]}" = "$want" ]; then
      echo "$i"; return 0
    fi
  done
  echo "-1"; return 1
}

update_and_push_patch() {
  local patch="$1"
  local safe_patch="${patch//\//-}"
  local branch="${BR_PREFIX}${safe_patch}"

  echo "  -> refresh $patch"
  stg refresh --patch "$patch"

  local sha
  sha="$(stg id "$patch")"
  [ -n "$sha" ] || die "failed to get SHA for '$patch'"

  echo "  -> branch -f $branch $sha"
  git branch -f "$branch" "$sha"

  echo "  -> push $REMOTE $branch (force-with-lease)"
  git push --force-with-lease "$REMOTE" "$branch"

  echo "     patch:  $sha"
  echo "     local:  $(git rev-parse "$branch")"
  echo "     remote: $(git rev-parse "$REMOTE/$branch")"
}

# -------- args --------
if [ $# -lt 1 ]; then
  cat <<'USAGE'
Usage:
  stg-update-push.sh [--remote origin] [--branch-prefix feature/] <patch> [<patch>...]

Behavior:
  For each <patch>, also updates every patch that comes AFTER it in `stg series`
  (i.e., higher in the stack), then moves & pushes their branches.

Examples:
  stg-update-push.sh user-repo
  stg-update-push.sh --branch-prefix feature/ user-model
  stg-update-push.sh --remote origin --branch-prefix feature/ user-model user-repo
USAGE
  exit 1
fi

PATCHES=()
while (( $# )); do
  case "$1" in
    --remote)        REMOTE="${2:-}"; shift 2 ;;
    --branch-prefix) BR_PREFIX="${2:-}"; shift 2 ;;
    *) PATCHES+=("$1"); shift ;;
  esac
done

git remote get-url "$REMOTE" >/dev/null 2>&1 || die "remote '$REMOTE' not found"
stg series >/dev/null 2>&1 || die "current branch is not stgit-initialized (run: stg init)"

read_series_into_array

echo "Stack (as printed by 'stg series'):"
printf '  %s\n' "${SERIES[@]}"

# Collect targets: each requested patch and all after it
declare -A TARGETS=()
for p in "${PATCHES[@]}"; do
  idx="$(index_in_series "$p")"
  if [ "$idx" -lt 0 ]; then
    echo "Requested patch '$p' not found. Available:"
    printf '  %s\n' "${SERIES[@]}"
    exit 1
  fi
  for (( i=idx; i<${#SERIES[@]}; i++ )); do
    TARGETS["${SERIES[$i]}"]=1
  done
done

# Preserve stack order
ORDERED_TARGETS=()
for name in "${SERIES[@]}"; do
  if [[ -n "${TARGETS[$name]+x}" ]]; then
    ORDERED_TARGETS+=("$name")
  fi
done

echo "Will update & push (cascade):"
printf '  %s\n' "${ORDERED_TARGETS[@]}"

for patch in "${ORDERED_TARGETS[@]}"; do
  echo "== Processing patch: $patch"
  update_and_push_patch "$patch"
  echo "✓ Done: $patch"
done

echo "All specified patches and their successors have been updated & pushed."
