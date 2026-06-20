// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pathutil

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// InDir checks whether path is in the file tree rooted at dir.
// It checks only the lexical form of the file names.
// It does not consider symbolic links.
//
// Copied from go/src/cmd/go/internal/search/search.go.
func InDir(dir, path string) bool {
	pv := strings.ToUpper(filepath.VolumeName(path))
	dv := strings.ToUpper(filepath.VolumeName(dir))
	path = path[len(pv):]
	dir = dir[len(dv):]
	switch {
	default:
		return false
	case pv != dv:
		return false
	case len(path) == len(dir):
		if path == dir {
			return true
		}
		return false
	case dir == "":
		return path != ""
	case len(path) > len(dir):
		if dir[len(dir)-1] == filepath.Separator {
			if path[:len(dir)] == dir {
				return path[len(dir):] != ""
			}
			return false
		}
		if path[len(dir)] == filepath.Separator && path[:len(dir)] == dir {
			if len(path) == len(dir)+1 {
				return true
			}
			return path[len(dir)+1:] != ""
		}
		return false
	}
}

var systemTempDirOnce sync.Once
var systemTempDirValue string

// SystemTempDir returns the operating system's default temporary directory,
// ignoring TMP/TEMP/TMPDIR overrides. Tests often redirect os.TempDir to an
// isolated directory, but module and workspace discovery walks up to the real
// system temp directory and must ignore stray go.mod and go.work files there.
func SystemTempDir() string {
	systemTempDirOnce.Do(func() {
		if runtime.GOOS == "windows" {
			systemTempDirValue = filepath.Join(os.Getenv("LOCALAPPDATA"), "Temp")
			return
		}
		envs := []string{"TMP", "TEMP", "TMPDIR"}
		saved := make([]string, len(envs))
		for i, k := range envs {
			saved[i] = os.Getenv(k)
			os.Unsetenv(k)
		}
		systemTempDirValue = os.TempDir()
		for i, k := range envs {
			if saved[i] != "" {
				os.Setenv(k, saved[i])
			}
		}
	})
	return systemTempDirValue
}

// InSystemTempDir reports whether path is the system temp directory or a
// subdirectory of it.
func InSystemTempDir(path string) bool {
	return InDir(SystemTempDir(), path)
}
