package submission

import (
	"strings"
	"testing"
)

func TestGoRunner_ValidCode(t *testing.T) {
	code := `
package main
import (
	"fmt"
)
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

	// Menguji dengan Test Case Input "Dunia"
	result := runner.RunTestCase(compileDir, "Dunia\n", 15000)

	if result.Error != "" {
		t.Errorf("Expected no error, got: %s", result.Error)
	}

	if result.Output != "Halo Dunia" {
		t.Errorf("Expected Output 'Halo Dunia', got '%s'", result.Output)
	}
}

func TestGoRunner_SyntaxError(t *testing.T) {
	code := `
package main
import "fmt"
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
	code := `
package main
import "time"
func main() {
	time.Sleep(3 * time.Second)
}
`
	runner := &GoRunner{}
	compileDir, err := runner.Compile(code)
	if err != nil {
		t.Fatalf("Compile error: %v", err)
	}
	defer runner.Cleanup(compileDir)

	result := runner.RunTestCase(compileDir, "", 1000)

	if result.Error != "Time Limit Exceeded" {
		t.Errorf("Expected Error 'Time Limit Exceeded', got '%s'", result.Error)
	}
}
