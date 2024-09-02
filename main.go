package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"s4/cmd"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

var database *sql.DB

func loadDb() {
	var err error
	database, err = sql.Open("sqlite3", "./database.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}

	loadDb()
	var isSessionOn bool
	var sessionUser = os.Getenv("USER_ID")
	if sessionUser != "" {
		isSessionOn = true
	}
	if !isSessionOn {
		username := ""
		psw := ""
		fmt.Println("Enter username:")
		fmt.Scanln(&username)

		row := database.QueryRow("SELECT psw FROM users WHERE username=?", username)

		var q_psw string
		err := row.Scan(&q_psw)
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			fmt.Println("Register ? [Y/N]")
			var choice string
			fmt.Scanln(&choice)
			if choice == "N" {
				return
			}
			if choice == "Y" {
				fmt.Println("Enter username:")
				fmt.Scanln(&username)
				fmt.Println("password:")
				psw = newPassword()
				fmt.Println(psw)
				_, err := database.Exec("INSERT INTO users(username, psw) VALUES(?, ?)", username, psw)
				if err != nil {
					log.Fatalln(err)
				}
				fmt.Println("User registered successfully")
				return
			}
		} else {
			fmt.Println("Enter password:")
			fmt.Scanln(&psw)
			if psw != q_psw {
				fmt.Println("Invalid password")
				return
			}
		}

		envFile := ".env"

		envMap, err := LoadEnvFile(envFile)
		if err != nil {
			fmt.Println("Error loading .env file:", err)
			return
		}

		envMap["USER_ID"] = username

		err = WriteEnvFile(envFile, envMap)
		if err != nil {
			fmt.Println("Error writing to .env file:", err)
		}

		x := os.Getenv("USER_ID")
		fmt.Println("Welcome", x)
	}

	cmd.Execute()
}

func newPassword() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
	password := ""

	for i := 0; i < 8; i++ {
		password += string(charset[rand.Intn(len(charset))])
	}
	return password
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
