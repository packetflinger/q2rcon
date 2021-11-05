# q2rcon

A command line out-of-band management client for Quake 2.

rcon password will be fetched from ~/.q2rcon

Build: go build q2rcon.go

Usage: q2rcon <server:port> <rcon string>

Examples: 
q2rcon 10.2.2.3:27910 status
q2rcon 10.2.2.3:27910 set dmflags 1023
q2rcon 10.2.2.3:27910 sv !mute cl 1 300  
