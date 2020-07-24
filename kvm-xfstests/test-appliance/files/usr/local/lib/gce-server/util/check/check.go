/*
Package check executes external commands, performs File I/Os and issues
OS related commands.

It also checks for errors and writes messages into logger.
*/
package check

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// RootDir points to the root of go server source code
// The compiled go executables are located in GOPATH/bin
const RootDir = "/usr/local/lib/gce-server"

// EmptyEnv provides a placeholder for default exec environment.
var EmptyEnv = map[string]string{}

// Run executes an external command and checks the return status.
// Returns true on success and false otherwise.
func Run(cmd *exec.Cmd, workDir string, env map[string]string, stdout io.Writer, stderr io.Writer) error {
	cmd.Dir = workDir
	cmd.Env = parseEnv(env)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	return err
}

// Output executes an external command, checks the return status, and
// returns the command stdout.
func Output(cmd *exec.Cmd, workDir string, env map[string]string, stderr io.Writer) (string, error) {
	cmd.Dir = workDir
	cmd.Env = parseEnv(env)
	cmd.Stderr = stderr
	out, err := cmd.Output()
	return string(out), err
}

// CombinedOutput executes an external command, checks the return status, and
// returns the combined stdout and stderr.
func CombinedOutput(cmd *exec.Cmd, workDir string, env map[string]string) (string, error) {
	cmd.Dir = workDir
	cmd.Env = parseEnv(env)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// parseEnv adds user specified environment to os.Environ.
func parseEnv(env map[string]string) []string {
	newEnv := os.Environ()
	for key, value := range env {
		newEnv = append(newEnv, key+"="+value)
	}
	return newEnv
}

// CreateDir creates a directory recursively with default permissions.
func CreateDir(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}

// RemoveDir removes a directory and all contents in it.
// Do nothing if the target path doesn't exist.
func RemoveDir(path string) error {
	err := os.RemoveAll(path)
	return err
}

// FileExists returns true if a file exists.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err == nil && !info.IsDir() {
		return true
	}
	return false
}

// DirExists returns true is a directory exists.
func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}

// ReadLines read a whole file into a slice of strings split by newlines.
// Removes '\n' and empty lines
func ReadLines(filename string) ([]string, error) {
	lines := []string{}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines, err
	}
	lines = strings.Split(string(content), "\n")
	nonEmptyLines := lines[:0]
	for i, line := range lines {
		if line != "" {
			nonEmptyLines = append(nonEmptyLines, lines[i:i+1]...)
		}
	}
	return nonEmptyLines, nil
}

// CopyFile copies the content of file src to file dst
// Overwrites dst if it already exists
func CopyFile(dst string, src string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

// Panic checks an error and log a panic entry with given msg
func Panic(err error, log *logrus.Entry, msg string) {
	if msg == "" {
		msg = "Something bad happended"
	}
	if err != nil {
		log.WithError(err).Panic(msg)
	}
}

// NoError checks an error and log a error entry with given msg
// return true if error is nil
func NoError(err error, log *logrus.Entry, msg string) bool {
	if msg == "" {
		msg = "Something bad happended"
	}
	if err != nil {
		log.WithError(err).Error(msg)
		return false
	}
	return true
}