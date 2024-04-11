package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"os"
	"encoding/json"
	"time"
	"io"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	FullName  string   `json:"fullname"`
	Emails    []string `json:"emails"`
	Addresses []string `json:"addresses"`
}

var users []User
var clientKeys map[*net.Conn]string

func main() {
	fmt.Println("Server is running...")


	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	

	// Load users from JSON file
	if err := loadUsers("users.json"); err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	// Initialize client keys map
	clientKeys = make(map[*net.Conn]string)


	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}

func loadUsers(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return err
	}

	return nil
}

func saveUsers(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		return err
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var key string
	// Authentication
	authenticated := false
	for authenticated == false {

		username, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		username = strings.TrimSpace(username)

		password, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		password = strings.TrimSpace(password)

		for _, u := range users {
			if u.Username == username && u.Password == password {
				authenticated = true
				break
			}
		}

		if authenticated == false {
			fmt.Fprintf(conn, "Invalid credentials. Please try again.\n")
		} else {

			// Generate a unique key for the authenticated user
			rand.Seed(time.Now().UnixNano())
			key = strconv.Itoa(rand.Intn(1000000))
			// Send successful signal
			fmt.Fprintf(conn, "successful\n")
			// Send the key to the client
			clientRes, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(clientRes)
			fmt.Fprintf(conn, "%s\n", key)
			//clientKeys[&c] = key
			break
		}
	}

	for {
		conn.Write([]byte("Enter 1 to play game, 2 to download user file, 3 to quit\n"))
		menuNum := getMenuChoice(conn)

		switch menuNum {
		case 1:
			handleGameSession(conn, key)
		case 2:
			handleFileDownload(conn, "users.json")

		case 3:
			fmt.Fprintf(conn, "%s_Stop\n", key)
			fmt.Println(key + " has quit")
			return
		default:
			continue
		}
	}
}

// func handleFileDownload(conn net.Conn, filename string) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Fprintf(conn, "Error opening file: %s\n", err)
// 		return
// 	}
// 	defer file.Close()

// 	// Send file name to client
// 	//fmt.Fprintf(conn, "%s\n", filename)

// 	// Send the file contents to the client

// 	fmt.Fprintf(conn, "File sent successfully.\n")
// 	// Send file name to client
// 	fmt.Fprintf(conn, "%s\n", filename)
	
// 	_, err = io.Copy(conn, file)
// 	if err != nil {
// 		fmt.Fprintf(conn, "Error sending file: %s\n", err)
// 		return
// 	}
// }

func handleFileDownload(conn net.Conn, filename string) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(conn, "Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	// Send file name to client
	//fmt.Fprintf(conn, "%s\n", filename)

	// Prompt the client to provide the desired file path

	//C:\Users\Asus\Downloads\test\downloaded_file.txt
	//conn.Write([]byte("Please enter the path to save the file: "))
	fmt.Fprintf(conn, "Please enter the path to save the file: \n")
	path, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Fprintf(conn, "Error reading file path: %s\n", err)
		return
	}
	path = strings.TrimSpace(path)

	// Send file name to client
	//fmt.Fprintf(conn, "%s\n", filename)

	// Create the file at the provided path
	newFile, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(conn, "Error creating file: %s\n", err)
		return
	}
	defer newFile.Close()

	fmt.Fprintf(conn, "File saved successfully at %s. Enter to return menu\n", path)

	// Send the file contents to the client
	_, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Fprintf(conn, "Error sending file: %s\n", err)
		return
	}

	
}

func getMenuChoice(conn net.Conn) int {
	reader := bufio.NewReader(conn)
	menu, _ := reader.ReadString('\n')
	menuNum, _ := strconv.Atoi(strings.TrimSpace(menu))
	return menuNum
}

func handleGameSession(conn net.Conn, key string) {
	for {
		secretNumber := rand.Intn(100) + 1
		fmt.Println("Target number: ", secretNumber)

		conn.Write([]byte("Welcome to the Guessing Game! Guess a number between 1 and 100\n"))

		for {
			guessStr := readInput(conn)

			if strings.TrimSpace(guessStr) == "stop" {
				fmt.Println(key + " requested to stop the game.")
				return
			}

			guess, err := strconv.Atoi(strings.TrimSpace(guessStr))
			if err != nil {
				conn.Write([]byte("Invalid input. Please enter a number.\n"))
				continue
			}

			if guess < secretNumber {
				// conn.Write([]byte("Too low! Try again.\n"))
				fmt.Fprintf(conn, "%s_To Low\n", key)
			} else if guess > secretNumber {
				fmt.Fprintf(conn, "%s_To Hight\n", key)
			} else {
				fmt.Fprintf(conn, "%s_Congratulations! You guessed it! Enter 'next' to play again or 'quit' to exit.\n", key)
				choose := readInput(conn)
				if strings.TrimSpace(choose) == "next" {
					break
				}
				return
			}
		}
	}
}

func readInput(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	input, _ := reader.ReadString('\n')
	return input
}
