package main

import "github.com/kercylan98/minotaur/server/client"

func main() {
	cli := client.NewWebsocket("ws://127.0.0.1:9999")
	cli.RegConnectionOpenedEvent(func(conn *client.Client) {
		conn.WriteWS(2, []byte("6"))
	})
	if err := cli.Run(true); err != nil {
		panic(err)
	}
}
