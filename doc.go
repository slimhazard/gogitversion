/*-
 * Copyright (c) 2018 UPLEX Nils Goroll Systemoptimierung
 * All rights reserved
 *
 * Author: Geoffrey Simmons <geoffrey.simmons@uplex.de>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

/* 
gogitversion generates a Go source file that declares a const for a
version string, whose contents are the output of 'git describe'.  It
is intended to be useful for 'go generate'.

To use it with go generate, include a directive like this in one of
your source files:

	//go:generate gogitversion

There must be no space between "//" and "go". Then issue the 'go
generate' command before 'go build'.

By default, gogitversion takes the base name of the current working
directory as the package for which it generates code. If the package
name is "foo", then by default a file named "foo_version.go" is
generated with content like this:

	package foo

	const version = "v1.0.4-14-g2414721"

... where "v1.0.4-14-g2414721" is the output from 'git describe'. The
version const can now be used in package foo.

See git-describe(1) for details about the git command.

Options

The code generation can be configured with command-line flags:

	-p pkg

Use 'pkg' as the package name. This name also appears in the name of
the generated file. The default package name is the base name of the
current working directory. For example, use '-p main' to declare the
const in package main.

	-c name

Use 'name' as the generated const name. For example, use '-c Version'
to create the const Version (which is then exported). The default
const name is 'version' (and hence is not exported).

	-s suffix

Use 'suffix' as the part of the file name that follows the package
name and precedes ".go". For example, if you use '-p foo' and '-s
Version', then the generated file is named 'fooVersion.go'. Use -s ''
(the empty string) to leave out the suffix. The default suffix is
'_version' (to create file names like 'foo_version.go').

	-u fallback

Use "fallback" as the version string if the 'git describe' invocation
fails. The default fallback string is "unknown".

	-d dir

Create the generated Go source in directory 'dir'. By default, it is
created in the current working directory.

	--opts opt1,opt2,opt3=val

Use '--opt1', '--opt2' and '--opt3=val' as command line options in the
'git describe' invocation. The contents of --opts are comma-separated,
and should not have '--' (otherwise they may be interpreted as options
for gogitversion, probably leading to failure). Options that have an
argument (such as --abbrev and --match) should include the equals sign
followed by the argument.

For example, '-opts all,abbrev=8,match=prod' will cause 'git describe'
to be called with '--all --abbrev=8 --match=prod'.

The default value of 'opts' is 'always' (so that 'git describe' will
always return a value, even if you have never created a tag). If you
use -opts to add other options and want to retain '--always', remember
to include 'always' in the list. Use -opts '' (the empty string) to
call 'git describe' without any options.

	--args arg1,arg2,arg3

Use 'arg1', 'arg2' and 'arg3' as arguments in the 'git describe'
invocation, following the options.  The contents of -args are
comma-separated. By default, 'git describe' is called with no
arguments (only the '--always' option).

	--path /opt/bin

Use '/opt/bin' as the path for 'git' (in other words, call
'/opt/bin/git').  By default, 'git' is invoked without a path, so the
binary is found via the $PATH variable.

The --version option causes gogitversion to print its version and exit
(this tool eats its own dog food).

Examples

Generate the file "main_version.go" in which the const 'version'
containing the output of 'git describe --always' is declared in
package main:

	gogitversion -p main

Generate the file "/tmp/fooversion.go" in which the exported const
'Version' containing the output of 'git describe --always' is declared
in package foo:

	gogitversion -p foo -s version -v Version -d /tmp

With the base name of the current working directory as the package
name, generate a version const using options and an argument for 'git
describe':

	gogitversion --opts all,abbrev=4 --args HEAD^

The version string contains the output of:

	git describe --all --abbrev=4 HEAD^

Building gogitversion

Bootstrapping is a bit tricky here, because the tool uses itself to
generate its version string. For best results, use the accompanying
Makefile (default target or 'all'), which takes care of things if you
have never built or installed gogitversion.

You must call 'go generate' before 'go build' (the Makefile does this
for you).

If you have been testing gogitversion and have been generating Go
sources in your working directory, make sure to delete them before
building gogitversion again (with make or 'go build'). Otherwise, 'go
build' will attempt to compile them as sources for gogitversion,
likely leading to errors. Using the -d option helps to avoid this
problem.
*/
package main
