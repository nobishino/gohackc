package compilation_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-multierror"
	"github.com/josharian/txtarfs"
	"github.com/nobishino/gohackc/compilation"
	"golang.org/x/tools/txtar"
)

func TestCompilationEngine(t *testing.T) {
	const OK, NG = false, true
	testcase := []struct {
		testfile    string
		shouldError bool
	}{
		{"bare_class", OK},
		{"bare_class_error", NG},
		{"class_var_dec", OK},
		{"subroutine_dec", OK},
		{"subroutine_dec_2", OK},
		{"subroutine_dec_3", OK},
	}
	for _, tc := range testcase {
		t.Run(tc.testfile, func(t *testing.T) {
			src, want := readTestCase(t, tc.testfile)
			var dst strings.Builder
			e := compilation.New(src, &dst)

			e.CompileClass()

			err := e.Error()
			if !tc.shouldError && err != nil {
				fatalError(t, err)
			}
			if tc.shouldError && err == nil {
				t.Fatal("should return non-nil error but got nil")
			}
			if tc.shouldError {
				return
			}

			got := dst.String()
			if diff := cmp.Diff(want, got); diff != "" {
				t.Error("\n" + got)
				t.Log("DIFF:\n", diff)
			}
		})
	}
}

func readTestCase(t *testing.T, caseName string) (fs.File, string) {
	const testDataDir = "./testdata"
	t.Helper()
	f, err := os.Open(filepath.Join(testDataDir, caseName+".txtar"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { f.Close() })

	txtarData, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	fs := txtarfs.As(txtar.Parse(txtarData))

	src, err := fs.Open("source.jack")
	if err != nil {
		t.Fatal(err, "test case seems broken")
	}

	want, err := fs.Open("expect.xml")
	if err != nil {
		t.Fatal(err, "test case seems broken")
	}

	var buf strings.Builder
	if _, err := io.Copy(&buf, want); err != nil {
		t.Fatal(err)
	}

	return src, buf.String()
}

func fatalError(t *testing.T, err error) {
	t.Helper()
	if merr, ok := err.(*multierror.Error); ok {
		for _, e := range merr.WrappedErrors() {
			t.Logf("%+v\n", e)
		}
		t.FailNow()
	}
}
