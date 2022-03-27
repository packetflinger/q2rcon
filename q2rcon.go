package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		Usage()
		return
	}

	passwordfile := flag.String("config", ".q2rcon", "The file containing the rcon password")

	flag.Parse()
	server := flag.Arg(0)
	cmd := strings.Join(flag.Args()[1:], " ")

	if cmd == "" {
		Usage()
		return
	}

	password := LoadRCONPassword(*passwordfile)
	fmt.Printf("password: %s\n", password)

	stub := fmt.Sprintf("rcon %s ", password)
	p := make([]byte, 1500)

	if !strings.Contains(server, ":") {
		server = server + ":27910"
	}
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
 * Look in rconfile for a password
 */
func LoadRCONPassword(rconfile string) string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	pwfile := fmt.Sprintf("%s%c%s", user.HomeDir, os.PathSeparator, rconfile)
	pwdata, err := os.ReadFile(pwfile)
	if err != nil {
		// problems with rcon file, try environment variable
		pw := os.Getenv("RCON")
		if pw == "" {
			panic(err)
		}
		return pw
	}

	return strings.Trim(string(pwdata), " \n\t")
}

func Usage() {
	fmt.Printf("Usage: %s [-config=rconfile] <serverip[:port]> <rcon command>\n", os.Args[0])
	return
}
