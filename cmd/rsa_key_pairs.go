/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// rsaKeyPairsCmd represents the rsaKeyPairs command
var rsaKeyPairsCmd = &cobra.Command{
	Use:     "rsa-key-pairs",
	Short:   "Perform actions on RSA Key Pairs",
	Long:    "RSA Key Pairs related actions.",
	Aliases: []string{"rsa"},
}

func init() {
	rootCmd.AddCommand(rsaKeyPairsCmd)
}
