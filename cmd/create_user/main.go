package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/OutClimb/Redirect/internal/store"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	username := ""
	password := ""
	name := ""
	email := ""

	flag.StringVar(&username, "u", "", "Specify username")
	flag.StringVar(&password, "p", "", "Specify password")
	flag.StringVar(&name, "n", "", "Specify name")
	flag.StringVar(&email, "e", "", "Specify email")
	flag.Parse()

	if username == "" || password == "" || name == "" || email == "" {
		fmt.Println("Please provide username, password, name and email")
		return
	}

	cost, err := strconv.Atoi(os.Getenv("PASSWORD_COST"))
	if err != nil {
		fmt.Println("Invalid PASSWORD_COST")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	storeLayer := store.New()
	user, err := storeLayer.CreateUser(username, string(hashedPassword), name, email, "admin")
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return
	}

	fmt.Println("User created successfully")
	fmt.Println("ID: ", user.ID)
	fmt.Println("Username: ", user.Username)
	fmt.Println("Role: ", user.Role)
	fmt.Println("Name: ", user.Name)
	fmt.Println("Email: ", user.Email)
}
