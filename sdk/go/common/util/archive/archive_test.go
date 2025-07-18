// Copyright 2016-2018, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIgnoreSimple(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, ".",
		fileContents{name: ".gitignore", contents: []byte("node_modules/pulumi/"), shouldRetain: true},
		fileContents{name: "included.txt", shouldRetain: true},
		fileContents{name: "node_modules/included.txt", shouldRetain: true},
		fileContents{name: "node_modules/pulumi/excluded.txt", shouldRetain: false},
		fileContents{name: "node_modules/pulumi/excluded/excluded.txt", shouldRetain: false})
}

func TestIgnoreNegate(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("Skipped on Windows: TODO[pulumi/pulumi#8648] handle Windows paths in test logic")
	}

	doArchiveTest(t, ".",
		fileContents{name: ".gitignore", contents: []byte("/*\n!/foo\n/foo/*\n!/foo/bar"), shouldRetain: false},
		fileContents{name: "excluded.txt", shouldRetain: false},
		fileContents{name: "foo/excluded.txt", shouldRetain: false},
		fileContents{name: "foo/baz/exlcuded.txt", shouldRetain: false},
		fileContents{name: "foo/bar/included.txt", shouldRetain: true})
}

func TestNested(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, ".",
		fileContents{name: ".gitignore", contents: []byte("node_modules/pulumi/"), shouldRetain: true},
		fileContents{name: "node_modules/.gitignore", contents: []byte("@pulumi/"), shouldRetain: true},
		fileContents{name: "included.txt", shouldRetain: true},
		fileContents{name: "node_modules/included.txt", shouldRetain: true},
		fileContents{name: "node_modules/pulumi/excluded.txt", shouldRetain: false},
		fileContents{name: "node_modules/@pulumi/pulumi-cloud/excluded.txt", shouldRetain: false})
}

func TestTypicalPythonPolicyPackDir(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, ".",
		fileContents{name: "__main__.py", shouldRetain: true},
		fileContents{name: ".gitignore", contents: []byte("*.pyc\nvenv/\n"), shouldRetain: true},
		fileContents{name: "PulumiPolicy.yaml", shouldRetain: true},
		fileContents{name: "requirements.txt", shouldRetain: true},
		fileContents{name: "venv/bin/activate", shouldRetain: false},
		fileContents{name: "venv/bin/pip", shouldRetain: false},
		fileContents{name: "venv/bin/python", shouldRetain: false},
		fileContents{name: "__pycache__/__main__.cpython-37.pyc", shouldRetain: false})
}

func TestIgnoreContentOfDotGit(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, ".",
		fileContents{name: ".git/HEAD", shouldRetain: false},
		fileContents{name: ".git/objects/00/02ae827766d77ee9e2082fee9adeaae90aff65", shouldRetain: false},
		fileContents{name: "__main__.py", shouldRetain: true},
		fileContents{name: "PulumiPolicy.yaml", shouldRetain: true},
		fileContents{name: "requirements.txt", shouldRetain: true})
}

func TestNestedPath(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, "pkg/",
		fileContents{name: "excluded.txt", shouldRetain: false},
		fileContents{name: "pkg/.gitignore", contents: []byte("node_modules/pulumi/"), shouldRetain: true},
		fileContents{name: "pkg/node_modules/included.txt", shouldRetain: true},
		fileContents{name: "pkg/node_modules/pulumi/excluded.txt", shouldRetain: false},
		fileContents{name: "pkg/node_modules/pulumi/excluded/excluded.txt", shouldRetain: false})
}

func TestIgnoreNestedGitignore(t *testing.T) {
	t.Parallel()

	doArchiveTest(t, "pkg/",
		fileContents{name: ".gitignore", contents: []byte("*.ts"), shouldRetain: false},
		fileContents{name: "excluded.txt", shouldRetain: false},
		fileContents{name: "pkg/.gitignore", contents: []byte("node_modules/pulumi/"), shouldRetain: true},
		fileContents{name: "pkg/node_modules/excluded.ts", shouldRetain: false},
		fileContents{name: "pkg/node_modules/included.txt", shouldRetain: true},
		fileContents{name: "pkg/node_modules/pulumi/excluded.txt", shouldRetain: false},
		fileContents{name: "pkg/node_modules/pulumi/excluded/excluded.txt", shouldRetain: false})
}

func doArchiveTest(t *testing.T, path string, files ...fileContents) {
	doTest := func(prefixPathInsideTar, path string) {
		tarball, err := archiveContents(t, prefixPathInsideTar, path, files...)
		require.NoError(t, err)

		tarReader := bytes.NewReader(tarball)
		gzr, err := gzip.NewReader(tarReader)
		require.NoError(t, err)
		r := tar.NewReader(gzr)

		checkFiles(t, prefixPathInsideTar, path, files, r)
	}
	for _, prefix := range []string{"", "package"} {
		doTest(prefix, path)
	}
}

func archiveContents(t *testing.T, prefixPathInsideTar, path string, files ...fileContents) ([]byte, error) {
	dir := t.TempDir()

	for _, file := range files {
		name := file.name
		if os.PathSeparator != '/' {
			name = strings.ReplaceAll(name, "/", string(os.PathSeparator))
		}

		err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(filepath.Join(dir, name), file.contents, 0o600)
		if err != nil {
			return nil, err
		}
	}

	return TGZ(filepath.Join(dir, path), prefixPathInsideTar, true /*useDefaultExcludes*/)
}

func checkFiles(t *testing.T, prefixPathInsideTar, path string, expected []fileContents, r *tar.Reader) {
	var expectedFiles []string
	var actualFiles []string

	for _, f := range expected {
		if f.shouldRetain {
			name := f.name
			if path != "." {
				name = strings.Replace(name, path, "", 1)
			}
			if prefixPathInsideTar != "" {
				// Joining with '/' rather than platform-specific `filepath.Join` because we expect
				// the name in the tar to be using '/'.
				name = fmt.Sprintf("%s/%s", prefixPathInsideTar, name)
			}
			expectedFiles = append(expectedFiles, name)
		}
	}

	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		// Ignore anything other than regular files (e.g. directories) since we only care
		// that the files themselves are correct.
		if header.Typeflag != tar.TypeReg {
			continue
		}

		actualFiles = append(actualFiles, header.Name)
	}

	sort.Strings(expectedFiles)
	sort.Strings(actualFiles)

	assert.Equal(t, expectedFiles, actualFiles)
}

type fileContents struct {
	name         string
	contents     []byte
	shouldRetain bool
}
