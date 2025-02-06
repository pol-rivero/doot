package cmd

import (
	"github.com/pol-rivero/doot/lib/common/log"
	"github.com/spf13/cobra"
)

func SetUpLogger(cmd *cobra.Command) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		panic(err)
	}
	quiet, err := cmd.Flags().GetBool("quiet")
	if err != nil {
		panic(err)
	}

	log.Init(verbose, quiet)
}
