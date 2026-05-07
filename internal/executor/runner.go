package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// RunResult holds the output of a command execution
type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Duration time.Duration
	Err      error
}

// Runner executes shell commands with optional proxy and timeout
type Runner struct {
	Proxy   string
	Timeout time.Duration
	LogFunc func(string) // callback for real-time log output
}

// NewRunner creates a runner with defaults
func NewRunner(proxy string) *Runner {
	return &Runner{
		Proxy:   proxy,
		Timeout: 5 * time.Minute,
	}
}

// Run executes a command string via bash
func (r *Runner) Run(command string) RunResult {
	return r.RunCtx(context.Background(), command)
}

// RunCtx executes a command with context
func (r *Runner) RunCtx(ctx context.Context, command string) RunResult {
	start := time.Now()

	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	// Inject proxy env vars if configured
	cmd.Env = r.buildEnv()

	var stdout, stderr bytes.Buffer

	if r.LogFunc != nil {
		// Stream output in real-time
		stdoutPipe, _ := cmd.StdoutPipe()
		stderrPipe, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
			return RunResult{
				Err:      err,
				ExitCode: -1,
				Duration: time.Since(start),
			}
		}

		// Read stdout
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdoutPipe.Read(buf)
				if n > 0 {
					text := string(buf[:n])
					stdout.WriteString(text)
					r.LogFunc(text)
				}
				if err != nil {
					break
				}
			}
		}()

		// Read stderr
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderrPipe.Read(buf)
				if n > 0 {
					text := string(buf[:n])
					stderr.WriteString(text)
					r.LogFunc(text)
				}
				if err != nil {
					break
				}
			}
		}()

		err := cmd.Wait()
		return RunResult{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: cmd.ProcessState.ExitCode(),
			Duration: time.Since(start),
			Err:      err,
		}
	}

	// Non-streaming mode
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return RunResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Duration: time.Since(start),
		Err:      err,
	}
}

// Check runs a check command and returns true if it succeeds (exit 0)
func (r *Runner) Check(command string) bool {
	if command == "" {
		return false
	}
	result := r.Run(command)
	return result.ExitCode == 0
}

// RunScript runs a multi-line script
func (r *Runner) RunScript(script string) RunResult {
	// Write script to temp file
	f, err := os.CreateTemp("", "dotsetup-*.sh")
	if err != nil {
		return RunResult{Err: err, ExitCode: -1}
	}
	defer os.Remove(f.Name())

	if _, err := io.WriteString(f, script); err != nil {
		f.Close()
		return RunResult{Err: err, ExitCode: -1}
	}
	f.Close()

	return r.Run(fmt.Sprintf("bash %s", f.Name()))
}

func (r *Runner) buildEnv() []string {
	env := os.Environ()
	if r.Proxy == "" {
		return env
	}

	// Add proxy environment variables
	proxyVars := []string{
		fmt.Sprintf("http_proxy=%s", r.Proxy),
		fmt.Sprintf("https_proxy=%s", r.Proxy),
		fmt.Sprintf("HTTP_PROXY=%s", r.Proxy),
		fmt.Sprintf("HTTPS_PROXY=%s", r.Proxy),
		fmt.Sprintf("ALL_PROXY=%s", r.Proxy),
	}

	// Replace existing proxy vars or append
	result := make([]string, 0, len(env)+len(proxyVars))
	for _, e := range env {
		key := strings.SplitN(e, "=", 2)[0]
		lower := strings.ToLower(key)
		if lower == "http_proxy" || lower == "https_proxy" || lower == "all_proxy" {
			continue // skip, we'll add our own
		}
		result = append(result, e)
	}
	result = append(result, proxyVars...)

	return result
}
