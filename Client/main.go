package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// Read username and password from user
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')

	// Send username and password to server
	conn.Write([]byte(fmt.Sprintf("%s|%s", strings.TrimSpace(username), strings.TrimSpace(password))))

	// Read server response
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(response)

	if strings.Contains(response, "Invalid") {
		return
	}

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Deposit")
		fmt.Println("2. Withdraw")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Print("Enter amount to deposit: $")
			depositAmount, _ := reader.ReadString('\n')
			conn.Write([]byte("deposit\n"))
			conn.Write([]byte(depositAmount))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print(response)
		case "2":
			fmt.Print("Enter amount to withdraw: $")
			withdrawAmount, _ := reader.ReadString('\n')
			conn.Write([]byte("withdraw\n"))
			conn.Write([]byte(withdrawAmount))
			response, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print(response)
		case "3":
			conn.Write([]byte("exit\n"))
			fmt.Println("Thank you for using the banking system. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
