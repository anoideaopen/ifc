package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anoideaopen/ifc/proc"
	"github.com/anoideaopen/ifc/utils"
	"github.com/spf13/cobra"
)

var workBlockCmd = &cobra.Command{
	Use:   "work",
	Short: "start a work process",
	Run: func(cmd *cobra.Command, args []string) {
		if utils.ConnectionFile1 == "" || utils.Org1 == "" || utils.User1 == "" ||
			utils.ConnectionFile2 == "" || utils.Org2 == "" || utils.User2 == "" {
			panic(utils.ErrorNotFindHLFConf)
		}

		var (
			terminate   = make(chan os.Signal, 1)
			ctx, cancel = context.WithCancel(context.Background())
		)

		// Start processing
		go func() {
			signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
			<-terminate
			cancel()
		}()

		proc.Process(ctx,
			utils.ConnectionFile1, utils.Org1, utils.User1,
			utils.ConnectionFile2, utils.Org2, utils.User2,
		)
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(workBlockCmd)
}
