#!/usr/bin/env sh

compare()
{
	[ "$2" = "$3" ] && printf "\033[38;5;10mSuccess:\033[39;49m %s\n" "$1" ||
		printf "\033[38;5;9mFail:\033[39;49m %s <expected \`\`%s'' but got \`\`%s''>\n" \
			"$1" "$2" "$3"
}

cd "${0%/*}"

trap 'rm -f tests tests2 tests2.go' EXIT
go build tests.go

compare "no args" "" "$(./tests)"
compare "-a with valid arg" "Valid flag 'a' with arg 'testy'" "$(./tests -a testy test)"
compare "-a with no arg" "Valid flag 'a' with no arg" "$(./tests -a)"
compare "-a with valid arg and no space" "Valid flag 'a' with arg 'testy'" "$(./tests -atesty test)"
compare "-x with no args" "Valid flag 'x'" "$(./tests -x)"
compare "-x with args" "Valid flag 'x'" "$(./tests -x testy test)"
compare "-x and -a with args" "Valid flag 'x'
Valid flag 'a' with arg 'testy test'" "$(./tests -x -a 'testy test')"
compare "-xa with args" "Valid flag 'x'
Valid flag 'a' with arg 'testy test'" "$(./tests -xa 'testy test')"
compare "-ax with args" "Valid flag 'a' with arg 'x'" "$(./tests -ax 'testy test')"
compare "-ax with args" "Valid flag 'a' with arg 'x'" "$(./tests -ax 'testy test')"
compare "-x after --" "" "$(./tests -- -x)"
compare "-a with args after --" "" "$(./tests -- -a testy test)"
compare "-a with args then -x after --" "Valid flag 'a' with arg 'testy'" \
	"$(./tests -a testy test -- -x)"
compare "-a with args then -x after empty string" "Valid flag 'a' with arg 'testy'" \
	"$(./tests -a testy test '' -x)"
compare "-x chained 3 times" "Valid flag 'x'
Valid flag 'x'
Valid flag 'x'" "$(./tests -xxx)"
compare "-x as arg to -a" "Valid flag 'a' with arg '-x'" "$(./tests -a -x)"
compare "invalid flag -b" "Invalid flag 'b'" "$(./tests -b)"
compare "invalid flag -b with args" "Invalid flag 'b'" "$(./tests -b testy test)"
compare "-x after non option arg" "" "$(./tests testy -x)"
compare "-x after -" "" "$(./tests testy - -x)"

sed '/Getopt(/s/:a:x/a:x/' tests.go >tests2.go
go build tests2.go

compare "-a with no arg and optstring[0] != ':'" "option requires an argument -- a
Invalid flag 'a'" "$(2>&1 ./tests2 -a)"
compare "invalid flag -b and optstring[0] != ':'" "unknown option -- b
Invalid flag 'b'" "$(2>&1 ./tests2 -b)"
