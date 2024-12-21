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
func Execute(version string, commit string, date string) {
	fmt.Printf("start ifc version - %s, commit - %s,date - %s\n", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&utils.ConnectionFile1,
		"connection1",
		"c",
		"",
		"file connection to HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.Org1,
		"org1",
		"o",
		"",
		"organization HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.User1,
		"user1",
		"u",
		"",
		"user HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.ConnectionFile2,
		"connection2",
		"d",
		"",
		"file connection to HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.Org2,
		"org2",
		"p",
		"",
		"organization HLF",
	)

	rootCmd.PersistentFlags().StringVarP(
		&utils.User2,
		"user2",
		"v",
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
