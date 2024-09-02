package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	envFile := ".env"

	envMap, err := LoadEnvFile(envFile)
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	envMap["USER_ID"] = ""

	err = WriteEnvFile(envFile, envMap)
	if err != nil {
		fmt.Println("Error writing to .env file:", err)
	}

	if err != nil {
		fmt.Println("Error logging out")
	}
	fmt.Println("Logged out")
}

func LoadEnvFile(filename string) (map[string]string, error) {
	envMap := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return envMap, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) == 2 {
			envMap[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
		}
	}

	return envMap, scanner.Err()
}

func WriteEnvFile(filename string, envMap map[string]string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for key, value := range envMap {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return err
		}
	}

	return nil
}
