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

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

//go:generate ./bootstrap.sh

var (
	pkg = flag.String("p", "",
		"go package in which the version is declared "+
			"(default base of cwd)")
	unknown = flag.String("u", "unknown",
		"fallback version string if git describe fails")
	suffix = flag.String("s", "_version",
		"go file suffix (.go is appended)")
	dir    = flag.String("d", ".",
		"output directory for the generated Go file")
	name  = flag.String("c", "version", "name of the version const")
	vers  = flag.Bool("version", false, "print version and exit")
	pathf = flag.String("path", "", "path of the git command")
	optsf = flag.String("opts", "always",
		"comma-separated list of options for the git describe "+
			"invocation\n(each will be prefixed with '--')")
	argsf = flag.String("args", "",
		"comma-separated list of arguments for the git describe "+
			"invocation")
)

func main() {
	flag.Parse()

	if *vers {
		fmt.Printf("%s version: %s\n", os.Args[0], version)
		os.Exit(0)
	}

	versionVal := *unknown
	git := "git"

	if *pkg == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Cannot get cwd, giving up:", err)
			os.Exit(-1)
		}
		*pkg = path.Base(dir)
	}
	stat, err := os.Stat(*dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if !stat.IsDir() {
		fmt.Printf("%s is not a directory, giving up\n", *dir)
		os.Exit(-1)
	}
	fname := path.Join(*dir, *pkg+*suffix+".go")
	file, err := os.Create(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	gitargs := []string{"describe"}
	if *optsf != "" {
		opts := strings.Split(*optsf, ",")
		for _, opt := range opts {
			gitargs = append(gitargs, "--"+opt)
		}
	}
	if *argsf != "" {
		args := strings.Split(*argsf, ",")
		if len(args) > 0 {
			gitargs = append(gitargs, args...)
		}
	}
	if *pathf != "" {
		git = path.Join(*pathf, git)
	}

	bytes, err := exec.Command(git, gitargs...).CombinedOutput()
	if err != nil {
		fmt.Println("git describe failed:", err)
		if bytes != nil {
			fmt.Print("git describe output:\n", string(bytes))
		}
		fmt.Println(os.Args[0], "falling back to version:", versionVal)
	} else {
		versionVal = strings.TrimSpace(string(bytes))
	}

	fmt.Fprintf(file, "package %s\n\n", *pkg)
	fmt.Fprintf(file, "const %s = \"%s\"\n", *name, versionVal)
	os.Exit(0)
}
