#!/bin/sh

# =============================================================================
# Build Script for Go Wallet SDK
# =============================================================================
# Usage: ./build.sh [mode] [-i=module1,module2,...]
#
# Modes:
#   all           - Run all modules
#   failed        - Run only previously failed modules
#   mod1,mod2,... - Run only specified modules (comma-separated)
#   (none)        - Interactive prompt if failures exist
#
# Options:
#   -i=x,y  - Ignore modules matching these names (comma-separated)
#
# Examples:
#   ./build.sh                      # Interactive mode
#   ./build.sh all                  # Run all modules
#   ./build.sh failed               # Run previously failed
#   ./build.sh bitcoin,ethereum     # Run only bitcoin and ethereum
#   ./build.sh all -i=zksync        # Run all except zksync
# =============================================================================

set -e

# -----------------------------------------------------------------------------
# Configuration
# -----------------------------------------------------------------------------
ROOT_DIR=$(pwd)
LOG_FILE="$ROOT_DIR/build_failures.log"
TMP_FILE="$ROOT_DIR/build_failures.tmp"

# Counters
SUCCESS_COUNT=0
FAIL_COUNT=0
IGNORE_COUNT=0

# Lists for tracking
FAILED_LIST=""
FAILED_DISPLAY=""
IGNORED_LIST=""
IGNORED_DISPLAY=""

# -----------------------------------------------------------------------------
# Helper Functions
# -----------------------------------------------------------------------------

# Check if a module name is in a comma-separated list
# Usage: is_in_list "module_name" "item1,item2,item3"
# Returns: 0 if found, 1 if not found
is_in_list() {
  local name="$1"
  local list="$2"
  
  [ -z "$list" ] && return 1
  
  echo "$list" | tr ',' '\n' | while read pattern; do
    [ "$name" = "$pattern" ] && exit 0
  done
  return $?
}

# Get previous failures from log file
get_prev_failures() {
  [ -f "$LOG_FILE" ] || return
  grep -E "^# (coins|crypto|util|example)/" "$LOG_FILE" 2>/dev/null | sed 's/^# //'
}

# Run a build step and return result
# Usage: run_step "step_num" "step_name" "command"
run_step() {
  local step_num="$1"
  local step_name="$2"
  shift 2
  
  printf "  [%s] %-14s ... " "$step_num" "$step_name"
  
  if output=$("$@" 2>&1); then
    echo "✓ PASS"
    return 0
  else
    echo "✗ FAILED"
    STEP_OUTPUT="$output"
    return 1
  fi
}

# Build a single module
# Returns: 0 on success, 1 on failure
build_module() {
  local dir="$1"
  local name=$(basename "$dir")
  local relpath="${dir#$ROOT_DIR/}"
  
  echo "=========================================="
  echo "[$relpath]"
  echo "=========================================="
  cd "$dir"
  
  STEP_OUTPUT=""
  local failed_step=""
  
  # Step 1: go mod tidy
  if ! run_step "1/4" "go mod tidy" go mod tidy; then
    failed_step="go mod tidy"
  fi
  
  # Step 2: go mod edit
  if [ -z "$failed_step" ]; then
    if ! run_step "2/4" "go mod edit" go mod edit --toolchain=none; then
      failed_step="go mod edit"
    fi
  else
    echo "  [2/4] go mod edit    ... SKIPPED"
  fi
  
  # Step 3: go build
  if [ -z "$failed_step" ]; then
    if ! run_step "3/4" "go build" go build ./...; then
      failed_step="go build"
    fi
  else
    echo "  [3/4] go build       ... SKIPPED"
  fi
  
  # Step 4: go test
  if [ -z "$failed_step" ]; then
    if ! run_step "4/4" "go test" go test -v ./...; then
      failed_step="go test"
    fi
  else
    echo "  [4/4] go test        ... SKIPPED"
  fi
  
  echo ""
  
  # Record result
  if [ -n "$failed_step" ]; then
    echo "  ✗ $name FAILED at: $failed_step"
    FAILED_DISPLAY="$FAILED_DISPLAY  - $relpath ($failed_step)\n"
    FAILED_LIST="$FAILED_LIST# $relpath\n"
    FAIL_COUNT=$((FAIL_COUNT + 1))
    
    # Log failure details
    {
      echo "=========================================="
      echo "FAILED: $relpath ($failed_step)"
      echo "=========================================="
      echo "$STEP_OUTPUT"
      echo ""
    } >> "$TMP_FILE"
    
    echo ""
    return 1
  else
    echo "  ✓ $name ALL STEPS PASSED"
    SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    echo ""
    return 0
  fi
}

# Print summary
print_summary() {
  echo ""
  echo "=========================================="
  echo "                SUMMARY"
  echo "=========================================="
  echo "Success: $SUCCESS_COUNT"
  echo "Failed:  $FAIL_COUNT"
  echo "Ignored: $IGNORE_COUNT"
  
  if [ $FAIL_COUNT -gt 0 ]; then
    echo ""
    echo "Failed modules:"
    printf "$FAILED_DISPLAY"
    echo ""
    echo "See $LOG_FILE for details."
    echo "To rerun only failed: ./build.sh failed"
  fi
  
  if [ $IGNORE_COUNT -gt 0 ]; then
    echo ""
    echo "Ignored modules:"
    printf "$IGNORED_DISPLAY"
  fi
  
  echo ""
  echo "Usage: ./build.sh [all|failed|mod1,mod2,...] [-i=module1,module2,...]"
}

# Write log file
write_log() {
  if [ $FAIL_COUNT -gt 0 ] || [ $IGNORE_COUNT -gt 0 ]; then
    {
      [ $IGNORE_COUNT -gt 0 ] && {
        echo "# === IGNORED MODULES ==="
        printf "$IGNORED_LIST"
        echo ""
      }
      [ $FAIL_COUNT -gt 0 ] && {
        echo "# === FAILED MODULES ==="
        printf "$FAILED_LIST"
        echo ""
        cat "$TMP_FILE"
      }
    } > "$LOG_FILE"
  else
    rm -f "$LOG_FILE"
  fi
  rm -f "$TMP_FILE"
}

# -----------------------------------------------------------------------------
# Parse Arguments
# -----------------------------------------------------------------------------
MODE=""
IGNORE_LIST=""
INCLUDE_LIST=""

for arg in "$@"; do
  case "$arg" in
    -i=*)   IGNORE_LIST="${arg#-i=}" ;;
    all)    MODE="all" ;;
    failed) MODE="failed" ;;
    *,*)    MODE="select"; INCLUDE_LIST="$arg" ;;  # Contains comma = module list
    *)      
      # Single module name (no comma)
      if [ -z "$MODE" ] && [ -n "$arg" ]; then
        MODE="select"
        INCLUDE_LIST="$arg"
      fi
      ;;
  esac
done

# -----------------------------------------------------------------------------
# Discover Modules
# -----------------------------------------------------------------------------
echo "Working directory: $ROOT_DIR"

ALL_DIRS=$(find "$ROOT_DIR/coins" "$ROOT_DIR/crypto" "$ROOT_DIR/util" "$ROOT_DIR/example" \
  -name "go.mod" -type f -exec dirname {} \; | sort)

# -----------------------------------------------------------------------------
# Filter Modules (by include list and ignore list)
# -----------------------------------------------------------------------------
DIRS_TO_RUN=""

[ -n "$IGNORE_LIST" ] && echo "Ignoring modules: $IGNORE_LIST"
[ -n "$INCLUDE_LIST" ] && echo "Selecting modules: $INCLUDE_LIST"

for dir in $ALL_DIRS; do
  name=$(basename "$dir")
  relpath="${dir#$ROOT_DIR/}"
  
  # Check if ignored
  if [ -n "$IGNORE_LIST" ] && is_in_list "$name" "$IGNORE_LIST"; then
    IGNORED_DISPLAY="$IGNORED_DISPLAY  - $relpath\n"
    IGNORED_LIST="$IGNORED_LIST# $relpath\n"
    IGNORE_COUNT=$((IGNORE_COUNT + 1))
    continue
  fi
  
  # Check if in include list (when in select mode)
  if [ "$MODE" = "select" ] && [ -n "$INCLUDE_LIST" ]; then
    if is_in_list "$name" "$INCLUDE_LIST"; then
      DIRS_TO_RUN="$DIRS_TO_RUN $dir"
    fi
  else
    DIRS_TO_RUN="$DIRS_TO_RUN $dir"
  fi
done

# -----------------------------------------------------------------------------
# Determine Which Modules to Run
# -----------------------------------------------------------------------------
case "$MODE" in
  all)
    echo "Mode: all (running all modules)"
    ;;
  select)
    echo "Mode: select (running specified modules: $INCLUDE_LIST)"
    ;;
  failed)
    PREV_FAILED=$(get_prev_failures)
    if [ -n "$PREV_FAILED" ]; then
      echo "Mode: failed (running previously failed modules only)"
      DIRS_TO_RUN=""
      for relpath in $PREV_FAILED; do
        name=$(basename "$relpath")
        if ! is_in_list "$name" "$IGNORE_LIST"; then
          DIRS_TO_RUN="$DIRS_TO_RUN $ROOT_DIR/$relpath"
        fi
      done
    else
      echo "No previous failures found. Running all modules."
    fi
    ;;
  *)
    # Interactive mode
    PREV_FAILED=$(get_prev_failures)
    if [ -n "$PREV_FAILED" ]; then
      echo "=== Previous failures found ==="
      for relpath in $PREV_FAILED; do
        echo "  - $relpath"
      done
      echo ""
      printf "Run only failed modules? [y/N]: "
      read answer
      if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        echo "Running previously failed modules only."
        DIRS_TO_RUN=""
        for relpath in $PREV_FAILED; do
          name=$(basename "$relpath")
          if ! is_in_list "$name" "$IGNORE_LIST"; then
            DIRS_TO_RUN="$DIRS_TO_RUN $ROOT_DIR/$relpath"
          fi
        done
      else
        echo "Running all modules."
      fi
    fi
    ;;
esac

# -----------------------------------------------------------------------------
# Display Plan
# -----------------------------------------------------------------------------
[ $IGNORE_COUNT -gt 0 ] && {
  echo ""
  echo "=== Ignored modules ==="
  printf "$IGNORED_DISPLAY"
}

echo ""
echo "=== Modules to build ==="
for dir in $DIRS_TO_RUN; do
  echo "  - ${dir#$ROOT_DIR/}"
done
echo ""

# -----------------------------------------------------------------------------
# Build Modules
# -----------------------------------------------------------------------------
> "$TMP_FILE"

set +e  # Don't exit on build failures
for dir in $DIRS_TO_RUN; do
  build_module "$dir"
done
set -e

# -----------------------------------------------------------------------------
# Finalize
# -----------------------------------------------------------------------------
write_log
print_summary
