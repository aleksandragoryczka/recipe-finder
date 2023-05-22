package commandlinearguments

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "recipeFinder",
	Short: "recipeFinder - program generating meals from ingredients in your fridge",
	Long: `recipeFinder is a program that helps people avoid wasting food 
			while maintaining a healthy lifestyle by generating a list of meals
			that can be prepared from ingredients in your fridge 
			with minimal number of missing ingredients`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, "Error while executing your CLI: '%s", err)
		os.Exit(1)
	}
}
