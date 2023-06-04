// Package compileinfo provides output to stdout about version, build and commit at the build time
package compileinfo

import "fmt"

// PrintCompileInfo prints state at the build time
func PrintCompileInfo(buildVersion string, buildDate string, buildCommit string) {
	// вывод информации о компиляции
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}
