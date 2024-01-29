package cmd

import (
  "fmt"
  "os"
  "strings"

  "github.com/spf13/cobra"
)

var SemVerVersion    string
var SemVerPrerelease string
var SemVerBuildMeta  string
var BuildDate        string
var BuildVcsUrl      string
var BuildVcsId       string
var BuildVcsIdDate   string
var BuildDescription string = "Carousel - cloudfoundry-community"

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Carousel",
  Long:  `Display the version and build number for Carousel`,
  Run: func(cmd *cobra.Command, args []string) {
	  version := strings.TrimSpace(SemVerVersion)
	  prerelease := strings.TrimSpace(SemVerPrerelease)
	  buildmeta := strings.TrimSpace(SemVerBuildMeta)
	  builddate := strings.TrimSpace(BuildDate)
	  description := strings.TrimSpace(BuildDescription)
	  vcsurl := strings.TrimSpace(BuildVcsUrl)
	  vcsid := strings.TrimSpace(BuildVcsId)
	  vcsiddate := strings.TrimSpace(BuildVcsIdDate)
	  if len(prerelease) > 0 {
		  prerelease = "-" + prerelease
	  }
	  if len(buildmeta) > 0 {
		  buildmeta = "+" + buildmeta
	  }
          fmt.Fprintf( os.Stdout, "v%v%v%v %v  %v\n%v  %v   %v\n",
                       version, prerelease, buildmeta, builddate, description,
		       vcsurl, vcsid, vcsiddate)
  },
}
