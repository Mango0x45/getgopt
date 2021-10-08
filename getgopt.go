/*
 * BSD Zero Clause License
 *
 * Copyright (c) 2021 Thomas Voss
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
 * REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY
 * AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
 * INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM
 * LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR
 * OTHER TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
 * PERFORMANCE OF THIS SOFTWARE.
 */

package getgopt

import (
	"fmt"
	"os"
)

const (
	errNoArg  = "option requires an argument -- %c\n"
	errBadArg = "unknown option -- %c\n"
)

var (
	opt    = 1
	optlen int
	parsed bool
	opts   [255]option
)

var (
	/* The argument to the matched flag */
	Optarg string
	/* If true, print error messages */
	Opterr = true
	/* During parsing this is the current index into os.Args being parsed. After parsing this is
	 * the index to the first non-option element of os.Args.
	 */
	Optind = 1
	/* The flag that caused the last parsing error */
	Optopt byte
)

type option struct {
	arg  bool
	used bool
}

func parseArgs(optstring string) {
	length := len(optstring)

	if (length > 0) && (optstring[0] == ':') {
		Opterr = false
	}

	for i, c := range optstring {
		if c != ':' {
			opts[c].used = true
			opts[c].arg = ((i < length-1) && (optstring[i+1] == ':'))
		}
	}
}

/* A function to parse command line flags. This function takes as it's arguments from first to last,
 * the count of command line arguments, the array of command line arguments (os.Args), an option
 * string, and a pointer to a byte where the current flag can be stored. When called the current
 * flag will be stored in `optptr` and the global variables `Optarg`, `Opterr`, `Optind`, and
 * `Optopt` may be set.
 *
 * If there are still more arguments to be parsed, the function will return true. Otherwise false is
 * returned. This makes it very easy to incorperate into a for/while loop.
 */
func Getopt(argc int, argv []string, optstring string, optptr *byte) bool {
	/* If we haven't parsed the optstring yet, parse it */
	if !parsed {
		parseArgs(optstring)
		parsed = true
	}

	/* Instantly return false if the follow cases are met */
	if Optind >= argc || argv[Optind] == "" || argv[Optind][0] != '-' ||
			argv[Optind] == "-" || optstring == "" {
		return false
	} else if argv[Optind] == "--" {
		Optind++
		return false
	}

	/* For each element of argv we calculate its length */
	if opt == 1 {
		optlen = len(argv[Optind])
	}

	/* The current flag */
	currFlag := argv[Optind][opt]

	if opts[currFlag].used {
		if opts[currFlag].arg {
			if opt == optlen-1 {
				Optind += 2
				if Optind > argc {
					Optopt = currFlag
					if Opterr {
						*optptr = '?'
						fmt.Fprintf(os.Stderr, errNoArg, Optopt)
					} else {
						*optptr = ':'
					}
				} else {
					Optarg = argv[Optind-1]
					*optptr = currFlag
				}
				opt = 1
				return true
			}

			/* If the opt takes an argument but it's not the last character in the
			 * string
			 */
			*optptr = currFlag
			Optarg = string(argv[Optind][opt+1:])
			Optind++
			opt = 1
		} else { /* If the opt doesn't take an argument */
			if opt == optlen-1 {
				opt = 1
				Optind++
			} else {
				opt++
			}

			*optptr = currFlag
		}
	} else { /* If the arg isn't in optstring */
		if Opterr {
			fmt.Fprintf(os.Stderr, errBadArg, argv[Optind][opt])
		}
		Optopt = currFlag
		*optptr = '?'

		if opt == optlen-1 {
			Optind++
			opt = 1
		} else {
			opt++
		}
	}

	return true
}
