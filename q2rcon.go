package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

//
// structure to hold our config file ~/.q2servers.json
//
type ServerJSON struct {
	Passwords []struct {
		Name     string `JSON:"name"`
		Password string `JSON:"password"`
	} `JSON:"passwords"`

	Servers []struct {
		Name     string `JSON:"name"`
		Groups   string `JSON:"groups"`
		Addr     string `JSON:"addr"`
		Password string `JSON:"password"`
		SSHHost  string `JSON:"sshhost"`
	} `JSON:"servers"`
}

type Server struct {
	Name     string
	Addr     string
	Password string
}

var Config ServerJSON

// flag
var Verbose *bool

//
// Start here
//
func main() {
	serverlookup := flag.Arg(0)
	server := GetServer(serverlookup)
	stub := fmt.Sprintf("rcon %s ", server.Password)
	p := make([]byte, 1500)

	// append the port if we need to
	if !strings.Contains(server.Addr, ":") {
		server.Addr = server.Addr + ":27910"
	}

	conn, err := net.Dial("udp4", server.Addr)
	if err != nil {
		fmt.Printf("Connection error %v", err)
		return
	}
	defer conn.Close()

	// 1 second (1000ms) timeout seems short,
	// but normal clients are talking with < 100ms RTT
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	cmd := strings.Join(flag.Args()[1:], " ")

	if cmd == "" {
		Usage()
		return
	}

	rconcmd := []byte{0xff, 0xff, 0xff, 0xff}
	rconcmd = append(rconcmd, stub...)
	rconcmd = append(rconcmd, cmd...)
	fmt.Fprintln(conn, string(rconcmd))

	length, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
		return
	}

	if length > 11 {
		fmt.Println(string(p[10:]))
	}
}

func Usage() {
	fmt.Printf("Usage: %s <server alias> <rcon command>\n", os.Args[0])
}

//
// Replace some characters in the password for displaying on screen
// Half of the string is replaced with * runes
//
func ObfuscatePassword(input string) string {
	inlength := len(input)
	position := 0

	runeinput := []rune(input)
	for i, _ := range runeinput {
		if position < inlength/2 {
			runeinput[i] = '*'
			position++
		} else {
			break
		}
	}

	return string(runeinput)
}

//
// Find the password linked with a particular server
//
func GetPassword(alias string) string {
	for _, p := range Config.Passwords {
		if p.Name == alias {
			return p.Password
		}
	}

	return "couldntfinerconpassword"
}

//
// Return a struct representing a specific server
//
func GetServer(alias string) Server {
	for _, s := range Config.Servers {
		if s.Name == alias {
			srv := Server{
				Name:     s.Name,
				Addr:     s.Addr,
				Password: GetPassword(s.Password),
			}
			return srv
		}
	}

	return Server{}
}

//
// Called before main()
//
func init() {
	if len(os.Args) < 3 {
		Usage()
		return
	}

	flag.Parse()

	homedir, err := os.UserHomeDir()
	sep := os.PathSeparator
	if err != nil {
		log.Fatal(err)
	}

	configfile := fmt.Sprintf("%s%c.q2servers.json", homedir, sep)
	raw, err := os.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
