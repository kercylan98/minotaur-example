package typedefines

import (
	"github.com/kercylan98/minotaur-example/protocol/protocol"
	"github.com/kercylan98/minotaur/server"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewPlayer 创建玩家
func NewPlayer(conn *server.Conn) *Player {
	player := &Player{
		conn: conn,
	}
	return player
}

// Player 玩家
type Player struct {
	id   int64
	conn *server.Conn
}

// GetID 获取玩家ID
func (slf *Player) GetID() int64 {
	return slf.id
}

// AffirmLogin 设置玩家为登录状态
func (slf *Player) AffirmLogin(id int64) {
	slf.id = id
}

// HasLogged 是否已登录
func (slf *Player) HasLogged() bool {
	return slf.id != 0
}

// Push 向玩家推送消息
func (slf *Player) Push(msgType, msgId protoreflect.Enum, message proto.Message, callback ...func(err error)) {
	var msg = &protocol.Message{
		Type: int32(msgType.Number()),
		Id:   int32(msgId.Number()),
	}

	var (
		bytes []byte
		err   error
	)

	if message != nil {
		bytes, err = proto.Marshal(message)
		if err != nil {
			panic(err)
		}
		msg.Body = bytes
	}

	bytes, err = proto.Marshal(msg)
	if err != nil {
		panic(err)
	}

	slf.conn.Write(bytes, callback...)
}
