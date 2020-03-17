/*
 * Copyright 2018-2020 the original author or authors.
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

package helper

import (
	"compress/bzip2"
	"os"
)

// ExtractTarBz2 extracts source BZIP'd TAR file to a destination directory.  An arbitrary number of top-level directory
// components can be stripped from each path.
func ExtractTarBz2(source string, destination string, stripComponents int) error {
	f, err := os.Open(source)
	if err != nil {
		return err
	}
	defer f.Close()

	return handleTar(bzip2.NewReader(f), destination, stripComponents)
}