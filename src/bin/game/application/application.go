package application

import (
	"fmt"
	"github.com/kercylan98/minotaur-example/protocol/protocol"
	"github.com/kercylan98/minotaur-example/src/pkg/packet"
	"github.com/kercylan98/minotaur-example/src/pkg/typedefines"
	"github.com/kercylan98/minotaur/server"
	"github.com/kercylan98/minotaur/server/router"
	"github.com/kercylan98/minotaur/utils/concurrent"
	"github.com/kercylan98/minotaur/utils/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"runtime/debug"
)

// 应用程序全局实例
var (
	MessageRouter = router.NewMultistage[func(player *typedefines.Player, reader packet.Reader)](router.WithRouteTrim[func(*typedefines.Player, packet.Reader)](onRouteTrim))
	SocketServer  = server.New(server.NetworkWebsocket)
)

// 应用程序私有实例
var (
	messagePool = concurrent.NewPool[*protocol.Message](1024*10, func() *protocol.Message {
		return new(protocol.Message)
	}, func(data *protocol.Message) {
		data.Type = 0
		data.Id = 0
		data.Body = nil
	})
)

func init() {
	SocketServer.RegConnectionOpenedEvent(onConnectionOpened, -1)
	SocketServer.RegConnectionReceivePacketEvent(onConnectionReceivePacket, -1)
}

func Start(address string) {
	if err := SocketServer.Run(address); err != nil {
		panic(err)
	}
}

func onRouteTrim(route any) any {
	enum, ok := route.(protoreflect.Enum)
	if !ok {
		return route
	}
	return int32(enum.Number())
}

func onConnectionOpened(srv *server.Server, conn *server.Conn) {
	conn.SetData(conn.GetIP(), typedefines.NewPlayer(conn))
}

func onConnectionReceivePacket(srv *server.Server, conn *server.Conn, packet []byte) {
	var player, exist = conn.GetData(conn.GetIP()).(*typedefines.Player)
	if !exist {
		panic(fmt.Errorf("abnormal connection, ip: %s", conn.GetIP()))
	}

	var msg = messagePool.Get()
	if err := proto.Unmarshal(packet, msg); err != nil {
		panic(err)
	}

	defer func() {
		messagePool.Release(msg)
		if err := recover(); err != nil {
			log.Error("Packet", zap.Any("error", err))
			debug.PrintStack()
		}
	}()

	if !player.HasLogged() && !((msg.Type == int32(protocol.MessageType_MT_System) && msg.Id == int32(protocol.MessageSystemID_MI_Heartbeat)) ||
		(msg.Type == int32(protocol.MessageType_MT_User) && msg.Id == int32(protocol.MessageUserID_MI_Handshake))) {
		conn.Close()
		return
	}

	handle := MessageRouter.Match(msg.Type, msg.Id)
	if handle == nil {
		panic(fmt.Errorf("not exist route: [%d] %d", msg.Type, msg.Id))
	}

	handle(player, func(message proto.Message) proto.Message {
		if err := proto.Unmarshal(msg.Body, message); err != nil {
			panic(err)
		}
		return message
	})
}
