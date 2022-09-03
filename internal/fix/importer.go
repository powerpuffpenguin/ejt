package fix

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-jsonnet"
)

type FileImporter struct {
	JPaths  []string
	fsCache map[string]*fsCacheEntry
}

type fsCacheEntry struct {
	exists   bool
	contents jsonnet.Contents
}

func (importer *FileImporter) tryPath(dir, importedPath string) (found bool, contents jsonnet.Contents, foundHere string, err error) {
	if importer.fsCache == nil {
		importer.fsCache = make(map[string]*fsCacheEntry)
	}
	var absPath string
	if filepath.IsAbs(importedPath) {
		absPath = importedPath
	} else {
		absPath = filepath.Join(dir, importedPath)
	}
	var entry *fsCacheEntry
	if cacheEntry, isCached := importer.fsCache[absPath]; isCached {
		entry = cacheEntry
	} else {
		contentBytes, err := ioutil.ReadFile(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				entry = &fsCacheEntry{
					exists: false,
				}
			} else {
				return false, jsonnet.Contents{}, "", err
			}
		} else {
			entry = &fsCacheEntry{
				exists:   true,
				contents: jsonnet.MakeContents(string(contentBytes)),
			}
		}
		importer.fsCache[absPath] = entry
	}
	return entry.exists, entry.contents, absPath, nil
}

// Import imports file from the filesystem.
func (importer *FileImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	// TODO(sbarzowski) Make sure that dir is absolute and resolving of ""
	// is independent from current CWD. The default path should be saved
	// in the importer.
	// We need to relativize the paths in the error formatter, so that the stack traces
	// don't have ugly absolute paths (less readable and messy with golden tests).
	dir, _ := filepath.Split(importedFrom)
	found, content, foundHere, err := importer.tryPath(dir, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}

	for i := len(importer.JPaths) - 1; !found && i >= 0; i-- {
		found, content, foundHere, err = importer.tryPath(importer.JPaths[i], importedPath)
		if err != nil {
			return jsonnet.Contents{}, "", err
		}
	}

	if !found {
		return jsonnet.Contents{}, "", fmt.Errorf("couldn't open import %#v: no match locally or in the Jsonnet library paths", importedPath)
	}
	return content, foundHere, nil
}
