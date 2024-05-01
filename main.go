package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

// func main() {
// 	fmt.Println("Welcome to the cool game.")
// 	rooms := CreateMap("testmap")
// 	world := newWorld(rooms)
// 	world.startWorld()
// 	player := newMob("Mysterious Stranger", world)
// 	var input string
// 	for world.running {
// 		fmt.Print(">")
// 		fmt.Scan(&input)
// 		availableActions := append(player.commands, player.location.getExitCommands()...)
// 		action := parseInput(input, availableActions)
// 		player.cmdQueue = append(player.cmdQueue, action)
// 	}
// }

// User represents a connected user
type User struct {
	Conn net.Conn
	// Add any additional user-related data you need to track here
	Mob *Mob
}

var (
	users     []*User
	usersLock sync.Mutex
	world     *World
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New connection established:", conn.RemoteAddr())

	// Create a new user and add it to the list
	user := &User{Conn: conn, Mob: newMob()}
	addUser(user)
	fmt.Println(user.Mob.name)
	// Send a welcome message to the user
	welcomeMessage := "Welcome to the Telnet Game!\nPlease select a name: \n"
	conn.Write([]byte(welcomeMessage))

	// Receive and process commands from the user
	for {
		// Create a buffer to hold the incoming data

		buffer := make([]byte, 1024)

		// Read data from the connection into the buffer
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			// Remove the user from the list if an error occurs
			removeUser(user)
			return
		}

		// Convert the received data to a string
		command := string(buffer[:bytesRead])
		command = strings.TrimRight(command, " \n\r")
		if user.Mob.name == "" {
			conn.Write([]byte(fmt.Sprintf("You shall be known as '%v'.\n", command)))
			output := make(chan string)
			user.Mob.connect(output)
			go func() {

				for {
					select {
					case msg := <-output:
						conn.Write([]byte(msg))
					}
				}
			}()
			user.Mob.spawn(command, world)
			continue
		}

		// Process the command
		processCommand(user, command)

	}
}

func addUser(user *User) {
	usersLock.Lock()
	defer usersLock.Unlock()
	users = append(users, user)
}

func removeUser(user *User) {
	usersLock.Lock()
	defer usersLock.Unlock()
	for i, u := range users {
		if u == user {
			// Remove the user from the list by swapping it with the last element and truncating the slice
			users[i] = users[len(users)-1]
			users = users[:len(users)-1]
			break
		}
	}
}

func processCommand(user *User, command string) {
	// Here you can implement logic to parse and handle the received command

	fmt.Println("Received command from", user.Conn.RemoteAddr(), ":", command)

	// for _, usr := range users {
	// 	usr.Conn.Write([]byte(fmt.Sprintf("%v says: %v", user.Conn.RemoteAddr(), command)))
	// }
	availableActions := append(user.Mob.commands, user.Mob.location.getExitCommands()...)
	action := parseInput(command, availableActions)
	user.Mob.cmdQueue = append(user.Mob.cmdQueue, action)
}

func main() {
	port := "0.0.0.0:8080" // Telnet default port

	rooms := CreateMap("testmap")
	world = newWorld(rooms)
	world.startWorld()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
