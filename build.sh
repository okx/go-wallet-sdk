#!/bin/sh

echo "Working directory: $(pwd)"
echo "Working directory: $(pwd)"

p=$(pwd)
logfile="$p/build_failures.log"
tmpfile="$p/build_failures.tmp"
success_count=0
fail_count=0
ignore_count=0

# Parse arguments
mode=""
ignore_list=""
for arg in "$@"; do
  case "$arg" in
    -i=*)
      ignore_list="${arg#-i=}"
      ;;
    all|failed)
      mode="$arg"
      ;;
  esac
done

# Find all directories containing go.mod under coins/, crypto/, and util/
all_dirs=$(find "$p/coins" "$p/crypto" "$p/util" -name "go.mod" -type f -exec dirname {} \; | sort)

# Filter out ignored directories
filtered_dirs=""
ignored_dirs=""
ignored_list=""
ignored_display=""

if [ -n "$ignore_list" ]; then
  echo "Ignoring modules matching: $ignore_list"
  for dir in $all_dirs; do
    name=$(basename "$dir")
    relpath="${dir#$p/}"
    is_ignored=false
    
    # Check if directory name matches any ignore pattern
    IFS=','
    for pattern in $ignore_list; do
      if [ "$name" = "$pattern" ]; then
        is_ignored=true
        break
      fi
    done
    unset IFS
    
    if [ "$is_ignored" = true ]; then
      ignored_dirs="$ignored_dirs $dir"
      ignored_display="$ignored_display  - $relpath\n"
      ignored_list="$ignored_list# $relpath\n"
      ignore_count=$((ignore_count + 1))
    else
      filtered_dirs="$filtered_dirs $dir"
    fi
  done
  all_dirs="$filtered_dirs"
fi

# Check for previous failures and determine which dirs to run
dirs_to_run=""

if [ "$mode" = "all" ]; then
  echo "Mode: all (running all modules)"
  dirs_to_run="$all_dirs"
elif [ "$mode" = "failed" ]; then
  if [ -f "$logfile" ]; then
    # Parse failed dirs from log file (lines starting with "# coins/", "# crypto/", or "# util/")
    prev_failed=$(grep -E "^# (coins|crypto|util)/" "$logfile" | sed 's/^# //')
    if [ -n "$prev_failed" ]; then
      echo "Mode: failed (running previously failed modules only)"
      dirs_to_run=""
      for relpath in $prev_failed; do
        # Skip if this path is in the ignore list
        name=$(basename "$relpath")
        is_ignored=false
        if [ -n "$ignore_list" ]; then
          IFS=','
          for pattern in $ignore_list; do
            if [ "$name" = "$pattern" ]; then
              is_ignored=true
              break
            fi
          done
          unset IFS
        fi
        if [ "$is_ignored" = false ]; then
          dirs_to_run="$dirs_to_run $p/$relpath"
        fi
      done
    else
      echo "No previous failures found. Running all modules."
      dirs_to_run="$all_dirs"
    fi
  else
    echo "No log file found. Running all modules."
    dirs_to_run="$all_dirs"
  fi
else
  # No arg provided - check for log file and prompt
  if [ -f "$logfile" ]; then
    prev_failed=$(grep -E "^# (coins|crypto|util)/" "$logfile" | sed 's/^# //')
    if [ -n "$prev_failed" ]; then
      echo "=== Previous failures found ==="
      for relpath in $prev_failed; do
        echo "  - $relpath"
      done
      echo ""
      printf "Run only failed modules? [y/N]: "
      read answer
      if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
        echo "Running previously failed modules only."
        dirs_to_run=""
        for relpath in $prev_failed; do
          # Skip if this path is in the ignore list
          name=$(basename "$relpath")
          is_ignored=false
          if [ -n "$ignore_list" ]; then
            IFS=','
            for pattern in $ignore_list; do
              if [ "$name" = "$pattern" ]; then
                is_ignored=true
                break
              fi
            done
            unset IFS
          fi
          if [ "$is_ignored" = false ]; then
            dirs_to_run="$dirs_to_run $p/$relpath"
          fi
        done
      else
        echo "Running all modules."
        dirs_to_run="$all_dirs"
      fi
    else
      dirs_to_run="$all_dirs"
    fi
  else
    dirs_to_run="$all_dirs"
  fi
fi

if [ $ignore_count -gt 0 ]; then
  echo ""
  echo "=== Ignored modules ==="
  printf "$ignored_display"
fi

echo ""
echo "=== Modules to build ==="
for dir in $dirs_to_run; do
  echo "  - ${dir#$p/}"
done
echo ""

# Initialize temp file for detailed output
> "$tmpfile"
failed_list=""
failed_display=""

for dir in $dirs_to_run; do
  name=$(basename "$dir")
  relpath="${dir#$p/}"
  echo "Building $name..."
  cd "$dir"
  
  # Run go mod tidy and check for failure
  tidy_output=$(go mod tidy 2>&1)
  tidy_exit=$?
  
  if [ $tidy_exit -ne 0 ]; then
    echo "✗ Build $name FAILED (go mod tidy failed)."
    echo ""
    failed_display="$failed_display  - $relpath\n"
    failed_list="$failed_list# $relpath\n"
    fail_count=$((fail_count + 1))
    
    # Append detailed output to temp file
    echo "==========================================" >> "$tmpfile"
    echo "FAILED: $relpath (go mod tidy)" >> "$tmpfile"
    echo "==========================================" >> "$tmpfile"
    echo "$tidy_output" >> "$tmpfile"
    echo "" >> "$tmpfile"
    continue
  fi
  
  go mod edit --toolchain=none
  
  # Capture test output
  output=$(go test -v ./... 2>&1)
  exit_code=$?
  
  if [ $exit_code -eq 0 ]; then
    echo "✓ Build $name success."
    echo ""
    success_count=$((success_count + 1))
  else
    echo "✗ Build $name FAILED."
    echo ""
    failed_display="$failed_display  - $relpath\n"
    failed_list="$failed_list# $relpath\n"
    fail_count=$((fail_count + 1))
    
    # Append detailed output to temp file
    echo "==========================================" >> "$tmpfile"
    echo "FAILED: $relpath" >> "$tmpfile"
    echo "==========================================" >> "$tmpfile"
    echo "$output" >> "$tmpfile"
    echo "" >> "$tmpfile"
  fi
done

# Write final log file with failed list at top
if [ $fail_count -gt 0 ] || [ $ignore_count -gt 0 ]; then
  {
    if [ $ignore_count -gt 0 ]; then
      echo "# === IGNORED MODULES ==="
      printf "$ignored_list"
      echo ""
    fi
    if [ $fail_count -gt 0 ]; then
      echo "# === FAILED MODULES ==="
      printf "$failed_list"
      echo ""
      cat "$tmpfile"
    fi
  } > "$logfile"
else
  # No failures or ignores - remove log file
  rm -f "$logfile"
fi
rm -f "$tmpfile"

echo ""
echo "=========================================="
echo "                SUMMARY"
echo "=========================================="
echo "Success: $success_count"
echo "Failed:  $fail_count"
echo "Ignored: $ignore_count"

if [ $fail_count -gt 0 ]; then
  echo ""
  echo "Failed modules:"
  printf "$failed_display"
  echo ""
  echo "See $logfile for details."
  echo ""
  echo "To rerun only failed: ./build.sh failed"
fi

if [ $ignore_count -gt 0 ]; then
  echo ""
  echo "Ignored modules:"
  printf "$ignored_display"
fi

echo ""
echo "Usage: ./build.sh [all|failed] [-i=module1,module2,...]"
