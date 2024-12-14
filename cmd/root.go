package cmd

import (
	"fmt"
	"os"

	"github.com/anoideaopen/ifc/utils"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "ifc",
	Short: "interfabric communication",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&utils.ConnectionFile1,
		"connection1",
		"c1",
		"",
		"file connection to HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.Org1,
		"org1",
		"o1",
		"",
		"organization HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.User1,
		"user1",
		"u1",
		"",
		"user HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.ConnectionFile2,
		"connection2",
		"c2",
		"",
		"file connection to HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.Org2,
		"org2",
		"o2",
		"",
		"organization HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.User2,
		"user2",
		"u2",
		"",
		"user HLF",
	)

	_ = rootCmd.MarkPersistentFlagRequired(utils.ConnectionFile1)
	_ = rootCmd.MarkPersistentFlagRequired(utils.Org1)
	_ = rootCmd.MarkPersistentFlagRequired(utils.User1)

	_ = rootCmd.MarkPersistentFlagRequired(utils.ConnectionFile2)
	_ = rootCmd.MarkPersistentFlagRequired(utils.Org2)
	_ = rootCmd.MarkPersistentFlagRequired(utils.User2)
}
