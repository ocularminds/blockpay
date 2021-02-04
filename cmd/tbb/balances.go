package main

import(
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"github.com/ocularminds/blockpay/database"
)

func balancesCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()
	
			fmt.Println("Accounts balances:")
			fmt.Println("__________________")
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account, balance))
			}
		},
	}
    var balancesCmd = &cobra.Command{
        Use:   "balances",
        Short: "Interact with balances (list...).",
        PreRunE: func(cmd *cobra.Command, args []string) error {
            return incorrectUsageErr()
        },
        Run: func(cmd *cobra.Command, args []string) {
        },
    }

    balancesCmd.AddCommand(balancesListCmd)
    return balancesCmd
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}