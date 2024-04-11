package main

import (
	"net"
	"os"
	"bufio"
	"fmt"
	"strings"
	"io"
)

const (
	HOST = "localhost"
	PORT = "8081"
	TYPE = "tcp"
)

func main() {

	connection, err := net.Dial(TYPE, HOST+":"+PORT)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)

	var key string
	for {
		// Authentication
		fmt.Print("Enter username: ")
		username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(connection, username)

		fmt.Print("Enter password: ")
		password, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprintf(connection, password)

		// Read the response from the server
		response, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Print("Server connection: " + response)
		connection.Write([]byte("Client Connected\n"))
		if strings.Contains(response, "successful") {
			key, err = bufio.NewReader(connection).ReadString('\n')
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("You are now connected to the server. Your key is %s \n", key)
			break
		}

	}

	// fmt.Print("Enter username: ")
	// username, _ := reader.ReadString('\n')
	// username = strings.TrimSpace(username)

	// fmt.Print("Enter password: ")
	// password, _ := reader.ReadString('\n')
	// password = strings.TrimSpace(password)


	//test to write data
	// _, err = connection.Write([]byte("Hello Server! Greetings.\n"))

	// if err != nil {
	// 	println("Write data failed:", err.Error())
	// 	os.Exit(1)
	// }

	message, _ := bufio.NewReader(connection).ReadString('\n')
	fmt.Print(message)


	for {

		// reader := bufio.NewReader(os.Stdin)
		//write the data
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(connection, text+"\n")

		
		key = strings.TrimSpace(key)

		//get time
		result, _ := bufio.NewReader(connection).ReadString('\n')
		if strings.TrimSpace(string(result)) == key+"_"+"Stop" {
			fmt.Println("TCP client exiting...")
			return
		}else if strings.TrimSpace(string(result)) == "File sent successfully." {
			fileName, _ := bufio.NewReader(connection).ReadString('\n')
			fileName = strings.TrimSpace(fileName)
			saveFile(connection, fileName)

			result, _ := bufio.NewReader(connection).ReadString('\n')
			fmt.Print(result)
			
		} else {
			fmt.Print(result)
		}


	}

}

func saveFile(conn net.Conn, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Lỗi khi tạo file:", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, conn)
	if err != nil {
		fmt.Println("Lỗi khi lưu file:", err)
		return
	}
	fmt.Printf("File đã được lưu thành công với tên %s\n", fileName)
}
