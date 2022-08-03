package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"
	"time"
)

const (
	PasswordFile string = ".q2rcon"
	AliasFile    string = ".q2info"
	PasswordENV  string = "RCON" // environment variable
)

//
// Structure for rcon passwords
//
type JSONData struct {
	Profiles []struct {
		Name     string `JSON:"name"`
		Password string `JSON:"password"`
		Default  bool   `JSON:"default"`
	} `JSON:"profiles"`
}

// flag
var Verbose *bool

//
// Start here
//
func main() {
	if len(os.Args) < 3 {
		Usage()
		return
	}

	Verbose = flag.Bool("v", false, "Show verbose output")
	passwordfile := flag.String("config", PasswordFile, "The file containing the rcon password")
	profile := flag.String("p", "", "The rcon password to use")

	flag.Parse()
	server := flag.Arg(0)
	cmd := strings.Join(flag.Args()[1:], " ")

	if cmd == "" {
		Usage()
		return
	}

	dirname, err := os.UserHomeDir()
	if err == nil {
		aliasfile := fmt.Sprintf("%s/%s", dirname, AliasFile)
		server = GetAlias(aliasfile, server)
		if *Verbose {
			fmt.Printf("** looking up servername in alias file %s: %s\n", aliasfile, server)
		}
	}

	password := LoadRCONPassword(*passwordfile, *profile)
	if password == "" {
		fmt.Println("Unable to locate a valid rcon password")
		return
	}

	if *Verbose {
		fmt.Printf("** using rcon password %s\n", ObfuscatePassword(password))
	}

	stub := fmt.Sprintf("rcon %s ", password)
	p := make([]byte, 1500)

	// append the port if we need to
	if !strings.Contains(server, ":") {
		server = server + ":27910"
	}

	conn, err := net.Dial("udp4", server)
	if err != nil {
		fmt.Printf("Connection error %v", err)
		return
	}
	defer conn.Close()

	// 1 second (1000ms) timeout seems short,
	// but normal clients are talking with < 100ms RTT
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

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

/**
 * Look in rconfile for a password
 */
func LoadRCONPassword(rconfile string, lookup string) string {
	safelookup := strings.Trim(lookup, " ")
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	pwfile := fmt.Sprintf("%s%c%s", user.HomeDir, os.PathSeparator, rconfile)
	pwdata, err := os.ReadFile(pwfile)
	if err != nil {
		if *Verbose {
			fmt.Printf("** password file %s not found\n", pwfile)
			fmt.Printf("** trying environment variable %s\n", PasswordENV)
		}

		// problems with rcon file, try environment variable
		pw := os.Getenv(PasswordENV)
		if pw == "" {
			panic(err)
		}
		return pw
	}

	if *Verbose {
		fmt.Printf("** reading password file %s\n", pwfile)
	}

	var profiles JSONData

	if err = json.Unmarshal(pwdata, &profiles); err != nil {
		fmt.Printf("error parsing %s: %s\n", pwfile, err.Error())
	}

	for _, p := range profiles.Profiles {
		if safelookup == "" && p.Default {
			return p.Password
		}

		if p.Name == safelookup {
			return p.Password
		}
	}
	return ""
}

func Usage() {
	fmt.Printf("Usage: %s [-v] [-p password_profile] <(serverip[:port])|(alias)> <rcon command>\n", os.Args[0])
}

/**
 * Find lookup in aliasfile
 */
func GetAlias(aliasfile string, lookup string) string {
	raw, err := os.ReadFile(aliasfile)
	if err != nil {
		return lookup
	}

	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		trimmedline := strings.TrimSpace(line)
		if trimmedline == "" {
			continue
		}

		if strings.HasPrefix(trimmedline, "#") {
			continue
		}

		if strings.HasPrefix(trimmedline, "//") {
			continue
		}

		alias := strings.Fields(line)

		if alias[0] == lookup {
			return alias[1]
		}
	}

	return lookup
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
