//go:build !windows

package script_test

import (
	"testing"

	"github.com/bitfield/script"
	"github.com/google/go-cmp/cmp"
)

func TestExecForEach_HandlesLongLines(t *testing.T) {
	t.Parallel()
	got, err := script.Echo(longLine).ExecForEach(`echo "{{.}}"`).String()
	if err != nil {
		t.Fatal(err)
	}
	if longLine != got {
		t.Error(cmp.Diff(longLine, got))
	}
}

func TestExecRunsShWithEchoHelloAndGetsOutputHello(t *testing.T) {
	t.Parallel()
	p := script.Exec("sh -c 'echo hello'")
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "hello\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecvWithDynamicArgs(t *testing.T) {
	t.Parallel()
	args := []string{"-c"}
	args = append(args, "echo hello")
	p := script.Execv("sh", args)
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "hello\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecvRunsShWithEchoHelloAndGetsOutputHello(t *testing.T) {
	t.Parallel()
	p := script.Execv("sh", []string{"-c", "echo hello"})
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "hello\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecRunsShWithinShWithEchoInceptionAndGetsOutputInception(t *testing.T) {
	t.Parallel()
	p := script.Exec("sh -c 'sh -c \"echo inception\"'")
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "inception\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecvRunsShWithinShWithEchoInceptionAndGetsOutputInception(t *testing.T) {
	t.Parallel()
	p := script.Execv("sh", []string{"-c", "sh -c \"echo inception\""})
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "inception\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecErrorsRunningShellCommandWithUnterminatedStringArgument(t *testing.T) {
	t.Parallel()
	p := script.Exec("sh -c 'echo oh no")
	p.Wait()
	if p.Error() == nil {
		t.Error("want error running 'sh' command line containing unterminated string")
	}
}

func TestExecvErrorsRunningShellCommandWithUnterminatedStringArgument(t *testing.T) {
	t.Parallel()
	p := script.Execv("sh", []string{"-c", "'echo oh no"})
	p.Wait()
	if p.Error() == nil {
		t.Error("want error running 'sh' command line containing unterminated string")
	}
}

func TestExecForEach_RunsEchoWithABCAndGetsOutputABC(t *testing.T) {
	t.Parallel()
	p := script.Echo("a\nb\nc\n").ExecForEach("echo {{.}}")
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "a\nb\nc\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecForEach_CorrectlyEvaluatesTemplateContainingIfStatement(t *testing.T) {
	t.Parallel()
	p := script.Echo("a").ExecForEach("echo {{if .}}it worked!{{end}}")
	if p.Error() != nil {
		t.Fatal(p.Error())
	}
	want := "it worked!\n"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func TestExecPipesDataToExternalCommandAndGetsExpectedOutput(t *testing.T) {
	t.Parallel()
	p := script.File("testdata/hello.txt").Exec("cat")
	want := "hello world"
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if want != got {
		t.Error(cmp.Diff(want, got))
	}
}

func ExampleExec_ok() {
	script.Exec("echo Hello, world!").Stdout()
	// Output:
	// Hello, world!
}

func ExampleExecv_ok() {
	script.Execv("echo", []string{"Hello, world!"}).Stdout()
	// Output:
	// Hello, world!
}

func ExampleFindFiles() {
	script.FindFiles("testdata/multiple_files_with_subdirectory").Stdout()
	// Output:
	// testdata/multiple_files_with_subdirectory/1.txt
	// testdata/multiple_files_with_subdirectory/2.txt
	// testdata/multiple_files_with_subdirectory/3.tar.zip
	// testdata/multiple_files_with_subdirectory/dir/.hidden
	// testdata/multiple_files_with_subdirectory/dir/1.txt
	// testdata/multiple_files_with_subdirectory/dir/2.txt
}

func ExampleIfExists_exec() {
	script.IfExists("./testdata/hello.txt").Exec("echo hello").Stdout()
	// Output:
	// hello
}

func ExampleIfExists_execv() {
	script.IfExists("./testdata/hello.txt").Execv("echo", []string{"hello"}).Stdout()
	// Output:
	// hello
}

func ExampleIfExists_noExec() {
	script.IfExists("doesntexist").Execv("echo", []string{"hello"}).Stdout()
	// Output:
	//
}

func ExampleListFiles() {
	script.ListFiles("testdata/multiple_files_with_subdirectory").Stdout()
	// Output:
	// testdata/multiple_files_with_subdirectory/1.txt
	// testdata/multiple_files_with_subdirectory/2.txt
	// testdata/multiple_files_with_subdirectory/3.tar.zip
	// testdata/multiple_files_with_subdirectory/dir
}

func ExamplePipe_Basename() {
	input := []string{
		"",
		"/",
		"/root",
		"/tmp/example.php",
		"/var/tmp/",
		"./src/filters",
		"C:/Program Files",
	}
	script.Slice(input).Basename().Stdout()
	// Output:
	// .
	// /
	// root
	// example.php
	// tmp
	// filters
	// Program Files
}

func ExamplePipe_Dirname() {
	input := []string{
		"",
		"/",
		"/root",
		"/tmp/example.php",
		"/var/tmp/",
		"./src/filters",
		"C:/Program Files",
	}
	script.Slice(input).Dirname().Stdout()
	// Output:
	// .
	// /
	// /
	// /tmp
	// /var
	// ./src
	// C:
}

func ExamplePipe_Exec() {
	script.Echo("Hello, world!").Exec("tr a-z A-Z").Stdout()
	// Output:
	// HELLO, WORLD!
}

func ExamplePipe_Execv() {
	script.Echo("Hello, world!").Execv("tr", []string{"a-z", "A-Z"}).Stdout()
	// Output:
	// HELLO, WORLD!
}

func ExamplePipe_ExecForEach() {
	script.Echo("a\nb\nc\n").ExecForEach("echo {{.}}").Stdout()
	// Output:
	// a
	// b
	// c
}
