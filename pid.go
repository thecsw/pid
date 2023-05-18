package pid

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/thecsw/rei"
)

const (
	// pidPath is the path where the pid files are stored.
	pidPath = "/tmp"
)

type proc struct {
	pid int
	loc string
}

// Stop deletes the pid file session.
func (p proc) Stop() {
	// we need to remove the pid file that we created
	if err := os.Remove(p.loc); err != nil {
		// not fatal, at least log it
		log.Printf("removing pid file %s: %v\n", p.loc, err)
	}
}

// Start starts a new pid session.
// The caller should call Stop() to cleanup pid data.
func Start(name string) interface{ Stop() } {
	// we don't run this on windows, because we can't check if the pid exists via a file
	// TODO: support windows, because we can check their pids
	if runtime.GOOS == "windows" {
		log.Fatal("windows pids are not supported")
	}

	// get the pid and put the file in
	pid := os.Getpid()

	// get the location of the pid file
	loc := filepath.Clean(pidPath + "/" + name + ".pid")

	// check if the pid file exists
	pidFileExists, err := rei.FileExists(loc)
	if err != nil {
		log.Fatalf("checking pid file existence: %v", err)
	}

	// if it exists, try to quickly check if there is actually a pid running
	if pidFileExists {
		// read the pid file
		targetPidData, err := os.ReadFile(filepath.Clean(loc))
		if err != nil {
			log.Fatalf("reading pid file %s: %v", loc, err)
		}
		// convert the pid file to an int
		targetPid, err := rei.Atoi(string(targetPidData))
		if err != nil {
			log.Fatalf("value in %s is not a pid: %v", loc, err)
		}
		// check if the pid exists
		pidProcessExists, err := process.PidExists(int32(targetPid))
		if err != nil {
			log.Fatalf("checking for pid %d existence: %v", targetPid, err)
		}
		// if the pid exists, we can't start
		if pidProcessExists {
			log.Fatalf("your app with pid %d is already running", targetPid)
		}
		// if the pid doesn't exist, we can delete the file and continue
		if err := os.Remove(loc); err != nil {
			log.Fatalf("removing stale pid file %s: %v", loc, err)
		}
	}

	// if the file doesn't exist or stale pid deleted, create one and leave this scope
	if err := os.WriteFile(loc, []byte(rei.Itoa(pid)), 0644); err != nil {
		log.Fatalf("writing pid file %s: %v", loc, err)
	}

	// return the process
	return proc{
		pid: pid,
		loc: loc,
	}
}
