package main

import "github.com/anoideaopen/ifc/cmd"

var (
	version = "none"
	commit  = "none"
	date    = "none"
)

func main() {
	cmd.Execute(version, commit, date)
}
