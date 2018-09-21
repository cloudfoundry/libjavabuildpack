/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/bouk/monkey"
	"github.com/cloudfoundry/libjavabuildpack"
)

// CaptureExitStatus returns a pointer to the exit status code when os.Exit() is called.  Returns a function for use
// with defer in order to clean up after capture.
//
// c, d := CaptureExitStatus(t)
// defer d()
func CaptureExitStatus(t *testing.T) (*int, func()) {
	t.Helper()

	code := math.MinInt64
	pg := monkey.Patch(os.Exit, func(c int) {
		code = c
	})

	return &code, func() {
		pg.Unpatch()
	}
}

// Console represents the standard console objects, stdin, stdout, and stderr.
type Console struct {
	errRead  *os.File
	errWrite *os.File
	inRead   *os.File
	inWrite  *os.File
	outRead  *os.File
	outWrite *os.File
}

// Err returns a string representation of captured stderr.
func (c Console) Err(t *testing.T) string {
	t.Helper()

	err := c.errWrite.Close()
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(c.errRead)
	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}

// In writes a string and closes the connection once complete.
func (c Console) In(t *testing.T, string string) {
	t.Helper()

	_, err := fmt.Fprint(c.inWrite, string)
	if err != nil {
		t.Fatal(err)
	}

	err = c.inWrite.Close()
	if err != nil {
		t.Fatal(err)
	}
}

// Out returns a string representation of captured stdout.
func (c Console) Out(t *testing.T) string {
	t.Helper()

	err := c.outWrite.Close()
	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(c.outRead)
	if err != nil {
		t.Fatal(err)
	}

	return string(bytes)
}

// FindRoot returns the root of project being tested.
func FindRoot(t *testing.T) string {
	t.Helper()

	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for {
		if dir == "/" {
			t.Fatalf("could not find go.mod in the directory hierarchy")
		}
		if exist, err := libjavabuildpack.FileExists(filepath.Join(dir, "go.mod")); err != nil {
			t.Fatal(err)
		} else if exist {
			return dir
		}
		dir, err = filepath.Abs(filepath.Join(dir, ".."))
		if err != nil {
			t.Fatal(err)
		}
	}
}

// ProtectEnv protects a collection of environment variables.  Returns a function for use with defer in order to reset
// the previous values.
//
// defer ProtectEnv(t, "alpha")()
func ProtectEnv(t *testing.T, keys ...string) func() {
	t.Helper()

	type state struct {
		value string
		ok    bool
	}

	previous := make(map[string]state)
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		previous[key] = state{value, ok}
	}

	return func() {
		for k, v := range previous {
			if v.ok {
				os.Setenv(k, v.value)
			} else {
				os.Unsetenv(k)
			}
		}
	}
}

// ReplaceArgs replaces the current command line arguments (os.Args) with a new collection of values.  Returns a
// function suitable for use with defer in order to reset the previous values
//
//  defer ReplaceArgs(t, "alpha")()
func ReplaceArgs(t *testing.T, args ...string) func() {
	t.Helper()

	previous := os.Args
	os.Args = args

	return func() { os.Args = previous }
}

// ReplaceConsole replaces the console files (os.Stderr, os.Stdin, os.Stdout).  Returns a function for use with defer in
// order to reset the previous values
//
// c, d := ReplaceConsole(t)
// defer d()
func ReplaceConsole(t *testing.T) (Console, func()) {
	t.Helper()

	var console Console
	var err error

	errPrevious := os.Stderr
	console.errRead, console.errWrite, err = os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = console.errWrite

	inPrevious := os.Stdin
	console.inRead, console.inWrite, err = os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdin = console.inRead

	outPrevious := os.Stdout
	console.outRead, console.outWrite, err = os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = console.outWrite

	return console, func() {
		os.Stderr = errPrevious
		os.Stdin = inPrevious
		os.Stdout = outPrevious
	}
}

// ReplaceEnv replaces an environment variable.  Returns a function for use with defer in order to reset the previous
// value.
//
// defer ReplaceEnv(t, "alpha", "bravo")()
func ReplaceEnv(t *testing.T, key string, value string) func() {
	t.Helper()

	previous, ok := os.LookupEnv(key)
	os.Setenv(key, value)

	return func() {
		if ok {
			os.Setenv(key, previous)
		} else {
			os.Unsetenv(key)
		}
	}
}

// ReplaceWorkingDirectory replaces the current working directory (os.Getwd()) with a new value.  Returns a function for
// use with defer in order to reset the previous value
//
// defer ReplaceWorkingDirectory(t, "alpha")()
func ReplaceWorkingDirectory(t *testing.T, dir string) func() {
	t.Helper()

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err = os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	return func() { os.Chdir(previous) }
}

// ScratchDir returns a safe scratch directory for tests to modify.
func ScratchDir(t *testing.T, prefix string) string {
	t.Helper()

	tmp, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(err)
	}

	abs, err := filepath.EvalSymlinks(tmp)
	if err != nil {
		t.Fatal(err)
	}

	return abs
}

// FixturePath returns the absolute path to the desired fixture.
func FixturePath(t *testing.T, fixture string) string {
	t.Helper()
	return filepath.Join(FindRoot(t), "fixtures", fixture)
}
