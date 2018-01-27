// Copyright 2016 hidez8891. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zip

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hidez8891/encstr"
)

type editorFileHeader struct {
	*FileHeader
	isNewFile bool
}

// Updater implements a zip file edit-writer.
type Updater struct {
	path    string
	r       *ReadCloser
	w       *Writer
	tmpw    *os.File
	File    []*editorFileHeader
	Comment *encstr.String
}

// OpenUpdater exist zip file for editing.
func OpenUpdater(path string) (*Updater, error) {
	r, err := OpenReader(path)
	if err != nil {
		return nil, err
	}

	u := &Updater{
		path:    path,
		r:       r,
		Comment: r.Comment,
	}

	u.File = make([]*editorFileHeader, 0)
	for _, file := range r.File {
		u.File = append(u.File, &editorFileHeader{
			FileHeader: &file.FileHeader,
			isNewFile:  false,
		})
	}

	return u, nil
}

// AppendFile add a file to the zip archive at last.
func (u *Updater) AppendFile(name string, streamMode bool) (io.Writer, error) {
	writer, err := u.createWriter(name, streamMode)
	if err != nil {
		return nil, err
	}

	u.File = append(u.File, &editorFileHeader{
		FileHeader: u.w.dir[len(u.w.dir)-1].FileHeader,
		isNewFile:  true,
	})

	return writer, nil
}

// InsertFile add a file to the zip archive at pos-th.
func (u *Updater) InsertFile(pos int, name string, streamMode bool) (io.Writer, error) {
	if len(u.File) >= pos {
		return nil, fmt.Errorf("zip: []File index out of range")
	}

	writer, err := u.createWriter(name, streamMode)
	if err != nil {
		return nil, err
	}

	u.File = append(u.File, nil)
	copy(u.File[pos+1:], u.File[pos:])
	u.File[pos] = &editorFileHeader{
		FileHeader: u.w.dir[len(u.w.dir)-1].FileHeader,
		isNewFile:  true,
	}

	return writer, nil
}

func (u *Updater) createWriter(name string, streamMode bool) (io.Writer, error) {
	if u.w == nil {
		var err error
		if u.tmpw, err = ioutil.TempFile("", "tmp_zip_updater"); err != nil {
			return nil, err
		}
		u.w = NewWriter(u.tmpw)
	}

	writer, err := u.w.Create(name, streamMode)
	if err != nil {
		return nil, err
	}

	return writer, nil
}

// SaveAs write all changes to new zip file and close file.
func (u *Updater) SaveAs(newpath string) error {
	newfile, err := os.Create(newpath)
	if err != nil {
		return err
	}

	defer func() {
		if newfile != nil {
			newfile.Close()
			os.Remove(newpath)
		}
	}()

	// reopen temporary file
	var tmpreader *Reader
	if u.w != nil {
		var err error
		if err = u.w.Close(); err != nil {
			return err
		}
		u.w = nil

		state, err := u.tmpw.Stat()
		if err != nil {
			return err
		}

		if _, err = u.tmpw.Seek(0, os.SEEK_SET); err != nil {
			return err
		}

		tmpreader, err = NewReader(u.tmpw, state.Size())
		if err != nil {
			return err
		}
	}

	// copy & write
	w := NewWriter(newfile)
	w.Comment = u.Comment
	for _, header := range u.File {
		var (
			file *File
			zipr *Reader
		)

		if header.isNewFile {
			zipr = tmpreader
		} else {
			zipr = &u.r.Reader
		}

		if zipr == nil {
			return fmt.Errorf("zip: zip reader has null pointer")
		}

		for _, f := range zipr.File {
			if f.Name.Str() == header.Name.Str() {
				file = f
				break
			}
		}

		if file == nil {
			return fmt.Errorf("zip: file %s does not exist", header.Name.Str())
		}

		if err := w.addFile(file); err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}

	// close file
	if err := newfile.Close(); err != nil {
		return err
	}
	newfile = nil

	if err := u.r.Close(); err != nil {
		return err
	}
	u.r = nil

	if u.tmpw != nil {
		if err := u.tmpw.Close(); err != nil {
			return err
		}
		if err := os.Remove(u.tmpw.Name()); err != nil {
			return err
		}
		u.tmpw = nil
	}

	return nil
}

// Save write all changes to current zip file and close file.
func (u *Updater) Save() error {
	tmpfile, err := ioutil.TempFile("", u.r.f.Name())
	if err != nil {
		return err
	}
	tmpfile.Close()
	tmpname := tmpfile.Name()

	// copy & write & close file
	if err := u.SaveAs(tmpname); err != nil {
		return err
	}

	// move & overwrite
	backuppath := u.path + ".bak"
	if err := os.Rename(u.path, backuppath); err != nil {
		return err
	}
	if err := os.Rename(tmpname, u.path); err != nil {
		return err
	}
	os.Remove(backuppath)

	return nil
}

// Close discard all changes and close file.
func (u *Updater) Close() error {
	if u.r != nil {
		return u.r.Close()
	}
	return nil
}
