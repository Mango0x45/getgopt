.. vi: tw=100 noet

Getgopt
=======

**Getgopt** is an extremely very simple implementation of POSIX C's ``getopt(3)`` function for
golang. The entire library consists of 4 exported global variables, 1 exported function, and less
than 100 lines of code (``sed -e '/^$/d' -e '/^\s*\*/d' -e '/^\s*\/\*/d' getgopt.go | wc -l``). I
wrote this library because I didn't like how all the other alternatives for normal command line
parsing were so overly complex. Go is all about keeping things simple, so let's keep flag parsing
simple.


Usage
=====

There is only 1 function for you to use, and that is ``getgopt.Getopt()``. The function works
*almost* the same way that the POSIX C ``getopt`` function works. Here is an example of it's usage:

.. code-block:: go

	package main

	import (
		"fmt"
		"os"

		"git.thomasvoss.com/getgopt"
	)

	func main() {
		for opt := byte(0); getgopt.Getopt(len(os.Args), os.Args, ":a:bcd", &opt); {
			switch opt {
			case 'a':
				fmt.Printf("Parsed the -a flag with the argument '%s'\n",
					getgopt.Optarg)
			case 'b':
				fmt.Println("Parsed the -b flag")
			case 'c':
				/* ... */
			case 'd':
				/* ... */
			case '?':
				fmt.Fprintf(os.Stderr, "Invalid flag '%c', read the manpage\n",
					getgopt.Optopt)
				os.Exit(1)
			case ':':
				fmt.Fprintf(os.Stderr, "The flag '%c' requires an argument\n",
					getgopt.Optopt)
				os.Exit(1)
			}
		}

		fmt.Printf("The first non-option argument is '%s'\n", os.Args[getgopt.Optind])
	}

After parsing a flag the ``Getopt()`` function returns true if there are still more flags to parse,
or false if there are none more. This means that we can use it in a ``for`` or ``while`` loop to
iterate over all of our arguments. As its arguments, the ``Getopt()`` function takes (in this
order), the count of command line arguments, the command line arguments, an *optstring*, and a
pointer to a byte where the parsed flag will be stored. After parsing a flag, the byte that was
passed as the last parameter will either have the value of the flag or one of ``':'`` and ``'?'``.
The value of ``opt`` is set to ``'?'`` if the user attempted to pass a flag that was not specified
by the given *optstring*. If the user specifies a flag that requires an argument without actually
passing an argument, then ``opt`` will be set to ``':'`` if the first character in the *optstring*
is ``':'`` and otherwise it will be ``'?'``.

The *optstring* is a string passed as the 3rd argument to ``Getopt()`` which specified which flags
you want to be able to handle. Each flag you want to handle is given as a single character in the
string in any order. For example if you want to support the ``-a``, ``-b``, and ``-x`` flags you can
do:

.. code-block:: go

	getgopt.Getopt(len(os.Args), os.Args, "abx", &opt)
	/* or */
	getgopt.Getopt(len(os.Args), os.Args, "bxa", &opt)

If you want a flag to take an argument, you should suffix the character with a ``':'``. So using the
above example, if we want the ``-b`` flag to take an argument, we could write:

.. code-block:: go

	getgopt.Getopt(len(os.Args), os.Args, "ab:x", &opt)

Finally, by default the ``Getopt()`` function will print diagnostic error messages to standard
output when the user fails to provide an argument to a flag that expects one or passes an invalid
flag. If you would like to not have these diagnostics printed you can either prefix the optstring
with ``':'`` or you can set the ``Opterr`` global variable to ``false``. Both of the following are
equivalent:

.. code-block:: go

	getgopt.Getopt(len(os.Args), os.Args, ":ab:x", &opt)
	/* or */
	getgopt.Opterr = false
	getgopt.Getopt(len(os.Args), os.Args, "ab:x", &opt)

There is a *slight* difference in behavior though which was explained above.

Finally, there are 3 other global variables you can access, these are ``Optarg``, ``Optind``, and
``Optopt``. When you parse a flag which requires an argument, that argument can be found as a string
in the ``Optarg`` variable. ``Optind`` is a variable which during the parsing of the flags holds the
index of command line argument being parsed. After the flags are parsed though it holds the index of
the first non-option argument in the provided argument list. ``Optopt`` functions similarly to the
byte you pass as the functions final argument, but it holds the flag which caused the last parsing
error.
