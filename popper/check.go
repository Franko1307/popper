package main

import (
	"log"
	"os"
	"strings"

	sh "github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
)

var environment []string
var volumes []string
var skip string
var timeout string

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Run experiment and report on its integrity (status)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		env := ""
		if len(environment) > 0 {
			env += " -e " + strings.Join(environment, " -e ")
		}
		if len(volumes) > 0 {
			env += " -v " + strings.Join(volumes, " -v ")
		}
		cmd_args := []string{"run", "--rm", "-i"}
		cmd_args = append(cmd_args, strings.Fields(env)...)
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		cmd_args = append(cmd_args, "--volume", dir+":"+dir, "--workdir", dir, "--volume", "/var/run/docker.sock:/var/run/docker.sock", "ivotron/poppercheck", "--timeout", timeout)
		if len(skip) > 0 {
			cmd_args = append(cmd_args, "--skip="+skip)
		}

		s := make([]interface{}, len(cmd_args))
		for i, v := range cmd_args {
			s[i] = v
		}
		if err := sh.Command("docker", s...).Run(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringSliceVarP(&environment, "environment", "e", []string{}, "Environment variables to be defined inside the test container.")
	checkCmd.Flags().StringSliceVarP(&volumes, "volume", "v", []string{}, "Volumes to be passed to the test container.")
	checkCmd.Flags().StringVarP(&skip, "skip", "s", "", "Comma-separated list of stages to skip.")
	checkCmd.Flags().StringVarP(&timeout, "timeout", "t", "36000", "Timeout limit for experiment (default: 10 hrs).")
}
