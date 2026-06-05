# Mezon Light SDK for Go

A lightweight Go SDK for [Mezon](https://mezon.ai) chat, ported from the
TypeScript package [`mezon-light-sdk`](../mezon-light-sdk).

## Features

- Authenticate with an ID token, or restore a session from stored tokens
- Refresh sessions (single-flight; concurrent callers share one refresh)
- Create DM / group DM channels
- Upload attachments
- Realtime messaging over WebSocket using the protobuf wire protocol
  (join/leave channels, send/receive messages, heartbeat with automatic
  dead-connection detection)

## Installation

```sh
go get github.com/quangledang23/mezon-light-sdk-go
```

## Quick start

```go
package main

import (
	"context"
	"log"

	mezonlight "github.com/quangledang23/mezon-light-sdk-go"
)

func main() {
	ctx := context.Background()

	// Authenticate with an ID token from an identity provider…
	client, err := mezonlight.Authenticate(ctx, mezonlight.AuthenticateConfig{
		IDToken:  "id-token-from-provider",
		UserID:   "user-123",
		Username: "johndoe",
	})
	if err != nil {
		log.Fatal(err)
	}

	// …or restore from previously stored tokens:
	// client, err := mezonlight.InitClient(mezonlight.ClientInitConfig{
	// 	Token:        "your-token",
	// 	RefreshToken: "your-refresh-token",
	// 	APIURL:       "https://api.mezon.ai",
	// 	WSURL:        "gw.mezon.ai",
	// 	UserID:       "user-123",
	// })

	// Create a DM channel.
	channel, err := client.CreateDM(ctx, "peer-user-id")
	if err != nil {
		log.Fatal(err)
	}

	// Connect the realtime socket.
	socket := mezonlight.NewLightSocket(client, client.Session())
	err = socket.Connect(ctx, mezonlight.SocketConnectOptions{
		OnError:      func(err error) { log.Println("socket error:", err) },
		OnDisconnect: func() { log.Println("disconnected") },
	})
	if err != nil {
		log.Fatal(err)
	}
	defer socket.Disconnect()

	// Receive messages. The returned function unsubscribes the handler.
	unsubscribe := socket.OnChannelMessage(func(msg *mezonlight.ChannelMessage) {
		log.Printf("received from %s: %v", msg.Username, msg.Content)
	})
	defer unsubscribe()

	// Join the DM channel and send a message.
	if err := socket.JoinDMChannel(ctx, channel.ChannelID); err != nil {
		log.Fatal(err)
	}
	err = socket.SendDM(ctx, mezonlight.SendMessagePayload{
		ChannelID: channel.ChannelID,
		Content:   map[string]string{"t": "Hello!"},
	})
	if err != nil {
		log.Fatal(err)
	}

	select {} // keep the process alive to receive messages
}
```

### Uploading attachments

```go
result, err := client.UploadAttachment(ctx, &mezonlight.ApiUploadAttachmentRequest{
	Filename: "image.png",
	Filetype: "image/png",
	Size:     1024,
	Width:    800,
	Height:   600,
})
// result.URL can be used in message attachments.
```

### Session management

```go
// Refresh before the token expires.
if client.IsSessionExpired() {
	if _, err := client.RefreshSession(ctx); err != nil {
		log.Fatal(err)
	}
}

// Persist the session for later restoration via InitClient.
config := client.ExportSession()
```

## Package layout

| Path        | Contents                                                            |
| ----------- | ------------------------------------------------------------------- |
| `.` (root)  | `LightClient`, `LightSocket`, `DefaultSocket`, `MezonApi`, `Session` |
| `proto`     | Hand-written protobuf wire codecs for the `mezon.api` and `mezon.realtime` messages used by the SDK (field numbers mirror the ts-proto generated code in `mezon-light-sdk/src/proto`) |

## Differences from the TypeScript SDK

- Methods take a `context.Context` and return `error` instead of promises.
- `WriteChatMessage` collects its many optional parameters in
  `ChatMessageOptions`.
- Callbacks (`OnChannelMessage`, `OnDisconnect`, …) are struct fields set
  before `Connect`.
- Snowflake IDs remain `string` in Go structs and are converted to/from
  int64 varints on the wire, matching ts-proto's int64-as-string behavior
  (an ID of `"0"` or `""` is omitted from the wire).
