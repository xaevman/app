//  ---------------------------------------------------------------------------
//
//  all_test.go
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
    "fmt"
    "testing"
)

var currentPid = NoPid


// TestExePathOps tests the various functions relating to the executing binary
// to make sure they return valid strings.
func TestExePathOps(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestExePathOps ============")

    appName := GetAppName()
    if len(appName) < 1 {
        t.Error("appName len < 1")
    }
    fmt.Println("appName = " + appName)

    exeDir := GetExeDir()
    if len(exeDir) < 1 {
        t.Error("exeDir len < 1")
    }
    fmt.Println("exeDir = " + exeDir)

    exeFile := GetExeFile()
    if len(exeFile) < 1 {
        t.Error("exeFile len < 1")
    }
    fmt.Println("exeFile = " + exeFile)

    exePath := GetExePath()
    if len(exePath) < 1 {
        t.Error("exePath len < 1")
    }
    fmt.Println("exePath = " + exePath)

    fmt.Println("Ok")
}

// TestRunStatusEarly tests the run status query when a PID file has not yet
// been created.
func TestRunStatusEarly(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestRunStatusEarly ============")

    proc := GetRunStatus()

    // should be nil at this point
    if proc != nil {
        fmt.Printf("%v\n", proc)
        t.Error("should be no existing process by this name")
    }

    fmt.Println("Ok")
}

// TestPidCreate tests the creation of the application's pid file.
func TestPidCreate(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestPidCreate ============")

    pid, err := CreatePidFile()
    if err != nil {
        t.Error(err)
    }

    currentPid = pid

    fmt.Printf("Pid: %d\n", currentPid)

    fmt.Println("Ok")
}

// TestReadPid tests reading the pid from the application's pid file.
func TestReadPid(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestReadPid ============")

    pid, err := ReadPidFile()
    if err != nil {
        t.Error(err)
    }   

    if pid != currentPid {
        t.Errorf("pids don't match: (current)%d != (read)%d", currentPid, pid)
    }

    fmt.Println("Ok")
}

// TestRunStatusLate tests the run status query when the application's pid file
// exists.
func TestRunStatusLate(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestRunStatusLate ============")

    proc := GetRunStatus()

    // now we should match an active process
    if proc == nil {
        t.Error("couldn't find test app pid")
    }

    fmt.Println("Ok")
}

// TestPidDelete deletes cleanup of the application's pid file.
func TestPidDelete(t *testing.T) {
    fmt.Println()
    fmt.Println("============ TestPidDelete ============")

    err := DeletePidFile()
    if err != nil {
        t.Error(err)
    }

    fmt.Println("Ok")
}
