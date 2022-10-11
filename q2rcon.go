package main

import (
	"bufio"
	"encoding/json"
	"errors"
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

//
// Temp server structure. All info needed for sending rcon msgs
//
type Server struct {
	Name     string
	Addr     string
	Password string
}

var (
	Config  ServerJSON
	SrvFile *string // flag
	Verbose *bool   // flag
)

//
// Start here
//
func main() {
	serverlookup := flag.Arg(0)
	server, err := GetServer(serverlookup)
	if err != nil {
		fmt.Println(err)
		return
	}

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

	rconcmd := []byte{0xff, 0xff, 0xff, 0xff}
	rconcmd = append(rconcmd, stub...)
	rconcmd = append(rconcmd, cmd...)
	timestart := time.Now()
	fmt.Fprintln(conn, string(rconcmd))

	length, err := bufio.NewReader(conn).Read(p)
	duration := time.Since(timestart)
	if err != nil {
		fmt.Printf("Read error: %s\n", err)
		return
	}

	if length > 11 {
		fmt.Println(string(p[10:]))
	}

	if *Verbose {
		log.Printf("Results fetched in %s\n", duration.String())
	}
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
func GetPassword(alias string) (string, error) {
	for _, p := range Config.Passwords {
		if p.Name == alias {
			return p.Password, nil
		}
	}

	return "", errors.New("Couldn't locate password tagged as " + alias)
}

//
// Return a struct representing a specific server
//
func GetServer(alias string) (Server, error) {
	for _, s := range Config.Servers {
		if s.Name == alias {
			actualpassword, err := GetPassword(s.Password)
			if err != nil {
				return Server{}, err
			}

			srv := Server{
				Name:     s.Name,
				Addr:     s.Addr,
				Password: actualpassword,
			}

			if *Verbose {
				log.Printf("Querying %s [%s] using password %s\n",
					srv.Name,
					srv.Addr,
					ObfuscatePassword(srv.Password),
				)
			}
			return srv, nil
		}
	}

	return Server{}, errors.New("unknown server")
}

//
// Called before main()
//
func init() {

	// parse args
	SrvFile = flag.String("c", "", "Specify a server data file")
	Verbose = flag.Bool("v", false, "Show some more info")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Printf("Usage: %s <flags> <serveralias> <rcon command>\n", os.Args[0])
		fmt.Printf("  flags:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *SrvFile == "" {
		homedir, err := os.UserHomeDir()
		sep := os.PathSeparator
		if err != nil {
			log.Fatal(err)
		}
		*SrvFile = fmt.Sprintf("%s%c.q2servers.json", homedir, sep)
	}

	if *Verbose {
		log.Printf("Loading passwords/servers from %s\n", *SrvFile)
	}

	raw, err := os.ReadFile(*SrvFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(raw, &Config)
	if err != nil {
		log.Fatal(err)
	}

	if *Verbose {
		log.Printf("  %d passwords\n", len(Config.Passwords))
		log.Printf("  %d servers\n", len(Config.Servers))
	}
}
