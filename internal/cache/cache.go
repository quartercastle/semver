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
	VersionNotInCache      Error = "version not in cache"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

type repo struct {
	hash     string
	versions map[string]struct{}
}

func (r repo) Add(version string) {
	r.versions[version] = struct{}{}
}

func (r repo) Has(version string) bool {
	_, found := r.versions[version]
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

func Add(path, version string) (string, error) {
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

	if !cache[path].Has(version) {
		cache[path].Add(version)
		if err := os.Mkdir(filepath.Join(dir, hash, version), 0755); err != nil {
			return "", err
		}
	}
	return filepath.Join(dir, hash, version), nil
}

func Get(path, version string) (string, error) {
	repo, ok := cache[path]

	if !ok {
		return "", RepoNotInCache
	}

	if !repo.Has(version) {
		return "", VersionNotInCache
	}

	return filepath.Join(dir, repo.hash, version), nil
}

func Has(path, version string) bool {
	repo, ok := cache[path]

	if !ok {
		return false
	}

	return repo.Has(version)
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
