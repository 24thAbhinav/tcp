# TCP Chat Server

A minimal group chat over raw TCP in Go. Multiple terminals connect to one server and every message is relayed to the other connected clients.

## Features

- Many clients at once, one goroutine per connection
- Broadcast to everyone except the sender
- Join and leave notices
- Mutex-protected client registry
- Client can receive messages while typing

## Running

Both files are `package main`, so run them individually. The server listens on `localhost:8080`.

Server:

```sh
go run server.go
```

Client (one per terminal):

```sh
go run client.go
```

Enter a name, then chat. To connect without Go, use `nc localhost 8080` and send your name on the first line.

## Structure

```
server.go   # listener, per-connection handler, broadcast
client.go   # stdin -> socket, socket -> screen
```

## How it works

### Server

```
                    server.go

         Accept loop (main goroutine)
                    │
        ┌───────────┼───────────┐
        ▼           ▼           ▼
  goroutine A  goroutine B  goroutine C    <- one per connection
        │           │           │
        └───────────┼───────────┘
                    ▼
           connectionMap (sync.Mutex)
             username -> net.Conn
                    │
            broadcast(sender, text)
                    │
      write to every conn except sender
```

The main goroutine does one thing: accept connections. Every accepted socket is handed to `handleConnection` in its own goroutine, so one client blocking on a read never stops the others.

Each handler reads the username, adds `username -> conn` to the shared `connectionMap`, then loops reading lines and calling `broadcast`. Since every handler goroutine touches the same map, all reads and writes go through `connectionMapMut`. Inside `broadcast` the lock is held only long enough to copy the recipient list, then released before writing — a slow client does not block the rest.

When a read fails with `io.EOF`, the handler broadcasts a leave notice, deletes the user from the map under the lock, and returns.

### Client

```
                    client.go

    main goroutine              reader goroutine
    ──────────────              ────────────────
    stdin ──▶ socket            socket ──▶ screen

    both print through terminalMu, so lines never mix
```

The main goroutine blocks on `stdin`, so you can type at any time. A second goroutine blocks on reads from the server, so incoming messages appear while you type. Both goroutines write to the terminal, so `terminalMu` keeps their output from interleaving.

## Notes

- Address and port are hardcoded (`localhost:8080`).
- Usernames are not checked for uniqueness; a duplicate name overwrites the earlier entry.
- Messages are relayed only, never stored.
