package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// User represents a bank customer
type User struct {
	Username string
	Password string
	Balance  float64
}

// BankingServer represents the banking server
type BankingServer struct {
	users map[string]User
}

// NewBankingServer creates a new instance of a BankingServer
func NewBankingServer() *BankingServer {
	return &BankingServer{
		users: make(map[string]User),
	}
}

// AddUser adds a new user to the banking server
func (bs *BankingServer) AddUser(username, password string, balance float64) {
	bs.users[username] = User{Username: username, Password: password, Balance: balance}
}

// Login authenticates a user
func (bs *BankingServer) Login(username, password string) (*User, bool) {
	user, ok := bs.users[username]
	if !ok || user.Password != password {
		return nil, false
	}
	return &user, true
}

// Deposit adds funds to a user's account
func (bs *BankingServer) Deposit(user *User, amount float64) {
	user.Balance += amount
}

// Withdraw deducts funds from a user's account
func (bs *BankingServer) Withdraw(user *User, amount float64) bool {
	if user.Balance < amount {
		return false
	}
	user.Balance -= amount
	return true
}

func handleConnection(conn net.Conn, server *BankingServer) {
	defer conn.Close()

	// Read username and password from client
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	data := strings.Split(string(buf[:n]), "|")
	username := data[0]
	password := data[1]

	// Authenticate user
	user, loggedIn := server.Login(username, password)
	if !loggedIn {
		conn.Write([]byte("Invalid username or password"))
		return
	}

	// Send welcome message and balance to client
	welcomeMsg := fmt.Sprintf("Welcome, %s! Your current balance is: $%.2f\n", user.Username, user.Balance)
	conn.Write([]byte(welcomeMsg))

	// Handle client requests
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			return
		}
		choice := strings.TrimSpace(string(buf[:n]))

		switch choice {
		case "deposit":
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			depositAmount, _ := strconv.ParseFloat(strings.TrimSpace(string(buf[:n])), 64)
			server.Deposit(user, depositAmount)
			conn.Write([]byte(fmt.Sprintf("Deposit successful. Your new balance is: $%.2f\n", user.Balance)))
		case "withdraw":
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				return
			}
			withdrawAmount, _ := strconv.ParseFloat(strings.TrimSpace(string(buf[:n])), 64)
			if success := server.Withdraw(user, withdrawAmount); success {
				conn.Write([]byte(fmt.Sprintf("Withdrawal successful. Your new balance is: $%.2f\n", user.Balance)))
			} else {
				conn.Write([]byte("Insufficient funds"))
			}
		case "exit":
			return
		default:
			conn.Write([]byte("Invalid choice. Please try again.\n"))
		}
	}
}

func main() {
	// Create a new banking server instance
	server := NewBankingServer()

	// Add some users to the banking server
	server.AddUser("user1", "password1", 1000.0)
	server.AddUser("user2", "password2", 500.0)

	// Start the server
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Banking server started. Listening on localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}
		go handleConnection(conn, server)
	}
}
