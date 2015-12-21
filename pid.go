//  ---------------------------------------------------------------------------
//
//  pid.go
//
//  Copyright (c) 2014, Jared Chavez. 
//  All rights reserved.
//
//  Use of this source code is governed by a BSD-style
//  license that can be found in the LICENSE file.
//
//  -----------

package app

import (
    "errors"
    "fmt"
    "os"
    "path"
    "strconv"
    "sync"
)

// NoPid is returned during any pid operations which fail to aquire a valid
// pid.
const NoPid = -1

// Module errors
var PidParseErr = errors.New("Unable to parse pid file")

// Module vars
var (
    basePidDir string
    pdMutex    sync.RWMutex
)

// CreatePidFile creates the calling application's pid file and returns the
// application's current pid for reference.
func CreatePidFile() (int, error) {
    f, err := os.Create(getPidFilePath())
    if err != nil {
        return NoPid, err
    }

    defer f.Close()

    pid := os.Getpid()

    _, err = f.WriteString(fmt.Sprintf("%d", pid))
    if err != nil {
        return pid, err
    }

    return pid, nil
}

// DeletePidFile attempts to remove the caller application's pid file (if any)
// and returns any errors from the underlying os.Remove call.
func DeletePidFile() error {
    return os.Remove(getPidFilePath())
}

// GetPidBaseDir reads the base directory for pid operations in a thread-safe way.
func GetPidBaseDir() string {
    pdMutex.RLock()
    defer pdMutex.RUnlock()

    return basePidDir
}

// ReadAsPidFile attempts to read and parse a given file path as a pid file,
// returning the pid recorded within the given file.
func ReadAsPidFile(path string) (int, error) {
    f, err := os.Open(path)
    if err != nil {
        return NoPid, err
    }

    defer f.Close()

    buf := make([]byte, 128)

    c, err := f.Read(buf)
    if err != nil {
        return NoPid, err
    }

    if c > len(buf) {
        return NoPid, PidParseErr
    }

    pid, err := strconv.ParseInt(string(buf[:c]), 10, 32)
    if err != nil {
        return NoPid, err
    }

    return int(pid), nil
}

// ReadPidFile calls ReadAsPidFile using the calling application's pid file path.
func ReadPidFile() (int, error) {
    return ReadAsPidFile(getPidFilePath())
}

// SetPidBaseDir sets the base directory for pid operations in a thread-safe way.
func SetPidBaseDir(path string) {
    pdMutex.Lock()
    defer pdMutex.Unlock()

    basePidDir = path
}

// getPidFilePath attempts to evaluate and return a standardized path to a pid file
// in the format <pid dir>/<app name>.pid.
func getPidFilePath() string {
    baseDir := GetPidBaseDir()

    return path.Join(baseDir, GetName() + ".pid")
}
