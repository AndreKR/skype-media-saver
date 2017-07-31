package main

import (
	"github.com/mitchellh/go-ps"
	"os"
	"path/filepath"
)

func findMyOtherProcesses() (pids []int) {

	processes, err := ps.Processes()
	if err != nil {
		log.Fatalln("Error finding process: " + err.Error())
	}

	for _, p := range processes {
		if p.Pid() == os.Getpid() {
			continue
		}
		if p.Executable() == filepath.Base(os.Args[0]) {
			pids = append(pids, p.Pid())
		}
	}

	return
}
