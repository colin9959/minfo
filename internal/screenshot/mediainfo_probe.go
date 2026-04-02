package screenshot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"minfo/internal/system"
)

type nativeDVDMediaInfoTrack struct {
	StreamID int
	ID       string
	Format   string
	Language string
	Title    string
	Source   string
}

type nativeDVDMediaInfoResult struct {
	Duration             float64
	Tracks               []nativeDVDMediaInfoTrack
	ProbePath            string
	SelectedVOBPath      string
	LanguageFallbackPath string
}

type nativeMediaInfoPayload struct {
	Media struct {
		Track []map[string]interface{} `json:"track"`
	} `json:"media"`
}

var nativeMediaInfoHexIDPattern = regexp.MustCompile(`(?i)\(0x([0-9a-f]+)\)`)

func probeDVDMediaInfo(ctx context.Context, mediainfoBin, path, probePath string) (nativeDVDMediaInfoResult, error) {
	if strings.TrimSpace(mediainfoBin) == "" {
		return nativeDVDMediaInfoResult{}, fmt.Errorf("mediainfo not available")
	}

	selectedVOBPath := nativeResolveDVDMediaInfoVOBPath(path, probePath)
	primaryPath := nativeResolveDVDMediaInfoProbePath(path, probePath)
	result, err := nativeProbeDVDMediaInfoOnce(ctx, mediainfoBin, primaryPath)
	if err != nil {
		return nativeDVDMediaInfoResult{}, err
	}
	result.ProbePath = primaryPath
	result.SelectedVOBPath = selectedVOBPath

	if nativeDVDMediaInfoNeedsLanguageFallback(result) {
		if bupPath, ok := nativeDVDMediaInfoBUPPath(primaryPath); ok {
			fallback, fallbackErr := nativeProbeDVDMediaInfoOnce(ctx, mediainfoBin, bupPath)
			if fallbackErr == nil {
				merged, used := nativeMergeDVDMediaInfoLanguageFallback(result, fallback)
				if used {
					merged.ProbePath = result.ProbePath
					merged.SelectedVOBPath = result.SelectedVOBPath
					merged.LanguageFallbackPath = bupPath
					result = merged
				}
			}
		}
	}
	return result, nil
}

func nativeProbeDVDMediaInfoOnce(ctx context.Context, mediainfoBin, path string) (nativeDVDMediaInfoResult, error) {
	stdout, stderr, err := system.RunCommand(ctx, mediainfoBin, "--Output=JSON", path)
	if err != nil {
		return nativeDVDMediaInfoResult{}, fmt.Errorf(system.BestErrorMessage(err, stderr, stdout))
	}
	if strings.TrimSpace(stdout) == "" {
		return nativeDVDMediaInfoResult{}, fmt.Errorf("mediainfo returned empty output")
	}

	payloads, err := nativeDecodeMediaInfoPayloads([]byte(stdout))
	if err != nil {
		return nativeDVDMediaInfoResult{}, err
	}

	result := nativeDVDMediaInfoResult{
		Tracks: make([]nativeDVDMediaInfoTrack, 0, 8),
	}

	for _, payload := range payloads {
		for _, track := range payload.Media.Track {
			trackType := strings.TrimSpace(nativeMediaInfoTrackString(track, "@type"))
			switch strings.ToLower(trackType) {
			case "general":
				if value, ok := nativeParseMediaInfoTrackDuration(nativeMediaInfoTrackString(track, "Duration")); ok && value > 0 && value > result.Duration {
					result.Duration = value
				}
			case "text":
				streamID, _ := nativeParseMediaInfoStreamID(nativeMediaInfoTrackString(track, "ID"))
				result.Tracks = append(result.Tracks, nativeDVDMediaInfoTrack{
					StreamID: streamID,
					ID:       nativeMediaInfoTrackString(track, "ID"),
					Format:   nativeMediaInfoTrackString(track, "Format", "Format/String"),
					Language: nativeMediaInfoTrackString(track, "Language", "Language/String"),
					Title:    nativeMediaInfoTrackString(track, "Title", "Title/String"),
					Source:   nativeMediaInfoTrackString(track, "Source"),
				})
			}
		}
	}

	sort.Slice(result.Tracks, func(i, j int) bool {
		if result.Tracks[i].StreamID != result.Tracks[j].StreamID {
			return result.Tracks[i].StreamID < result.Tracks[j].StreamID
		}
		return result.Tracks[i].ID < result.Tracks[j].ID
	})
	return result, nil
}

func nativeResolveDVDMediaInfoProbePath(path, probePath string) string {
	for _, candidate := range []string{probePath, path} {
		if resolved, ok := nativeDVDMediaInfoIFOPath(candidate); ok {
			return resolved
		}
	}
	return strings.TrimSpace(path)
}

func nativeResolveDVDMediaInfoVOBPath(path, probePath string) string {
	for _, candidate := range []string{probePath, path} {
		if resolved, ok := nativeDVDMediaInfoTitleVOBPath(candidate); ok {
			return resolved
		}
	}
	return ""
}

func nativeDVDMediaInfoIFOPath(path string) (string, bool) {
	cleaned := strings.TrimSpace(path)
	if cleaned == "" {
		return "", false
	}

	upperBase := strings.ToUpper(filepath.Base(cleaned))
	switch filepath.Ext(upperBase) {
	case ".IFO":
		if nativeFileExists(cleaned) {
			return cleaned, true
		}
	case ".BUP":
		ifoPath := strings.TrimSuffix(cleaned, filepath.Ext(cleaned)) + ".IFO"
		if nativeFileExists(ifoPath) {
			return ifoPath, true
		}
	case ".VOB":
		if strings.EqualFold(upperBase, "VIDEO_TS.VOB") {
			ifoPath := filepath.Join(filepath.Dir(cleaned), "VIDEO_TS.IFO")
			if nativeFileExists(ifoPath) {
				return ifoPath, true
			}
			return "", false
		}
		if len(upperBase) == len("VTS_00_1.VOB") &&
			strings.HasPrefix(upperBase, "VTS_") &&
			upperBase[6] == '_' &&
			upperBase[8:] == ".VOB" &&
			upperBase[7] >= '1' && upperBase[7] <= '9' {
			ifoPath := filepath.Join(filepath.Dir(cleaned), upperBase[:7]+"0.IFO")
			if nativeFileExists(ifoPath) {
				return ifoPath, true
			}
		}
	}
	return "", false
}

func nativeDVDMediaInfoTitleVOBPath(path string) (string, bool) {
	cleaned := strings.TrimSpace(path)
	if cleaned == "" {
		return "", false
	}

	upperBase := strings.ToUpper(filepath.Base(cleaned))
	switch filepath.Ext(upperBase) {
	case ".VOB":
		if nativeFileExists(cleaned) && nativeLooksLikeDVDSource(cleaned) {
			return cleaned, true
		}
	case ".IFO", ".BUP":
		if strings.EqualFold(upperBase, "VIDEO_TS.IFO") || strings.EqualFold(upperBase, "VIDEO_TS.BUP") {
			return "", false
		}
		if len(upperBase) == len("VTS_00_0.IFO") &&
			strings.HasPrefix(upperBase, "VTS_") &&
			upperBase[6] == '_' &&
			upperBase[7] == '0' {
			vobPath := filepath.Join(filepath.Dir(cleaned), upperBase[:7]+"1.VOB")
			if nativeFileExists(vobPath) {
				return vobPath, true
			}
		}
	}
	return "", false
}

func nativeDVDMediaInfoBUPPath(path string) (string, bool) {
	cleaned := strings.TrimSpace(path)
	if !strings.EqualFold(filepath.Ext(cleaned), ".ifo") {
		return "", false
	}
	bupPath := strings.TrimSuffix(cleaned, filepath.Ext(cleaned)) + ".BUP"
	if !nativeFileExists(bupPath) {
		return "", false
	}
	return bupPath, true
}

func nativeDVDMediaInfoNeedsLanguageFallback(result nativeDVDMediaInfoResult) bool {
	if len(result.Tracks) == 0 {
		return false
	}
	for _, track := range result.Tracks {
		if !nativeDVDMediaInfoHasLanguage(track.Language) {
			return true
		}
	}
	return false
}

func nativeMergeDVDMediaInfoLanguageFallback(primary, fallback nativeDVDMediaInfoResult) (nativeDVDMediaInfoResult, bool) {
	if len(primary.Tracks) == 0 || len(fallback.Tracks) == 0 {
		return primary, false
	}

	merged := primary
	merged.Tracks = append([]nativeDVDMediaInfoTrack(nil), primary.Tracks...)

	fallbackByStreamID := make(map[int][]nativeDVDMediaInfoTrack, len(fallback.Tracks))
	fallbackIndexByStreamID := make(map[int][]int, len(fallback.Tracks))
	fallbackOrdered := make([]nativeDVDMediaInfoTrack, 0, len(fallback.Tracks))
	for _, track := range fallback.Tracks {
		if !nativeDVDMediaInfoHasLanguage(track.Language) {
			continue
		}
		fallbackOrdered = append(fallbackOrdered, track)
		orderedIndex := len(fallbackOrdered) - 1
		if track.StreamID > 0 {
			fallbackByStreamID[track.StreamID] = append(fallbackByStreamID[track.StreamID], track)
			fallbackIndexByStreamID[track.StreamID] = append(fallbackIndexByStreamID[track.StreamID], orderedIndex)
		}
	}
	if len(fallbackOrdered) == 0 {
		return primary, false
	}

	usedOrdered := make([]bool, len(fallbackOrdered))
	used := false
	for index := range merged.Tracks {
		if nativeDVDMediaInfoHasLanguage(merged.Tracks[index].Language) {
			continue
		}

		if mergedTrack, ok := nativeFillDVDMediaInfoLanguageByStreamID(merged.Tracks[index], fallbackByStreamID, fallbackIndexByStreamID, usedOrdered); ok {
			merged.Tracks[index] = mergedTrack
			used = true
			continue
		}

		for fallbackIndex, candidate := range fallbackOrdered {
			if usedOrdered[fallbackIndex] {
				continue
			}
			merged.Tracks[index] = nativeMergeDVDMediaInfoTrack(merged.Tracks[index], candidate)
			usedOrdered[fallbackIndex] = true
			used = true
			break
		}
	}

	if merged.Duration <= 0 && fallback.Duration > 0 {
		merged.Duration = fallback.Duration
	}
	return merged, used
}

func nativeFillDVDMediaInfoLanguageByStreamID(track nativeDVDMediaInfoTrack, fallbackByStreamID map[int][]nativeDVDMediaInfoTrack, fallbackIndexByStreamID map[int][]int, usedOrdered []bool) (nativeDVDMediaInfoTrack, bool) {
	if track.StreamID <= 0 {
		return track, false
	}
	candidates := fallbackByStreamID[track.StreamID]
	candidateIndexes := fallbackIndexByStreamID[track.StreamID]
	if len(candidates) == 0 || len(candidateIndexes) == 0 {
		return track, false
	}
	for index, candidate := range candidates {
		if index >= len(candidateIndexes) {
			break
		}
		orderedIndex := candidateIndexes[index]
		if orderedIndex < 0 || orderedIndex >= len(usedOrdered) || usedOrdered[orderedIndex] {
			continue
		}
		usedOrdered[orderedIndex] = true
		return nativeMergeDVDMediaInfoTrack(track, candidate), true
	}
	return track, false
}

func nativeMergeDVDMediaInfoTrack(track, fallback nativeDVDMediaInfoTrack) nativeDVDMediaInfoTrack {
	if !nativeDVDMediaInfoHasLanguage(track.Language) && nativeDVDMediaInfoHasLanguage(fallback.Language) {
		track.Language = fallback.Language
		track.Source = fallback.Source
	}
	if strings.TrimSpace(track.Title) == "" && strings.TrimSpace(fallback.Title) != "" {
		track.Title = strings.TrimSpace(fallback.Title)
	}
	if strings.TrimSpace(track.Format) == "" && strings.TrimSpace(fallback.Format) != "" {
		track.Format = strings.TrimSpace(fallback.Format)
	}
	return track
}

func nativeDVDMediaInfoHasLanguage(language string) bool {
	switch strings.ToLower(strings.TrimSpace(language)) {
	case "", "unknown", "und", "undefined", "null", "n/a", "na":
		return false
	default:
		return true
	}
}

func nativeFileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func nativeDecodeMediaInfoPayloads(data []byte) ([]nativeMediaInfoPayload, error) {
	var single nativeMediaInfoPayload
	if err := json.Unmarshal(data, &single); err == nil {
		return []nativeMediaInfoPayload{single}, nil
	}

	var multiple []nativeMediaInfoPayload
	if err := json.Unmarshal(data, &multiple); err == nil {
		return multiple, nil
	}

	return nil, fmt.Errorf("unsupported mediainfo JSON shape")
}

func nativeMediaInfoTrackString(track map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(nativeJSONString(track[key]))
		if value != "" {
			return value
		}
	}
	return ""
}

func nativeParseMediaInfoStreamID(raw string) (int, bool) {
	matches := nativeMediaInfoHexIDPattern.FindAllStringSubmatch(raw, -1)
	if len(matches) > 0 {
		last := matches[len(matches)-1]
		if len(last) >= 2 {
			value, err := strconv.ParseInt(last[1], 16, 64)
			if err == nil {
				return int(value), true
			}
		}
	}

	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	if value, err := strconv.Atoi(raw); err == nil {
		return value, true
	}
	return 0, false
}

func nativeParseMediaInfoTrackDuration(raw string) (float64, bool) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return 0, false
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed <= 0 {
		return 0, false
	}
	return parsed / 1000.0, true
}

func nativeResolveDVDMediaInfoTracks(raw []nativeSubtitleTrack, tracks []nativeDVDMediaInfoTrack) map[int]nativeDVDMediaInfoTrack {
	resolved := make(map[int]nativeDVDMediaInfoTrack, len(tracks))
	if len(raw) == 0 || len(tracks) == 0 {
		return resolved
	}

	rawPIDSet := make(map[int]struct{}, len(raw))
	for _, track := range raw {
		pid, ok := nativeNormalizeStreamPID(track.StreamID)
		if !ok {
			continue
		}
		rawPIDSet[pid] = struct{}{}
	}

	exactMatched := false
	for _, item := range tracks {
		if item.StreamID <= 0 {
			continue
		}
		if _, ok := rawPIDSet[item.StreamID]; ok {
			resolved[item.StreamID] = item
			exactMatched = true
		}
	}
	if exactMatched {
		return resolved
	}

	type rawTrackPID struct {
		pid int
	}
	rawPIDs := make([]rawTrackPID, 0, len(raw))
	for _, track := range raw {
		if strings.ToLower(strings.TrimSpace(track.Codec)) != "dvd_subtitle" {
			continue
		}
		pid, ok := nativeNormalizeStreamPID(track.StreamID)
		if !ok {
			continue
		}
		rawPIDs = append(rawPIDs, rawTrackPID{pid: pid})
	}
	sort.Slice(rawPIDs, func(i, j int) bool {
		return rawPIDs[i].pid < rawPIDs[j].pid
	})

	mediaInfoCopy := append([]nativeDVDMediaInfoTrack(nil), tracks...)
	sort.Slice(mediaInfoCopy, func(i, j int) bool {
		if mediaInfoCopy[i].StreamID != mediaInfoCopy[j].StreamID {
			return mediaInfoCopy[i].StreamID < mediaInfoCopy[j].StreamID
		}
		return mediaInfoCopy[i].ID < mediaInfoCopy[j].ID
	})

	limit := len(mediaInfoCopy)
	if len(rawPIDs) < limit {
		limit = len(rawPIDs)
	}
	for index := 0; index < limit; index++ {
		resolved[rawPIDs[index].pid] = mediaInfoCopy[index]
	}
	return resolved
}
