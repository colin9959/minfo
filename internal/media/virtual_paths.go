package media

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const virtualISOPrefix = "ISO:"

func ResolveInputPath(ctx context.Context, input string) (string, func(), error) {
	cleaned := strings.TrimSpace(strings.Trim(input, "\""))
	if cleaned == "" {
		return "", func() {}, fmt.Errorf("missing path")
	}

	if isVirtualISOPath(cleaned) {
		return resolveVirtualISOPath(ctx, cleaned)
	}

	cleaned = filepath.Clean(cleaned)
	if _, err := os.Stat(cleaned); err != nil {
		return "", func() {}, fmt.Errorf("path not found: %v", err)
	}
	return cleaned, func() {}, nil
}

func isVirtualISOPath(input string) bool {
	_, _, ok := parseVirtualISOPath(input)
	return ok
}

func parseVirtualISOPath(input string) (string, string, bool) {
	if !strings.HasPrefix(input, virtualISOPrefix) {
		return "", "", false
	}

	rest := strings.TrimPrefix(input, virtualISOPrefix)
	if rest == "" {
		return "", "", false
	}

	var isoPath string
	inner := "/"
	if strings.HasSuffix(rest, "!") {
		isoPath = rest[:len(rest)-1]
	} else {
		bang := strings.Index(rest, "!/")
		if bang < 0 {
			return "", "", false
		}
		isoPath = rest[:bang]
		inner = rest[bang+1:]
	}

	isoPath = filepath.Clean(strings.TrimSpace(isoPath))
	if isoPath == "." || isoPath == "" || !isISOFile(isoPath) {
		return "", "", false
	}

	inner = strings.ReplaceAll(strings.TrimSpace(inner), "\\", "/")
	if inner == "" {
		inner = "/"
	}
	if !strings.HasPrefix(inner, "/") {
		inner = "/" + inner
	}
	inner = path.Clean(inner)
	if inner == "." {
		inner = "/"
	}
	if !strings.HasPrefix(inner, "/") {
		return "", "", false
	}

	return isoPath, inner, true
}

func buildVirtualISOPath(isoPath, inner string, isDir bool) string {
	cleanISO := filepath.Clean(isoPath)
	cleanInner := strings.ReplaceAll(strings.TrimSpace(inner), "\\", "/")
	if cleanInner == "" {
		cleanInner = "/"
	}
	if !strings.HasPrefix(cleanInner, "/") {
		cleanInner = "/" + cleanInner
	}
	cleanInner = path.Clean(cleanInner)
	if cleanInner == "." {
		cleanInner = "/"
	}

	result := virtualISOPrefix + cleanISO + "!"
	if cleanInner != "/" {
		result += cleanInner
	}
	if isDir && !strings.HasSuffix(result, "/") {
		result += "/"
	}
	return result
}

func resolveVirtualISOPath(ctx context.Context, input string) (string, func(), error) {
	isoPath, inner, ok := parseVirtualISOPath(input)
	if !ok {
		return "", func() {}, fmt.Errorf("invalid ISO browser path")
	}
	if _, err := os.Stat(isoPath); err != nil {
		return "", func() {}, fmt.Errorf("path not found: %v", err)
	}

	mountDir, cleanup, err := mountISO(ctx, isoPath)
	if err != nil {
		return "", func() {}, err
	}

	target := mountDir
	if inner != "/" {
		target = filepath.Join(mountDir, filepath.FromSlash(strings.TrimPrefix(inner, "/")))
	}
	target = filepath.Clean(target)
	if !isSubpath(mountDir, target) {
		cleanup()
		return "", func() {}, fmt.Errorf("path is outside mounted ISO")
	}
	if _, err := os.Stat(target); err != nil {
		cleanup()
		return "", func() {}, fmt.Errorf("path not found: %v", err)
	}

	return target, cleanup, nil
}
