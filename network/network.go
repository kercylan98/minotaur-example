package network

import (
	"fmt"
	"github.com/kercylan98/minotaur-example/game/defines"
	"github.com/kercylan98/minotaur/server"
	"github.com/kercylan98/minotaur/utils/random"
	"github.com/kercylan98/minotaur/utils/sole"
)

var (
	// Server websocket 服务器
	Server = server.New(server.NetworkWebsocket,
		server.WithTicker(50, false),
	)

	// Router 路由器
	Router = newRouter[*defines.Player](func(srv *server.Server, conn *server.Conn) *defines.Player {
		return conn.GetData(defines.ConnPlayerKey).(*defines.Player)
	})
)

// Init 初始化
func Init(address string) {
	initRouter()
	initServer()

	if err := Server.Run(address); err != nil {
		panic(err)
	}
}

func initRouter() {
	Server.RegConnectionReceivePacketEvent(Router.HandleConnectionPacket)
}

func initServer() {
	sole.RegNameSpace(SoleNamespaceNotLoginPlayer)
	Server.RegConnectionOpenedEvent(onConnectionOpened)
}

func onConnectionOpened(srv *server.Server, conn *server.Conn) {
	userId := conn.GetData("user_id").(string)
	player := defines.NewPlayer(fmt.Sprintf("%s:%d", random.HostName(), sole.GetWith(SoleNamespaceNotLoginPlayer)), conn)
	defines.NotLoginPlayers.Set(userId, player)
}
