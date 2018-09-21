package libjavabuildpack_test

import (
	"testing"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"path/filepath"
	"github.com/cloudfoundry/libjavabuildpack/test"
	"os"
	"github.com/cloudfoundry/libjavabuildpack"
)

func TestExtractTarGz(t *testing.T) {
	spec.Run(t, "ExtractTarGz", testExtractTarGz, spec.Report(report.Terminal{}))
}

func testExtractTarGz(t *testing.T, when spec.G, it spec.S) {
	var outDir string

	it.Before(func() {
		outDir = test.ScratchDir(t, "out")
	})

	it.After(func() {
		os.RemoveAll(outDir)
	})

	it("extracts the archive", func() {
		err := libjavabuildpack.ExtractTarGz(test.FixturePath(t, "test-archive.tar.gz"), outDir, 0)
		if err != nil {
			t.Fatal(err)
		}
		test.BeFileLike(t, filepath.Join(outDir, "fileA.txt"), 0644, "file A\n")
		test.BeFileLike(t, filepath.Join(outDir, "dirA", "fileB.txt"), 0644, "file B\n")
		test.BeFileLike(t, filepath.Join(outDir, "dirA", "fileC.txt"), 0644, "file C\n")
	})

	it("skips stripped components", func() {
		err := libjavabuildpack.ExtractTarGz(test.FixturePath(t, "test-archive.tar.gz"), outDir, 1)
		if err != nil {
			t.Fatal(err)
		}
		test.BeFileLike(t, filepath.Join(outDir, "fileB.txt"), 0644, "file B\n")
		test.BeFileLike(t, filepath.Join(outDir, "fileC.txt"), 0644, "file C\n")
	})
}
