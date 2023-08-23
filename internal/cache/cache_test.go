package cache

import (
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	if err := Setup("testcache"); err != nil {
		t.Errorf("expected no error; got %s", err.Error())
	}

	if _, err := os.Stat("testcache"); os.IsNotExist(err) {
		t.Error("expected testcache to exist")
	}

	if err := Setup("testcache"); err != CacheIsAlreadySetup {
		t.Errorf("expected error of CacheIsAlreadySetup; got %s", err.Error())
	}

	if err := Destroy(); err != nil {
		t.Errorf("expected no error; got %s", err.Error())
	}

	os.Mkdir("testcache", 0755)
	if err := Setup("testcache"); err != FailedAtSettingUpCache {
		t.Errorf("expected error of FailedAtSettingUpCache; got %s", err.Error())
	}
	Destroy()
}

func TestAdd(t *testing.T) {
	t.Skip() // Fix oaths
	Setup("testcache")
	defer Destroy()

	expected := "testcache/8458bc5ba4df1237e49f36863f675a9ed1551c52/v1.0.0"

	wd, _ := os.Getwd()

	path, err := Add(wd, "v1.0.0")

	if err != nil {
		t.Errorf("expected no error; got %s", err.Error())
		return
	}

	if expected != path {
		t.Errorf("expected path of %s; got %s", expected, path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected path to exist")
	}
}

func TestGet(t *testing.T) {
	t.Skip() // Fix oaths
	Setup("testcache")
	defer Destroy()

	expected := "testcache/8458bc5ba4df1237e49f36863f675a9ed1551c52/v1.0.0"

	wd, _ := os.Getwd()
	Add(wd, "v1.0.0")

	path, err := Get(wd, "v1.0.0")

	if err != nil {
		t.Errorf("expected no error; got %s", err)
	}

	if expected != path {
		t.Errorf("expected path %s; got %s", expected, path)
	}
}

func TestHas(t *testing.T) {
	Setup("testcache")
	defer Destroy()

	wd, _ := os.Getwd()
	Add(wd, "v1.0.0")

	if !Has(wd, "v1.0.0") {
		t.Error("expected cache to contain repository and version")
	}

	if Has(wd, "v1.0.1") {
		t.Error("expected cache to not have version")
	}

	if Has("./other-repo", "v1.0.0") {
		t.Error("expected cache to not have repo")
	}
}
