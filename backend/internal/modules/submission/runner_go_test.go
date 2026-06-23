package submission

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// fakeExecCommandContext adalah fungsi yang menggantikan exec.CommandContext asli.
// Ia membelokkan eksekusi ke proses helper di dalam tes itu sendiri.
func fakeExecCommandContext(ctx context.Context, command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", "WANT_HELPER", command}
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	return cmd
}

// TestHelperProcess bukan pengujian sungguhan. Ini adalah trik "Aktor Pengganti".
// Saat exec.CommandContext memanggil file ini, proses pembantu ini dijalankan
// dan memanipulasi stdout, stderr, atau exit status sesuai kebutuhan tes.
func TestHelperProcess(t *testing.T) {
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	
	if len(args) == 0 || args[0] != "WANT_HELPER" {
		return
	}

	cmd, cmdArgs := args[1], args[2:]

	switch cmd {
	case "go":
		// Mock untuk "go build"
		if len(cmdArgs) > 0 && cmdArgs[0] == "build" {
			content, err := os.ReadFile("main.go")
			if err != nil {
				fmt.Fprintf(os.Stderr, "TestHelper Error reading main.go: %v\n", err)
				os.Exit(1)
			}
			
			if strings.Contains(string(content), "syntax_error") {
				fmt.Fprint(os.Stderr, "Compilation Error: syntax error")
				os.Exit(1)
			} else {
			    fmt.Fprintf(os.Stderr, "TestHelper: main.go content does not contain syntax_error. Content:\n%s\n", string(content))
			}
			// Sukses
			os.Exit(0)
		}
	case "docker":
		// Mock untuk "docker run"
		if len(cmdArgs) > 0 && cmdArgs[0] == "run" {
			// Membaca Stdin (Input pengguna)
			inputBytes, _ := io.ReadAll(os.Stdin)
			input := string(inputBytes)

			// Cek tipe skenario
			for _, a := range cmdArgs {
				if strings.Contains(a, "timeout_test") {
					// Simulasi timeout dengan delay (meskipun context timeout yang mematikannya)
					// Atau kita bisa paksa return code
					// Biarkan ia hang agar context deadline terpicu, tapi dalam unit test
					// kita ingin cepat. Jadi kita pura-pura lama tapi sebenarnya OS kill.
					// Opsi lebih baik: tidak perlu sleep lama, context di test dibuat sangat kecil.
				}
			}

			// Simulasi Output sukses
			if strings.TrimSpace(input) == "Dunia" {
				fmt.Fprint(os.Stdout, "Halo Dunia")
			} else {
				fmt.Fprint(os.Stdout, "Output untuk: "+input)
			}
			os.Exit(0)
		}
	}

	// Default
	os.Exit(0)
}

func TestGoRunner_ValidCode(t *testing.T) {
	// 1. Timpa dengan mock
	execCommandContext = fakeExecCommandContext
	defer func() { execCommandContext = exec.CommandContext }() // kembalikan ke aslinya nanti

	code := `
package main
import "fmt"
func main() {
	var input string
	fmt.Scanln(&input)
	fmt.Print("Halo " + input)
}
`
	runner := &GoRunner{}
	compileDir, err := runner.Compile(code)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	defer runner.Cleanup(compileDir)

	result := runner.RunTestCase(compileDir, "Dunia\n", 15000)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}
	if result.Output != "Halo Dunia" {
		t.Errorf("Expected Output 'Halo Dunia', got '%s'", result.Output)
	}
}

func TestGoRunner_SyntaxError(t *testing.T) {
	execCommandContext = fakeExecCommandContext
	defer func() { execCommandContext = exec.CommandContext }()

	// Kita menyisipkan komentar aneh agar TestHelperProcess tahu ini adalah tes error
	code := `
package main
// syntax_error
func main() {
	fmt.Print("Missing bracket"
}
`
	runner := &GoRunner{}
	compileDir, err := runner.Compile(code)
	runner.Cleanup(compileDir)

	if err == nil {
		t.Errorf("Expected compilation error, got nil")
	} else if !strings.Contains(err.Error(), "Compilation Error") {
		t.Errorf("Expected compilation error to mention 'Compilation Error', got: %v", err)
	}
}

func TestGoRunner_TimeLimitExceeded(t *testing.T) {
	execCommandContext = func(ctx context.Context, command string, args ...string) *exec.Cmd {
		// Mock langsung khusus untuk tes ini
		cmd := exec.CommandContext(ctx, "sleep", "1") // paksa delay
		return cmd
	}
	defer func() { execCommandContext = exec.CommandContext }()

	runner := &GoRunner{}
	// Kita by-pass Compile dengan membuat direktori bohongan
	compileDir := "/tmp/mock_compile_dir"

	// Kita berikan batas waktu sangat singkat (1ms) agar cepat Time Limit
	result := runner.RunTestCase(compileDir, "", 1)

	if result.Error != "Time Limit Exceeded" {
		t.Errorf("Expected Error 'Time Limit Exceeded', got '%s'", result.Error)
	}
}
