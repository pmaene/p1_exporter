package main

import (
	"github.com/pmaene/p1_exporter/cmd"
)

var (
	version = ""
	commit  = ""
)

func main() {
	cmd.Execute()
}

func init() {
	cmd.SetVersion(version)
	cmd.SetCommit(commit)
}
