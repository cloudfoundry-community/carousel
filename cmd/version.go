package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Carousel",
  Long:  `Display the version and build number for Carousel`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("v0.8.0 Carousel - cloudfoundry-community")
  },
}
