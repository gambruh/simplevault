package main

import (
	"log"

	"github.com/gambruh/simplevault/internal/compileinfo"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	compileinfo.PrintCompileInfo(buildVersion, buildDate, buildCommit)

	if err := newRootCmd(runClient).Execute(); err != nil {
		log.Fatal(err)
	}
}
