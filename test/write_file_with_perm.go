/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/cloudfoundry/libcfbuildpack/test_helper"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WriteFileWithPerm writes a file with specific permissions during testing.
func WriteFileWithPerm(t test_helper.TInterface, filename string, perm os.FileMode, format string, args ...interface{}) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, []byte(fmt.Sprintf(format, args...)), perm); err != nil {
		t.Fatal(err)
	}
}
