# Q2rcon
A command line out-of-band management client for Quake 2.

## Config File
Add your rcon password to default file `~/.q2rcon`. You can
specify any file in your home directory with the `-config` option.

## Usage
`q2rcon [-config=<file>] ip:[port] consolecmd`

The port is optional, the default is `27910`

## Examples
```
q2rcon 192.0.2.45 status
q2rcon 192.0.2.10:27913 set dmflags 1023
q2rcon 192.0.2.11:27910 sv !mute cl 2 300
q2rcon -config=pfservers.txt 192.0.2.4:27999 recycle
```
