package cache

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	PathIsAlreadyCached    Error = "path is already cached"
	PathIsNotCached        Error = "path is not cached"
	CacheIsAlreadySetup    Error = "cache is already setup"
	CacheIsNotSetup        Error = "cache is not setup"
	FailedAtSettingUpCache Error = "failed at setting up cache"
	RepoNotInCache         Error = "repository not in cache"
	RefNotInCache          Error = "reference not in cache"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type repo struct {
	hash       string
	references map[string]struct{}
}

func (r repo) Add(reference string) {
	r.references[reference] = struct{}{}
}

func (r repo) Has(reference string) bool {
	_, found := r.references[reference]
	return found
}

var (
	dir   string
	cache map[string]repo
)

func Setup(path string) error {
	if cache != nil {
		return CacheIsAlreadySetup
	}

	dir = path
	cache = map[string]repo{}

	if err := os.Mkdir(dir, 0755); err != nil {
		return FailedAtSettingUpCache
	}

	return nil
}

func Add(path, reference string) (string, error) {
	if cache == nil {
		return "", CacheIsNotSetup
	}

	h := sha1.New()
	io.WriteString(h, path)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	if _, ok := cache[path]; !ok {
		cache[path] = repo{hash, map[string]struct{}{}}
		if err := os.Mkdir(filepath.Join(dir, hash), 0755); err != nil {
			return "", err
		}
	}

	if !cache[path].Has(reference) {
		cache[path].Add(reference)
		if err := os.Mkdir(filepath.Join(dir, hash, reference), 0755); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, hash, reference), nil
}

func Get(path, reference string) (string, error) {
	repo, ok := cache[path]

	if !ok {
		return "", RepoNotInCache
	}

	if !repo.Has(reference) {
		return "", RefNotInCache
	}

	return filepath.Join(dir, repo.hash, reference), nil
}

func Has(path, reference string) bool {
	repo, ok := cache[path]

	if !ok {
		return false
	}

	return repo.Has(reference)
}

func Clean() error {
	if cache == nil {
		return CacheIsNotSetup
	}
	return nil
}

func Destroy() error {
	if cache == nil {
		return CacheIsNotSetup
	}

	cache = nil
	return os.RemoveAll(dir)
}
