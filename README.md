# Q2rcon
A command line out-of-band management client for Quake 2.

## Config File
Servers, passwords and mappings are done via the config file `.q2servers.config`
in your home directory. This file is used by other programs as well 
so it might contain more information than is strictly needed for this
utility.

The config file is a text-format protobuf. The proto it implements is
https://github.com/packetflinger/libq2/blob/main/proto/servers_file.proto

You can have as many `password` and `server` stanzas as you like. Example:
```
password {
  identifier: "passwordname1"
  secret: "the-actual-rcon-password"
}

server {
  identifier: "server1"
  address: "192.0.2.55:27911"
  rcon_password: "passwordname1"
}
```

The example above shows 1 server and 1 password. Passwords
are defined in their own blocks so you can easily change them in one place
and not have to edit every single server definition. Passwords are linked
to server via the "identifier" string.

Config files should be kept in secure locations (your home directory) with 
appropriate file permissions (600). You can not supply the server address and password
via the command line for security reasons.

## Usage
```
q2rcon [--config <alternate_configfile>] <server_name> <console_cmd>`
```

The `<server_name>` argument matches the "name" field in the server definitions.
The `<console_cmd>` argument is the command you would normally enter into
the server console.

It's not necessary to wrap the whole command in quotes, unless the command needs to be
quoted as it's executed on the server. They'll need to be escaped, see the last example.

## Examples
```
# get the server and player info for server1
q2rcon server1 status

# set the flags for server1
q2rcon server1 set dmflags 1023

# instruct q2admin to mute client #2 for 5 minutes
q2rcon server1 sv !mute cl 2 300

# load custom config and tell myserver to shutdown 
q2rcon --config=./pfservers.txt myserver recycle

# set the hostname that includes quotes
q2rcon server 1 "set hostname \"Joe's Server\""
```

## Dependencies
```
go get google.golang.org/protobuf
go get github.com/packetflinger/libq2
```

## Building
```
go build .
```
