package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hidez8891/zip"
)

type zipFile struct {
	*zip.Writer
}

type archiveWriteFunc func(info os.FileInfo, file io.Reader, entryName string) (err error)

// GenerateZipContent return zip content.
// We should not use the standard zip package, see https://github.com/golang/go/issues/23301
func GenerateZipContent(folderName string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	zip := zipFile{zip.NewWriter(buf)}

	zip.AddAll(folderName, false)

	if err := zip.Close(); err != nil {
		return nil, fmt.Errorf("unable to generate zip content: %v", err)
	}

	return buf, nil
}

// AddAll adds all files from dir in archive, recursively.
// Directories receive a zero-size entry in the archive, with a trailing slash in the header name, and no compression
func (z *zipFile) AddAll(dir string, includeCurrentFolder bool) error {
	dir = path.Clean(dir)
	return addAll(dir, dir, includeCurrentFolder, func(info os.FileInfo, file io.Reader, entryName string) (err error) {
		// If we have a file to write (i.e., not a directory) then pipe the file into the archive writer
		if file != nil {
			writer, _ := z.Create(entryName, true)
			if _, err := io.Copy(writer, file); err != nil {
				return err
			}
		}

		return nil
	})
}

// addAll is used to recursively go down through directories and add each file and directory to an archive, based on an archiveWriteFunc given to it
func addAll(dir string, rootDir string, includeCurrentFolder bool, writerFunc archiveWriteFunc) error {
	// Get a list of all entries in the directory, as []os.FileInfo
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// Loop through all entries
	for _, info := range fileInfos {

		full := filepath.Join(dir, info.Name())

		// If the entry is a file, get an io.Reader for it
		var file *os.File
		var reader io.Reader
		if !info.IsDir() {
			file, err = os.Open(full)
			if err != nil {
				return err
			}
			reader = file
		}

		// Write the entry into the archive
		subDir := getSubDir(dir, rootDir, includeCurrentFolder)
		entryName := path.Join(subDir, info.Name())
		if err := writerFunc(info, reader, entryName); err != nil {
			if file != nil {
				file.Close()
			}
			return err
		}

		if file != nil {
			if err := file.Close(); err != nil {
				return err
			}

		}

		// If the entry is a directory, recurse into it
		if info.IsDir() {
			addAll(full, rootDir, includeCurrentFolder, writerFunc)
		}
	}

	return nil
}

func getSubDir(dir string, rootDir string, includeCurrentFolder bool) (subDir string) {

	subDir = strings.Replace(dir, rootDir, "", 1)
	// Remove leading slashes, since this is intentionally a subdirectory.
	if len(subDir) > 0 && subDir[0] == os.PathSeparator {
		subDir = subDir[1:]
	}
	subDir = path.Join(strings.Split(subDir, string(os.PathSeparator))...)

	if includeCurrentFolder {
		parts := strings.Split(rootDir, string(os.PathSeparator))
		subDir = path.Join(parts[len(parts)-1], subDir)
	}

	return
}
