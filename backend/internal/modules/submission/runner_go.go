package submission

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type GoRunner struct{}

func (r *GoRunner) Compile(code string) (string, error) {
	// 1. Membuat direktori sementara
	tmpDir, err := os.MkdirTemp("", "leetcode_go_*")
	if err != nil {
		return "", fmt.Errorf("Gagal membuat direktori sementara: %v", err)
	}

	// 2. Menulis kode pengguna ke file main.go
	codePath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(codePath, []byte(code), 0600)
	if err != nil {
		return tmpDir, fmt.Errorf("Gagal menulis file kode: %v", err)
	}

	// === FASE 1: KOMPILASI LOKAL (HOST) ===
	compileCtx, compileCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer compileCancel()

	compileCmd := exec.CommandContext(compileCtx, "go", "build", "-o", "main_exec", "main.go")
	compileCmd.Dir = tmpDir
	// Menyuntikkan CGO_ENABLED=0 agar file biner bisa berjalan di alpine tanpa error library glibc
	compileCmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")

	var compileErrBuf bytes.Buffer
	compileCmd.Stderr = &compileErrBuf

	err = compileCmd.Run()
	if err != nil {
		if compileCtx.Err() == context.DeadlineExceeded {
			return tmpDir, errors.New("Compilation Time Limit Exceeded")
		}
		errMsg := compileErrBuf.String()
		if errMsg == "" {
			errMsg = err.Error()
		}
		return tmpDir, errors.New("Compilation Error:\n" + errMsg)
	}

	return tmpDir, nil
}

func (r *GoRunner) RunTestCase(compileDir string, input string, timeLimitMs int) TestCaseResult {
	start := time.Now()
	result := TestCaseResult{}

	// === FASE 2: EKSEKUSI DOCKER (CONTAINER) ===
	execCtx, execCancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMs)*time.Millisecond)
	defer execCancel()

	// Mengeksekusi file biner di dalam container alpine super ringan
	// Menggunakan flag -i (interactive) agar container bisa menerima input dari Stdin
	/* #nosec G204 -- compileDir is safely generated via filepath.Join(os.TempDir(), ...) */
	runCmd := exec.CommandContext(execCtx, "docker", "run", "-i", "--rm",
		"--network", "none",
		"--memory", "128m",
		"--cpus", "1.0",
		"--pids-limit", "64",
		"-v", fmt.Sprintf("%s:/app", compileDir),
		"-w", "/app",
		"alpine:latest",
		"./main_exec",
	)

	// Menyuntikkan input test case ke dalam standard input (Stdin) container
	runCmd.Stdin = strings.NewReader(input)

	var runOutBuf, runErrBuf bytes.Buffer
	runCmd.Stdout = &runOutBuf
	runCmd.Stderr = &runErrBuf

	err := runCmd.Run()
	
	result.TimeTakenMs = time.Since(start).Milliseconds()
	result.Output = strings.TrimSpace(runOutBuf.String())
	
	if execCtx.Err() == context.DeadlineExceeded {
		result.Error = "Time Limit Exceeded"
		return result
	}

	if err != nil {
		result.Error = runErrBuf.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
	}

	return result
}

func (r *GoRunner) Cleanup(compileDir string) {
	if compileDir != "" {
		_ = os.RemoveAll(compileDir)
	}
}
