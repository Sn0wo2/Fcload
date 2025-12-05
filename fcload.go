package Fcload

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(src any, ls ...Loader) error {
	if len(ls) == 0 {
		return errors.New("no loaders provided")
	}
	loaderByExt := make(map[string]Loader)

	for _, l := range ls {
		for _, ext := range l.GetAllowFileExtensions() {
			loaderByExt["."+strings.ToLower(ext)] = l
		}
	}

	Path = os.Getenv("CONFIG_PATH")

	if Path != "" {
		if _, err := os.Stat(Path); err != nil {
			base := strings.TrimSuffix(Path, filepath.Ext(Path))
			for ext := range loaderByExt {
				tryPath := base + ext
				if _, err := os.Stat(tryPath); err == nil {
					Path = tryPath

					break
				}
			}
		}
	}

	// p1 fallback
	if Path == "" {
		searchPaths := []string{"./data/"}

	searchLoop:
		for _, p := range searchPaths {
			for ext := range loaderByExt {
				fullPath := filepath.Join(p, "config"+ext)
				if _, err := os.Stat(fullPath); err == nil {
					Path = fullPath

					break searchLoop
				}
			}
		}
	}

	// p1 & p2 fallback
	if Path == "" {
		Path = "./data/config.yml"

		return ErrConfigNotFound
	}

	ext := strings.ToLower(filepath.Ext(Path))

	loader, ok := loaderByExt[ext]
	retryIndex := 0

retryLoaders:
	if !ok {
		if retryIndex >= len(ls) {
			return fmt.Errorf("no loader found for config file %s", Path)
		}

		loader = ls[retryIndex]
		retryIndex++
		_, _ = fmt.Fprintf(os.Stderr, "failed to find config loader %s. Retrying with next loader: %s %d/%d\n", Path, loader.GetTag(), retryIndex, len(ls))
	}

	if err := loader.Load(src, Path); err != nil {
		if !ok {
			_, _ = fmt.Fprintf(os.Stderr, "loader %s failed to load config file %s: %v. Retrying with next loader... %d/%d: %v\n", loader.GetTag(), Path, err, retryIndex, len(ls), err)

			goto retryLoaders
		}

		return fmt.Errorf("failed to load config file %s: %w", Path, err)
	}

	return nil
}
