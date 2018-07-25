#!/usr/bin/env bash
src="test-data"

readarray -d $'\0' dirs < <(find "$src" -type d -print0 |sort -z)

for dir in "${dirs[@]}"; do
    [[ -f $dir/stdout.expect || -f $dir/stderr.expect ]] && touch "$dir/untouched"
done

args=()
while (( "$#" )); do
  args+=($1)
  if [[ "$1" == "-run" ]]; then
    shift
    args+=("TestFirst|$1")
  fi
  shift
done

go test "${args[@]}" |& tee test.log

readarray -t log_lines <test.log

edit() {
    if [[ $EDITOR ]] ; then
        cat "$dir/$ft".actual
        ${EDITOR} "$dir/$ft".expect
    else
        echo "Actual file: $dir/$ft.actual"
        cat "$dir/$ft".actual
        echo "Expect file: $dir/$ft.expect"
        cat "$dir/$ft".expect
        read -p "Fix $ft.actual by editing $ft.expect.  Done?" ans
    fi
}

path_re='([^%])Test path: '
for line in "${log_lines[@]}"; do
    [[ $line = *' Test path:'* ]] || continue
    echo "Test for ${line}"
    dir="${line##*: }"
    for ft in stderr stdout; do
        [[ -f $dir/untouched ]] && continue
        [[ $(cat "$dir/$ft".actual) == "" && ! -f "$dir/$ft".expect ]] && continue
        [[ -f "$dir/$ft".expect ]] || touch "$dir/$ft".expect
        [[ -f $dir/$ft.actual || -f $dir/$ft.expect ]] || continue
        if [[ -f $dir/want-usage && ft = stdout ]]; then
           if grep -q 'Usage:' "$dir/$ft".actual; then
               true
           else
               echo "Expected usage"
               read -p "Press a key to continue" ans
           fi
        fi
        if grep -q '^RE:' "$dir/$ft".expect; then
            go run ../cmds/regex_test/testRe.go "$dir/$ft".expect "$dir/$ft".actual && continue
            edit
            continue
        fi
        diff -u "$dir/$ft".expect "$dir/$ft".actual && continue
        read -p "Move $ft.actual to $ft.expect? (y/e/n)" ans
        case $ans in
            y) mv "$dir/$ft".actual "$dir/$ft".expect;;
            e) edit;;
        esac
    done
done
