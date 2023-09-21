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
	conn *server.Conn

	ID       int64  `gorm:"column:id;primary_key;<-:create"`   // 玩家ID
	Account  string `gorm:"column:account"`                    // 玩家账号
	Password string `gorm:"column:password;type:varchar(128)"` // 玩家密码

	LastLogin   int64 `gorm:"column:last_login"`   // 最后登录时间
	LastOffline int64 `gorm:"column:last_offline"` // 最后离线时间

	CreatedAt int64 `gorm:"column:created_at;autoCreateTime"` // 创建时间
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime"` // 更新时间
}

func (slf *Player) TableName() string {
	return "player"
}

// GetID 获取玩家ID
func (slf *Player) GetID() int64 {
	return slf.ID
}

// AffirmLogin 设置玩家为登录状态
func (slf *Player) AffirmLogin(player *Player) {
	slf.conn.SetData(slf.conn.GetIP(), player)
}

// HasLogged 是否已登录
func (slf *Player) HasLogged() bool {
	return slf.ID != 0
}

// IsOnline 是否在线
func (slf *Player) IsOnline() bool {
	return slf.conn != nil && !slf.conn.IsClosed()
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
