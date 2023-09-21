package user

import (
	"github.com/kercylan98/minotaur-example/protocol/protocol"
	"github.com/kercylan98/minotaur-example/src/bin/game/application"
	"github.com/kercylan98/minotaur-example/src/pkg/data"
	"github.com/kercylan98/minotaur-example/src/pkg/packet"
	"github.com/kercylan98/minotaur-example/src/pkg/typedefines"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	application.MessageRouter.Register(protocol.MessageType_MT_User, protocol.MessageUserID_MI_Handshake).Bind(onHandshake)
}

func onHandshake(player *typedefines.Player, reader packet.Reader) {
	var message = new(protocol.CMessageUserHandshake)
	reader.ReadTo(message)

	loginPlayer := data.GetPlayerByAccount(message.Account)
	if loginPlayer == nil {
		return // 未注册的用户
	}

	if err := bcrypt.CompareHashAndPassword([]byte(loginPlayer.Password), []byte(message.Password)); err != nil {
		return // 密码错误
	}

	player.AffirmLogin(loginPlayer)

	player.Push(protocol.MessageType_MT_User, protocol.MessageUserID_MI_Handshake, &protocol.SMessageUserHandshake{Id: player.GetID()})
}
