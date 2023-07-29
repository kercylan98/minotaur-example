package defines

import (
	"github.com/kercylan98/minotaur-example/protobuf/protobuf"
	"github.com/kercylan98/minotaur/game/builtin"
	"github.com/kercylan98/minotaur/server"
	"github.com/kercylan98/minotaur/utils/str"
	"github.com/kercylan98/minotaur/utils/super"
	"google.golang.org/protobuf/proto"
)

// NewPlayer 创建玩家
//   - 当玩家存在时，将会复用玩家对象
func NewPlayer(id string, conn *server.Conn) *Player {
	player, exist := Players.GetExist(id)
	if !exist {
		player = &Player{
			Player: builtin.NewPlayer(id, conn),
		}
	} else {
		player.UseConn(conn)
	}
	player.GetConn().SetData(ConnPlayerKey, player)
	return player
}

type Player struct {
	*builtin.Player[string]
	loginId string // 登录后真实的 ID
}

// SetLogged 设置为已登录状态
func (slf *Player) SetLogged(id string) string {
	slf.loginId = id
	return slf.loginId
}

// GetID 获取玩家ID
func (slf *Player) GetID() string {
	return super.If(slf.loginId == str.None, slf.Player.GetID(), slf.loginId)
}

// Send 发送消息
func (slf *Player) Send(messageId protobuf.MessageID, msg proto.Message) {
	var message = &protobuf.Message{
		Id: messageId,
	}
	if msg != nil {
		body, err := proto.Marshal(msg)
		if err != nil {
			panic(err)
		}
		message.Body = body
	}
	resBinary, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}
	slf.GetConn().Write(server.Packet{Data: resBinary, WebsocketType: server.WebsocketMessageTypeBinary})
}
