// Copyright 2012 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package pebble

import (
	"fmt"

	"github.com/cockroachdb/pebble/internal/base"
	"github.com/cockroachdb/pebble/vfs"
)

type fileType = base.FileType

const (
	fileTypeLog      = base.FileTypeLog
	fileTypeLock     = base.FileTypeLock
	fileTypeTable    = base.FileTypeTable
	fileTypeManifest = base.FileTypeManifest
	fileTypeCurrent  = base.FileTypeCurrent
	fileTypeOptions  = base.FileTypeOptions
)

func setCurrentFile(dirname string, fs vfs.FS, fileNum uint64) error {
	newFilename := base.MakeFilename(dirname, fileTypeCurrent, fileNum)
	oldFilename := fmt.Sprintf("%s.%06d.dbtmp", newFilename, fileNum)
	fs.Remove(oldFilename)
	f, err := fs.Create(oldFilename)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(f, "MANIFEST-%06d\n", fileNum); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return fs.Rename(oldFilename, newFilename)
}
