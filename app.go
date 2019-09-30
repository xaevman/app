//  ---------------------------------------------------------------------------
//
//  app.go
//
//  Copyright (c) 2014, Jared Chavez.
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

// Package app provides extended facilities for gathering information about the
// executing application.
package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

// GetExeDir returns the base directory of the executing binary.
func GetExeDir() string {
	return filepath.Dir(GetExePath())
}

// GetExeFile returns the file name of the executing binary.
func GetExeFile() string {
	return filepath.Base(GetExePath())
}

// GetExePath attempts to find and return the executing binary's full path.
func GetExePath() string {
	rootPath := filepath.Dir(os.Args[0])

	if rootPath == "." {
		rootPath, err := exec.LookPath(os.Args[0])
		if err != nil {
			return os.Args[0]
		}

		rootPath, err = filepath.Abs(rootPath)
		if err != nil {
			return os.Args[0]
		}

		return rootPath
	}

	return os.Args[0]
}

// GetName returns the name of the executing application, as parsed from
// the executing binary's file name.
func GetName() string {
	fn := GetExeFile()
	ext := filepath.Ext(fn)

	return fn[:len(fn)-len(ext)]
}

// GetRunStatus attempts to read the application pid file and check the status
// of the given pid. It will return the active process, if one exists, or nil if
// there is no pid file, or the documented pid is no longer active. If a pid file
// is found, but the pid is no longer active, GetRunStatus will automatically
// attempt to clean up the lingering pid file.
func GetRunStatus() *os.Process {
	pid, err := ReadPidFile()
	if err != nil {
		// can't find pid file
		return nil
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		// pid file exists, but pid no longer does
		DeletePidFile()
		return nil
	}

	proc, err := ps.FindProcess(pid)
	if err != nil {
		return p
	}

	if proc == nil || (strings.ToLower(GetExeFile()) != strings.ToLower(proc.Executable())) {
		// pid exists, but is tied to a different application now
		DeletePidFile()
		return nil
	}

	return p
}

// init initializes the package.
func init() {
	SetPidBaseDir(GetExeDir())
}
