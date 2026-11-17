# TCP Chat Server

A learning project in Go. Multiple terminals connect to one server and broadcast messages to each other, like a barebones group chat.

## Status

- [x] Server accepts multiple connections
- [x] Each connection handled in its own goroutine
- [x] Client sends lines from stdin
- [x] Server logs received messages
- [ ] Server broadcasts messages to other clients
- [ ] Client displays received messages
- [ ] Usernames and join/leave notices

## Running

Both files are `package main`, so run them individually.

Server:

```sh
go run server.go
```

Client (one per terminal):

```sh
go run client.go
```

Or connect with netcat:

```sh
nc localhost 8080
```

## Structure

```
server.go   # listener + per-connection handler
client.go   # interactive terminal client
```

## How it works

- `server.go` listens on `:8080`, accepts connections in a loop, and spawns a goroutine per client. Each handler reads newline-terminated messages with `bufio` and logs them.
- `client.go` dials `localhost:8080` and writes each stdin line to the socket.

The client currently only sends, and the server only logs — no relay yet.

## Next steps

1. Keep a mutex-guarded registry of connections and write incoming messages to everyone except the sender.
2. Read from the socket in a goroutine on the client while reading stdin on the main one.
3. Add usernames, join/leave notices, and cleanup of dead connections.
