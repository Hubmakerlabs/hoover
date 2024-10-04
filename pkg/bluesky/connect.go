package bluesky

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

func Connect(c context.Context, api string) (conn *websocket.Conn, err error) {
	// api := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	// fmt.Println("dialing: ", api)
	d := websocket.DefaultDialer
	conn, _, err = d.Dial(api, http.Header{})
	if err != nil {
		err = errorf.E("dial failure: %w", err)
	}
	go func() {
		<-c.Done()
		_ = conn.Close()
	}()
	return
}
