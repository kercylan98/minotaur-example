package network

import (
	"github.com/kercylan98/minotaur-example/game/defines"
	"github.com/kercylan98/minotaur-example/protobuf/protobuf"
	"github.com/kercylan98/minotaur/game"
	"github.com/kercylan98/minotaur/server"
	mr "github.com/kercylan98/minotaur/server/router"
	"google.golang.org/protobuf/proto"
)

type PackerReader func(caseTo proto.Message) any
type Handle[T any] func(T, PackerReader)

// ReadTo 读取数据到指定的 proto.Message
func (slf PackerReader) ReadTo(param proto.Message) {
	slf(param)
}

// newRouter 创建路由器
func newRouter[T game.Player[string]](converter func(srv *server.Server, conn *server.Conn) T) *router[T] {
	return &router[T]{
		Level1Router: mr.NewLevel1Router[protobuf.MessageID, Handle[T]](),
		converter:    converter,
	}
}

// Router 路由器
type router[T game.Player[string]] struct {
	*mr.Level1Router[protobuf.MessageID, Handle[T]]
	converter func(srv *server.Server, conn *server.Conn) T
}

// HandleConnectionPacket 处理连接数据包
func (slf *router[T]) HandleConnectionPacket(srv *server.Server, conn *server.Conn, packet server.Packet) {
	var req = new(protobuf.Message)
	if err := proto.Unmarshal(packet.Data, req); err != nil {
		conn.Close()
		return
	}

	player := slf.converter(srv, conn)

	switch req.Id {
	case protobuf.MessageID_System_Heartbeat, protobuf.MessageID_System_Login:
	default:
		if !defines.IsLogin(player) {
			conn.Close()
			return
		}
	}

	slf.Match(req.Id)(player, func(caseTo proto.Message) any {
		if err := proto.Unmarshal(req.Body, caseTo); err != nil {
			panic(err)
		}
		return caseTo
	})
}
