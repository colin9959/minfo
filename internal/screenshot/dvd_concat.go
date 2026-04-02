package screenshot

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var nativeDVDTitleVOBPattern = regexp.MustCompile(`(?i)^VTS_(\d{2})_([1-9]\d*)\.VOB$`)

func nativeIsDVDTitleVOB(path string) bool {
	return nativeDVDTitleVOBPattern.MatchString(filepath.Base(strings.TrimSpace(path)))
}

func prepareSourceForFFmpeg(sourcePath string) (string, string, func(), error) {
	return sourcePath, sourcePath, func() {}, nil
}

func nativeCollectDVDTitleSetVOBs(sourcePath string) ([]string, error) {
	videoTSDir := filepath.Dir(sourcePath)
	match := nativeDVDTitleVOBPattern.FindStringSubmatch(filepath.Base(sourcePath))
	if len(match) != 3 {
		return nil, fmt.Errorf("dvd source is not a title-set VOB: %s", sourcePath)
	}

	titleSet := match[1]
	entries, err := os.ReadDir(videoTSDir)
	if err != nil {
		return nil, err
	}

	type part struct {
		index int
		path  string
	}

	parts := make([]part, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		submatch := nativeDVDTitleVOBPattern.FindStringSubmatch(entry.Name())
		if len(submatch) != 3 || submatch[1] != titleSet {
			continue
		}
		partIndex, err := strconv.Atoi(submatch[2])
		if err != nil {
			continue
		}
		parts = append(parts, part{
			index: partIndex,
			path:  filepath.Join(videoTSDir, entry.Name()),
		})
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("no title-set VOB files found for %s", sourcePath)
	}

	sort.Slice(parts, func(i, j int) bool {
		if parts[i].index != parts[j].index {
			return parts[i].index < parts[j].index
		}
		return parts[i].path < parts[j].path
	})

	result := make([]string, 0, len(parts))
	for _, item := range parts {
		result = append(result, item.path)
	}
	return result, nil
}

func nativeEscapeFFconcatPath(value string) string {
	return strings.ReplaceAll(value, "'", "'\\''")
}

func nativeIsFFconcatSource(path string) bool {
	return strings.EqualFold(filepath.Ext(strings.TrimSpace(path)), ".ffconcat")
}

func nativeReadFFconcatFiles(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dir := filepath.Dir(path)
	files := make([]string, 0, 8)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.EqualFold(line, "ffconcat version 1.0") {
			continue
		}
		if !strings.HasPrefix(strings.ToLower(line), "file ") {
			continue
		}

		raw := strings.TrimSpace(line[5:])
		resolved, ok := nativeParseFFconcatFileLine(raw)
		if !ok {
			continue
		}
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(dir, resolved)
		}
		files = append(files, filepath.Clean(resolved))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("ffconcat file contains no media entries: %s", path)
	}
	return files, nil
}

func nativeParseFFconcatFileLine(value string) (string, bool) {
	if value == "" {
		return "", false
	}
	if value[0] == '\'' && value[len(value)-1] == '\'' && len(value) >= 2 {
		return strings.ReplaceAll(value[1:len(value)-1], "'\\''", "'"), true
	}
	return value, true
}
