package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func resolveBin(envKey, fallback string) (string, error) {
	bin := getenv(envKey, fallback)
	if _, err := exec.LookPath(bin); err != nil {
		return "", fmt.Errorf("%s not found; set %s or add to PATH", bin, envKey)
	}
	return bin, nil
}

func runCommand(ctx context.Context, bin string, args ...string) (string, string, error) {
	cmd := exec.Command(bin, args...)
	setCommandProcessGroup(cmd)

	stdoutFile, err := os.CreateTemp("", "minfo-stdout-*")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(stdoutFile.Name())
	defer stdoutFile.Close()

	stderrFile, err := os.CreateTemp("", "minfo-stderr-*")
	if err != nil {
		return "", "", err
	}
	defer os.Remove(stderrFile.Name())
	defer stderrFile.Close()

	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	if err := cmd.Start(); err != nil {
		return "", "", err
	}

	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
	}()

	var waitErr error
	select {
	case waitErr = <-waitCh:
	case <-ctx.Done():
		killCommandProcessGroup(cmd)
		waitErr = ctx.Err()
		<-waitCh
	}

	stdoutData, _ := os.ReadFile(stdoutFile.Name())
	stderrData, _ := os.ReadFile(stderrFile.Name())
	return string(stdoutData), string(stderrData), waitErr
}
