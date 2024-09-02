package cmd

import (
	"fmt"
	"math/rand"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "iuboknl;",
	Long:  "igiuinmo",

	Run: genPsw,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().IntP("length", "l", 8, "tfg8nyuhi")
	generateCmd.Flags().BoolP("digits", "d", false, "tfg8nykjnuhi")
	generateCmd.Flags().BoolP("spChar", "s", false, "tfijpog8nyuhi")

}

func genPsw(cmd *cobra.Command, args []string) {
	length, _ := cmd.Flags().GetInt("length")
	isDigits, _ := cmd.Flags().GetBool("digits")
	isspChar, _ := cmd.Flags().GetBool("spChar")

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if isDigits {
		charset += "0123456789"
	}
	if isspChar {
		charset += "!@#$%^&*()_+"
	}
	password := ""

	for i := 0; i < length; i++ {
		password += string(charset[rand.Intn(len(charset))])
	}

	fmt.Println(password)
}
