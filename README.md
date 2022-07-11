# DRPC
[![GoDoc](https://godoc.org/github.com/altfoxie/drpc?status.svg)](https://godoc.org/github.com/altfoxie/drpc)

Discord RPC implementation in Go.

```bash
$ go get -u github.com/altfoxie/drpc
```

## Usage
You need to create a client.
```go
client, err := drpc.New("APP_ID")
if err != nil { /* handle error */ }
```

Now you can use the client to set player's activity. DRPC maintains a persistent connection to Discord RPC, so you don't need to worry about reconnecting.
```go
err := client.SetActivity(drpc.Activity{
    State: "Playing",
    Details: "Unranked PvP",
    Timestamps: &drpc.Timestamps{
        Start: time.Now(),
    },
    // Other fields...
})
if err != nil { /* handle error */ }
```

## Roadmap
 - [x] Support Windows
 - [ ] Support events
 - [ ] Match responses with method calls by `nonce`

## License
MIT