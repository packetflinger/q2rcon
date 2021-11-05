# Q2rcon
A command line out-of-band management client for Quake 2.

## Config File
Add your rcon password to `~/.q2rcon`

## Usage
`q2rcon ip:port consolecmd`

## Examples
```
q2rcon 192.0.2.45:27910 status
q2rcon 192.0.2.10:27913 set dmflags 1023
q2rcon 192.0.2.11:27910 sv !mute cl 2 300
```

