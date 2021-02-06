package tarfs

import (
	"bytes"
	"io"
	"io/fs"
	"time"
)

type file struct {
	entry
	r          io.ReadSeeker
	readDirPos int
}

func newFile(e entry) *file {
	return &file{e, bytes.NewReader(e.b), 0}
}

var _ fs.File = &file{}

func (f *file) Stat() (fs.FileInfo, error) {
	return f.h.FileInfo(), nil
}

func (f *file) Read(b []byte) (int, error) {
	if f.IsDir() {
		return 0, newErrDir("read", f.Name())
	}

	return f.r.Read(b)
}

func (f *file) Close() error {
	return nil
}

var _ io.Seeker = &file{}

func (f *file) Seek(offset int64, whence int) (int64, error) {
	if f.IsDir() {
		return 0, newErrDir("seek", f.Name())
	}

	return f.r.Seek(offset, whence)
}

var _ fs.ReadDirFile = &file{}

func (f *file) ReadDir(n int) ([]fs.DirEntry, error) {
	if !f.IsDir() {
		return nil, newErrNotDir("readdir", f.Name())
	}

	if f.readDirPos >= len(f.entries) {
		if n <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}

	var entries []fs.DirEntry

	if n > 0 && f.readDirPos+n <= len(f.entries) {
		entries = f.entries[f.readDirPos : f.readDirPos+n]
		f.readDirPos += n
	} else {
		entries = f.entries[f.readDirPos:]
		f.readDirPos += len(entries)
	}

	return entries, nil
}

type rootFile struct {
	tfs        *tarfs
	readDirPos int
}

var _ fs.File = &rootFile{}

func (rf *rootFile) Stat() (fs.FileInfo, error) {
	return rf, nil
}

func (*rootFile) Read([]byte) (int, error) {
	return 0, newErrDir("read", ".")
}

func (*rootFile) Close() error {
	return nil
}

var _ fs.ReadDirFile = &rootFile{}

func (rf *rootFile) ReadDir(n int) ([]fs.DirEntry, error) {
	if rf.readDirPos >= len(rf.tfs.rootEntries) {
		if n <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}

	var entries []fs.DirEntry

	if n > 0 && rf.readDirPos+n <= len(rf.tfs.rootEntries) {
		entries = rf.tfs.rootEntries[rf.readDirPos : rf.readDirPos+n]
		rf.readDirPos += n
	} else {
		entries = rf.tfs.rootEntries[rf.readDirPos:]
		rf.readDirPos += len(entries)
	}

	return entries, nil
}

var _ fs.FileInfo = &rootFile{}

func (rf *rootFile) Name() string {
	return "."
}

func (rf *rootFile) Size() int64 {
	return 0
}

func (rf *rootFile) Mode() fs.FileMode {
	return fs.ModeDir | 0755
}

func (rf *rootFile) ModTime() time.Time {
	return time.Time{}
}

func (rf *rootFile) IsDir() bool {
	return true
}

func (rf *rootFile) Sys() interface{} {
	return nil
}
