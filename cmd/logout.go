package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var logout = &cobra.Command{
	Use:   "logout",
	Short: "iuboknl;",
	Long:  "igiuinmo",

	Run: logoutFun,
}

func init() {
	rootCmd.AddCommand(logout)
}

func logoutFun(cmd *cobra.Command, args []string) {
	err := os.Setenv("USER-ID", "")
	if err != nil {
		fmt.Println("Error logging out")
	}
	fmt.Println("Logged out")
}
