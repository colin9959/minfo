package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func mountISO(ctx context.Context, isoPath string) (string, func(), error) {
	mountBin, err := resolveBin("MOUNT_BIN", "mount")
	if err != nil {
		return "", noop, err
	}
	umountBin, err := resolveBin("UMOUNT_BIN", "umount")
	if err != nil {
		return "", noop, err
	}

	mountDir, err := os.MkdirTemp("", "minfo-iso-mount-*")
	if err != nil {
		return "", noop, err
	}

	mountCtx, cancel := context.WithTimeout(ctx, mountTimeout)
	defer cancel()

	modErr := loadUDFModule(mountCtx)
	_, stderr, err := runCommand(mountCtx, mountBin, "-o", "loop,ro", isoPath, mountDir)
	if err != nil {
		msg := strings.TrimSpace(stderr)
		if msg == "" {
			msg = err.Error()
		}
		if isUnknownUDFMountError(msg) {
			if retryModErr := loadUDFModule(mountCtx); retryModErr == nil {
				_, retryStderr, retryErr := runCommand(mountCtx, mountBin, "-o", "loop,ro", isoPath, mountDir)
				if retryErr == nil {
					return mountDir, buildMountCleanup(mountDir, umountBin), nil
				}
				retryMsg := strings.TrimSpace(retryStderr)
				if retryMsg == "" {
					retryMsg = retryErr.Error()
				}
				_ = os.RemoveAll(mountDir)
				return "", noop, fmt.Errorf("mount iso failed after modprobe udf: %s", retryMsg)
			}
		}
		_ = os.RemoveAll(mountDir)
		return "", noop, fmt.Errorf("mount iso failed: %s", explainISOmountError(msg, modErr))
	}

	return mountDir, buildMountCleanup(mountDir, umountBin), nil
}

func explainISOmountError(message string, modErr error) string {
	if isUnknownUDFMountError(message) {
		if modErr != nil {
			return message + "; auto `modprobe udf` failed: " + modErr.Error() + ". Ensure host supports udf and mount `/lib/modules:/lib/modules:ro` into container"
		}
		return message + "; attempted auto `modprobe udf`, please check host kernel module availability"
	}
	return message
}

func isUnknownUDFMountError(message string) bool {
	lower := strings.ToLower(message)
	return strings.Contains(lower, "unknown filesystem type 'udf'") || strings.Contains(lower, "unknown filesystem type \"udf\"")
}

func loadUDFModule(ctx context.Context) error {
	modprobeBin, err := resolveBin("MODPROBE_BIN", "modprobe")
	if err != nil {
		return err
	}
	_, stderr, err := runCommand(ctx, modprobeBin, "udf")
	if err != nil {
		msg := strings.TrimSpace(stderr)
		if msg == "" {
			msg = err.Error()
		}
		return fmt.Errorf("modprobe udf failed: %s", msg)
	}
	return nil
}

func buildMountCleanup(mountDir, umountBin string) func() {
	return func() {
		umountCtx, cancel := context.WithTimeout(context.Background(), umountTimeout)
		defer cancel()
		if _, _, err := runCommand(umountCtx, umountBin, mountDir); err != nil {
			_, _, _ = runCommand(umountCtx, umountBin, "-l", mountDir)
		}
		_ = os.RemoveAll(mountDir)
	}
}
