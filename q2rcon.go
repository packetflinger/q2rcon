package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"
	"time"
)

func main() {
	server := os.Args[1]
	cmd := strings.Join(os.Args[2:], " ")

	password := LoadRCONPassword()
	stub := fmt.Sprintf("rcon %s ", password)

	p := make([]byte, 1500)

	// only use IPv4
	conn, err := net.Dial("udp4", server)
	if err != nil {
		fmt.Printf("Connection error %v", err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	rconcmd := []byte{0xff, 0xff, 0xff, 0xff}
	rconcmd = append(rconcmd, stub...)
	rconcmd = append(rconcmd, cmd...)
	fmt.Fprintln(conn, string(rconcmd))

	length, err := bufio.NewReader(conn).Read(p)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	if length > 11 {
		fmt.Println(string(p[10:]))
	}
}

/**
 * Look in ~/.q2rcon for a password
 */
func LoadRCONPassword() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	pwfile := fmt.Sprintf("%s%c%s", user.HomeDir, os.PathSeparator, ".q2rcon")
	pwdata, err := os.ReadFile(pwfile)
	if err != nil {
		panic(err)
	}

	// only match "=" as a delimiter
	delim := func(c rune) bool {
		return c == '='
	}

	passwordfields := strings.FieldsFunc(string(pwdata), delim)
	if len(passwordfields) < 2 {
		panic("Invalid rcon password file")
	}

	// remove common whitespace from both sides
	password := strings.Trim(passwordfields[1], " \n\t")
	return password
}
