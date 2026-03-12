package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func mediaRoots() []string {
	if roots := detectMountedRoots(); len(roots) > 0 {
		return roots
	}
	return []string{defaultRoot}
}

func detectMountedRoots() []string {
	content, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return nil
	}

	ignoredFS := map[string]struct{}{
		"overlay":    {},
		"proc":       {},
		"sysfs":      {},
		"tmpfs":      {},
		"devpts":     {},
		"mqueue":     {},
		"cgroup":     {},
		"cgroup2":    {},
		"nsfs":       {},
		"tracefs":    {},
		"debugfs":    {},
		"configfs":   {},
		"securityfs": {},
		"pstore":     {},
		"fusectl":    {},
		"hugetlbfs":  {},
		"rpc_pipefs": {},
	}
	ignoredMountNames := map[string]struct{}{
		"proc":  {},
		"sys":   {},
		"dev":   {},
		"run":   {},
		"tmp":   {},
		"etc":   {},
		"usr":   {},
		"lib":   {},
		"lib64": {},
		"bin":   {},
		"sbin":  {},
		"boot":  {},
		"var":   {},
	}

	lines := strings.Split(string(content), "\n")
	roots := make([]string, 0, len(lines))
	seen := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " - ", 2)
		if len(parts) != 2 {
			continue
		}
		left := strings.Fields(parts[0])
		right := strings.Fields(parts[1])
		if len(left) < 5 || len(right) < 1 {
			continue
		}

		mountPoint := decodeMountInfoField(left[4])
		if mountPoint == "" || mountPoint == "/" {
			continue
		}
		if _, ok := ignoredFS[right[0]]; ok {
			continue
		}
		if !isTopLevelMount(mountPoint) {
			continue
		}
		name := strings.Trim(mountPoint, "/")
		if _, ok := ignoredMountNames[name]; ok {
			continue
		}
		info, err := os.Stat(mountPoint)
		if err != nil || !info.IsDir() {
			continue
		}
		clean := filepath.Clean(mountPoint)
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		roots = append(roots, clean)
	}

	sort.Strings(roots)
	return roots
}

func isTopLevelMount(path string) bool {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return false
	}
	return !strings.Contains(trimmed, "/")
}

func decodeMountInfoField(value string) string {
	if value == "" || !strings.Contains(value, "\\") {
		return value
	}

	var builder strings.Builder
	builder.Grow(len(value))
	for index := 0; index < len(value); index++ {
		if value[index] == '\\' && index+3 < len(value) && isOctal(value[index+1]) && isOctal(value[index+2]) && isOctal(value[index+3]) {
			decoded := (int(value[index+1]-'0') << 6) + (int(value[index+2]-'0') << 3) + int(value[index+3]-'0')
			builder.WriteByte(byte(decoded))
			index += 3
			continue
		}
		builder.WriteByte(value[index])
	}
	return builder.String()
}

func isOctal(ch byte) bool {
	return ch >= '0' && ch <= '7'
}
