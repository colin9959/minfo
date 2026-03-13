package media

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SuggestPaths(roots []string, prefix string, limit int) ([]string, string, error) {
	if len(roots) == 0 {
		return nil, "", errors.New("no MEDIA_ROOT configured")
	}

	if prefix == "" {
		if len(roots) == 1 {
			items, err := listDir(roots[0], "", limit)
			return items, roots[0], err
		}
		items := make([]string, 0, len(roots))
		for _, root := range roots {
			items = append(items, withDirSuffix(root))
		}
		return items, "", nil
	}

	cleaned := filepath.Clean(prefix)
	selectedRoot := ""
	var absPrefix string
	if filepath.IsAbs(cleaned) {
		var ok bool
		absPrefix = cleaned
		selectedRoot, ok = findContainingRoot(roots, absPrefix)
		if !ok {
			return nil, "", errors.New("path is outside MEDIA_ROOTS")
		}
	} else {
		if len(roots) != 1 {
			return nil, "", errors.New("relative path requires a single MEDIA_ROOT")
		}
		selectedRoot = roots[0]
		absPrefix = filepath.Join(selectedRoot, cleaned)
	}

	sep := string(filepath.Separator)
	if strings.HasSuffix(prefix, sep) || strings.HasSuffix(prefix, "/") || strings.HasSuffix(prefix, "\\") {
		if !isSubpath(selectedRoot, absPrefix) {
			return nil, "", errors.New("path is outside MEDIA_ROOTS")
		}
		items, err := listDir(absPrefix, "", limit)
		return items, selectedRoot, err
	}

	dir := filepath.Dir(absPrefix)
	base := filepath.Base(absPrefix)
	if !isSubpath(selectedRoot, dir) {
		return nil, "", errors.New("path is outside MEDIA_ROOTS")
	}
	items, err := listDir(dir, base, limit)
	return items, selectedRoot, err
}

func ResolveRoots(roots []string) ([]string, error) {
	resolved := make([]string, 0, len(roots))
	seen := make(map[string]struct{}, len(roots))
	for _, root := range roots {
		root = filepath.Clean(strings.TrimSpace(root))
		if root == "" {
			continue
		}
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			continue
		}
		if _, ok := seen[root]; ok {
			continue
		}
		seen[root] = struct{}{}
		resolved = append(resolved, root)
	}
	if len(resolved) == 0 {
		return nil, errors.New("no MEDIA_ROOT configured")
	}
	sort.Strings(resolved)
	return resolved, nil
}

func findContainingRoot(roots []string, path string) (string, bool) {
	for _, root := range roots {
		if isSubpath(root, path) {
			return root, true
		}
	}
	return "", false
}

func withDirSuffix(path string) string {
	if strings.HasSuffix(path, string(filepath.Separator)) {
		return path
	}
	return path + string(filepath.Separator)
}

func listDir(dir, base string, limit int) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	items := make([]string, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if base != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(base)) {
			continue
		}
		full := filepath.Join(dir, name)
		if entry.IsDir() {
			full = withDirSuffix(full)
		}
		items = append(items, full)
		if limit > 0 && len(items) >= limit {
			break
		}
	}
	sort.Strings(items)
	return items, nil
}

func isSubpath(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	if root == path {
		return true
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}
