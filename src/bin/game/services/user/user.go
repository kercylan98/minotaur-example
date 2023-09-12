package user

import (
	"github.com/kercylan98/minotaur-example/protocol/protocol"
	"github.com/kercylan98/minotaur-example/src/bin/game/application"
	"github.com/kercylan98/minotaur-example/src/pkg/packet"
	"github.com/kercylan98/minotaur-example/src/pkg/typedefines"
	"github.com/kercylan98/minotaur/utils/sole"
)

func init() {
	application.MessageRouter.Register(protocol.MessageType_MT_User, protocol.MessageUserID_MI_Handshake).Bind(onHandshake)
}

func onHandshake(player *typedefines.Player, reader packet.Reader) {
	var message = new(protocol.CMessageUserHandshake)
	reader.ReadTo(message)

	if message.Account != "minotaur" {
		return
	}

	// 查询ID或生成ID
	player.AffirmLogin(sole.SonyflakeID())

	player.Push(protocol.MessageType_MT_User, protocol.MessageUserID_MI_Handshake, &protocol.SMessageUserHandshake{Id: player.GetID()})
}
