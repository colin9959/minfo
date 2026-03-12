package main

import (
	"net/http"
	"strings"
)

func pathSuggestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writePathError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	roots, err := resolveRoots(mediaRoots())
	if err != nil {
		writePathError(w, http.StatusBadRequest, err.Error())
		return
	}
	prefix := strings.TrimSpace(r.URL.Query().Get("prefix"))
	prefix = strings.Trim(prefix, "\"")

	items, root, err := suggestPaths(roots, prefix, maxSuggestions)
	if err != nil {
		writePathError(w, http.StatusBadRequest, err.Error())
		return
	}

	writePathJSON(w, http.StatusOK, pathResponse{
		OK:    true,
		Root:  root,
		Roots: roots,
		Items: items,
	})
}
