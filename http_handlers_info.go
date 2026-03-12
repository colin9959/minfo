package main

import (
	"context"
	"fmt"
	"net/http"
)

func infoHandler(envKey, fallback string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ensurePost(w, r) {
			return
		}
		if err := parseForm(w, r); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanupMultipart(r)

		path, cleanup, err := inputPath(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanup()

		bin, err := resolveBin(envKey, fallback)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		stdout, stderr, err := runCommand(ctx, bin, path)
		if err != nil {
			writeError(w, http.StatusInternalServerError, bestErrorMessage(err, stderr, stdout))
			return
		}

		output := combineCommandOutput(stdout, stderr)
		if output == "" {
			writeError(w, http.StatusInternalServerError, "mediainfo returned empty output")
			return
		}

		writeJSON(w, http.StatusOK, infoResponse{OK: true, Output: output})
	}
}

func mediainfoHandler(envKey, fallback string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ensurePost(w, r) {
			return
		}
		if err := parseForm(w, r); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanupMultipart(r)

		path, cleanup, err := inputPath(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanup()

		bin, err := resolveBin(envKey, fallback)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		candidates, sourceCleanup, err := resolveMediaInfoCandidates(ctx, path, mediaInfoCandidateLimit)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer sourceCleanup()

		var lastErr string
		for _, sourcePath := range candidates {
			stdout, stderr, err := runCommand(ctx, bin, sourcePath)
			if err != nil {
				lastErr = bestErrorMessage(err, stderr, stdout)
				continue
			}

			output := combineCommandOutput(stdout, stderr)
			if output == "" {
				lastErr = fmt.Sprintf("mediainfo returned empty output for: %s", sourcePath)
				continue
			}

			writeJSON(w, http.StatusOK, infoResponse{OK: true, Output: output})
			return
		}

		if lastErr == "" {
			lastErr = "mediainfo returned empty output"
		}
		writeError(w, http.StatusInternalServerError, lastErr)
	}
}

func bdinfoHandler(envKey, fallback string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !ensurePost(w, r) {
			return
		}
		if err := parseForm(w, r); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanupMultipart(r)

		path, cleanup, err := inputPath(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer cleanup()

		bin, err := resolveBin(envKey, fallback)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		bdPath, bdCleanup, err := resolveBDInfoSource(ctx, path)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer bdCleanup()

		stdout, stderr, err := runCommand(ctx, bin, bdPath)
		if err != nil {
			writeError(w, http.StatusInternalServerError, bestErrorMessage(err, stderr, stdout))
			return
		}

		writeJSON(w, http.StatusOK, infoResponse{OK: true, Output: combineCommandOutput(stdout, stderr)})
	}
}
