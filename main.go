package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"s4/cmd"

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

	godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalf("Error loading .env file")
	// 	fmt.Println(err)
	// }

	loadDb()
	var isSessionOn bool
	var sessionUser = os.Getenv("USER-ID")
	println(sessionUser)
	if sessionUser != "" {
		isSessionOn = true
	}
	if !isSessionOn {
		username := ""
		psw := ""
		fmt.Println("Enter username:")
		fmt.Scanln(&username)

		row := database.QueryRow("SELECT psw FROM users WHERE user=?", username)

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
				_, err := database.Exec("INSERT INTO users(user, psw) VALUES(?, ?)", username, psw)
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
		os.Setenv("USER-ID", username)
		x := os.Getenv("USER-ID")
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
