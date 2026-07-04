#!/bin/sh
# Integration checks for yup-tac, run inside a Debian (GNU coreutils) container.
#
# parity CASE  — yup-tac must produce byte-identical output to GNU `tac`.
# assert WANT  — yup-tac must produce WANT exactly (used where yup-tac diverges
#                from GNU by design; see cmd-tac COMPATIBILITY.md).
set -eu

fails=0

# parity_stdin INPUT ARGS... — feed INPUT to both yup-tac and GNU tac on stdin.
parity_stdin() {
	input=$1
	shift
	ours=$(printf '%s' "$input" | yup-tac "$@" 2>/dev/null || true)
	gnu=$(printf '%s' "$input" | tac "$@" 2>/dev/null || true)
	if [ "$ours" = "$gnu" ]; then
		printf 'ok    parity  tac %s\n' "$*"
	else
		printf 'FAIL  parity  tac %s\n        gnu:  %s\n        ours: %s\n' "$*" "$gnu" "$ours"
		fails=$((fails + 1))
	fi
}

# parity_files ARGS... — both read the same file operand(s).
parity_files() {
	ours=$(yup-tac "$@" 2>/dev/null || true)
	gnu=$(tac "$@" 2>/dev/null || true)
	if [ "$ours" = "$gnu" ]; then
		printf 'ok    parity  tac %s\n' "$*"
	else
		printf 'FAIL  parity  tac %s\n        gnu:  %s\n        ours: %s\n' "$*" "$gnu" "$ours"
		fails=$((fails + 1))
	fi
}

# assert_files WANT ARGS... — yup-tac on file operand(s) must equal WANT exactly.
assert_files() {
	want=$1
	shift
	got=$(yup-tac "$@" 2>/dev/null || true)
	if [ "$got" = "$want" ]; then
		printf 'ok    assert  tac %s\n' "$*"
	else
		printf 'FAIL  assert  tac %s\n        want: %s\n        got:  %s\n' "$*" "$want" "$got"
		fails=$((fails + 1))
	fi
}

# assert_stdin WANT INPUT ARGS... — yup-tac on stdin must equal WANT exactly.
assert_stdin() {
	want=$1
	input=$2
	shift 2
	got=$(printf '%s' "$input" | yup-tac "$@" 2>/dev/null || true)
	if [ "$got" = "$want" ]; then
		printf 'ok    assert  tac %s\n' "$*"
	else
		printf 'FAIL  assert  tac %s\n        want: %s\n        got:  %s\n' "$*" "$want" "$got"
		fails=$((fails + 1))
	fi
}

# Default: reverse the order of lines read from stdin.
parity_stdin "alpha
beta
gamma
"
parity_stdin "one
two
three
four
five
"

# Single line in, single line out (nothing to reverse).
parity_stdin "solo
"

# Empty input yields empty output.
parity_stdin ""

# File operands: reverse each file's lines.
printf 'one\ntwo\nthree\n' >/tmp/a.txt
printf 'red\ngreen\nblue\n' >/tmp/b.txt
parity_files /tmp/a.txt

# Documented divergence: with multiple operands GNU tac reverses each file's
# lines independently and emits them in operand order, while yup-tac concatenates
# every file into one stream and reverses that stream as a whole. For operands
# b.txt (red,green,blue) then a.txt (one,two,three) the concatenation reverses to
# three,two,one,blue,green,red. The two are not byte-identical — see cmd-tac
# COMPATIBILITY.md. Assert yup-tac's contract.
assert_files "three
two
one
blue
green
red" /tmp/b.txt /tmp/a.txt

# Documented divergence: -s record separator. GNU tac treats the separator as a
# trailing record terminator, preserving it on each emitted record (so `a:b:c`
# reverses to `c\nb:a:`). yup-tac instead joins the input with newlines, splits
# on the literal separator, reverses the records, and rejoins them with newlines
# (yielding `c\nb\na\n`). The two are not byte-identical — see cmd-tac
# COMPATIBILITY.md. Assert yup-tac's contract.
assert_stdin "c
b
a" "a:b:c
" -s :

if [ "$fails" -ne 0 ]; then
	printf '\n%s check(s) failed\n' "$fails"
	exit 1
fi
printf '\nall checks passed\n'
