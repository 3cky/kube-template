// Copyright © 2015 Victor Antonovich <victor@antonovich.me>
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

package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/golang/glog"
)

// Normalize given path: evaluate symlinks, convert to absolute and clean
func NormPath(path string) (string, error) {
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	}
	p, err = filepath.Abs(p)
	if err != nil {
		return path, err
	}
	return filepath.Clean(p), nil
}

// Check string is present in given list
func IsPresent(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}

// Execute command using system shell with timeout
func Execute(command string, timeout time.Duration) error {
	// Set shell and command execution flag
	shell, flag := "/bin/sh", "-c"
	if runtime.GOOS == "windows" {
		shell, flag = "cmd", "/C"
	}

	// Create command execution and get stdout/stderr
	cmd := exec.Command(shell, flag, command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer CloseQuietly(stdout)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer CloseQuietly(stderr)

	// Start command execution
	if err := cmd.Start(); err != nil {
		return err
	}

	// Create command result channel
	result := make(chan error, 1)
	defer close(result)
	go func() {
		result <- cmd.Wait()
	}()

	// Log stdout/stderr
	outScanner := bufio.NewScanner(stdout)
	go func() {
		for outScanner.Scan() {
			glog.V(4).Infof("STDOUT: %s", outScanner.Text())
		}
		if err := outScanner.Err(); err != nil {
			glog.Errorf("STDOUT: error: %v", err)
		}
	}()
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for errScanner.Scan() {
			glog.V(4).Infof("STDERR: %s", errScanner.Text())
		}
		if err := errScanner.Err(); err != nil {
			glog.Errorf("STDERR: error: %v", err)
		}
	}()

	// Wait for result or timeout
	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				glog.Errorf("timeout: %q, not killed: %v", command, err)
			} else {
				glog.Warningf("timeout: %q, killed", command)
			}
		} else {
			glog.Warningf("timeout: %q, nothing to kill", command)
		}
		return fmt.Errorf("timeout: %q", command)
	case err := <-result:
		return err
	}
}

// Close given closer without error checking
func CloseQuietly(closer io.Closer) {
	_ = closer.Close()
}

// Unlink given file without error checking
func UnlinkQuietly(path string) {
	_ = syscall.Unlink(path)
}
