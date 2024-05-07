package main

import (
	"fmt"
	"net"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

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
	log.Info("New connection established from ", conn.RemoteAddr())

	// Create a new user and add it to the list
	user := &User{Conn: conn, Mob: newMob()}
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
			log.WithError(err).Warn("Error reading from connection.")
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
			addUser(user)
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

func getUserFromMob(m *Mob) (*User, error) {
	for _, u := range users {
		if u.Mob == m {
			return u, nil
		}
	}
	return nil, fmt.Errorf("Could not find user with mob '%v' in connection pool.", m)
}

func disconnectUserFromMob(m *Mob) {
	if user, err := getUserFromMob(m); err == nil {
		removeUser(user)
		user.Conn.Close()
	} else {
		log.WithError(err).Errorf("Could not find mob '%v'.", m)
	}
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

	log.WithFields(log.Fields{
		"mob_name":       user.Mob.name,
		"location":       user.Mob.location.name,
		"command":        command,
		"remote_address": user.Conn.RemoteAddr(),
	}).Info("Command received")

	firstPart, otherParts, _ := strings.Cut(command, " ")
	availableActions := append(user.Mob.commands, user.Mob.location.getExitCommands()...)
	action := readyCommand(firstPart, availableActions)
	user.Mob.cmdQueue = append(user.Mob.cmdQueue, action(user.Mob, otherParts))
}

func main() {
	port := "0.0.0.0:8080" // Telnet default port

	rooms := CreateMap("testmap")
	world = newWorld(rooms)
	world.startWorld()
	defer world.stopWorld()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.WithError(err).Fatal("Error listening on port", port)
		return
	}
	defer listener.Close()
	log.Infof("Server is listening on port %v", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.WithError(err).Fatal("Error accepting connection.")
			continue
		}
		go handleConnection(conn)
	}
}
