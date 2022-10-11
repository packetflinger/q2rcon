# Q2rcon
A command line out-of-band management client for Quake 2.

## Config File
Servers, passwords and mappings are done via the config file `.q2servers.json`
in your home directory. This file is used by other programs as well 
so it might contain more inforamation than is strictly needed for this
utility.

```
{
  "passwords": [
    {
      "name": "pass1",
      "password": "thisisareallybadpassword"
    },
    {
      "name": "pass2",
      "password": "ca00b050-495b-11ed-93de-cb7ecce5017d"
    }
  ],
  "servers": [
    {
      "name": "server1",
      "groups": "deathmatch usa",
      "addr": "100.64.3.5:27910",
      "password": "pass2"
    },
    {
      "name": "server2",
      "groups": "rocketarena usa",
      "addr": "192.0.2.44:27919",
      "password": "pass1"
    },
  ]
}
```
The example above shows 2 servers and 2 different passwords. Passwords
are defined in their own blocks so you can easily change them in one place
and not have to edit every single server definition.

Config files should be kept in secure locations (your home directory) with 
appropriate file permissions. You can not supply the server address and password
via the command line for security reasons.

## Usage
`q2rcon [-v] [-c <alternate_configfile>] <server_name> <console_cmd>`

The `-v` flag is for verbose output, giving you slightly more info. The `-c`
flag is to specify a config file other than `~/.q2server.json`.

The `<server_name>` argument matches the "name" field in the server definitions.
The `<console_cmd>` argument is the command you would normally enter into
the server console.

## Examples
```
# get the server and player info for server1
q2rcon -v server1 status

# set the flags for server1
q2rcon server1 set dmflags 1023

# instruct q2admin to mute client #2 for 5 minutes
q2rcon server2 sv !mute cl 2 300

# load custom config and tell myserver to shutdown 
q2rcon -c ./pfservers.txt myserver recycle
```
