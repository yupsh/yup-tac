package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestRun(t *testing.T) {
	cases := []struct {
		files      map[string]string
		name       string
		version    string
		stdin      string
		wantOut    string
		wantErrSub string
		args       []string
		wantCode   int
	}{
		{
			name:    "reverse stdin lines",
			args:    []string{"tac"},
			stdin:   "alpha\nbeta\ngamma\n",
			wantOut: "gamma\nbeta\nalpha\n",
		},
		{
			name:    "separator splits records",
			args:    []string{"tac", "-s", ":"},
			stdin:   "a:b:c\n",
			wantOut: "c\nb\na\n",
		},
		{
			name:    "file source",
			args:    []string{"tac", "/in.txt"},
			files:   map[string]string{"/in.txt": "one\ntwo\nthree\n"},
			wantOut: "three\ntwo\none\n",
		},
		{
			name:    "version flag reports injected version",
			version: "1.2.3",
			args:    []string{"tac", "--version"},
			wantOut: "tac version 1.2.3\n",
		},
		{
			name:       "unknown flag errors",
			args:       []string{"tac", "--nope"},
			wantCode:   1,
			wantErrSub: "tac:",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			for path, content := range tc.files {
				if err := afero.WriteFile(fs, path, []byte(content), 0o644); err != nil {
					t.Fatalf("write fixture %s: %v", path, err)
				}
			}

			var out, errOut bytes.Buffer
			code := run(tc.version, tc.args, strings.NewReader(tc.stdin), &out, &errOut, fs)

			if code != tc.wantCode {
				t.Fatalf("exit code = %d, want %d (stderr=%q)", code, tc.wantCode, errOut.String())
			}
			if tc.wantErrSub == "" && out.String() != tc.wantOut {
				t.Fatalf("stdout = %q, want %q", out.String(), tc.wantOut)
			}
			if tc.wantErrSub != "" && !strings.Contains(errOut.String(), tc.wantErrSub) {
				t.Fatalf("stderr = %q, want substring %q", errOut.String(), tc.wantErrSub)
			}
		})
	}
}

func Test_main(t *testing.T) {
	origExit, origRun := osExit, runCLI
	t.Cleanup(func() { osExit, runCLI = origExit, origRun })

	gotCode := -1
	osExit = func(code int) { gotCode = code }
	runCLI = func(string, []string, io.Reader, io.Writer, io.Writer, afero.Fs) int { return 7 }

	main()

	if gotCode != 7 {
		t.Fatalf("main propagated exit code %d, want 7", gotCode)
	}
}
