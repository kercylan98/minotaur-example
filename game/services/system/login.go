package system

import (
	"github.com/kercylan98/minotaur-example/game/defines"
	"github.com/kercylan98/minotaur-example/game/events"
	"github.com/kercylan98/minotaur-example/network"
	"github.com/kercylan98/minotaur-example/protobuf/protobuf"
	"time"
)

func init() {
	network.Router.Route(protobuf.MessageID_System_Heartbeat, onHeartbeat)
	network.Router.Route(protobuf.MessageID_System_ServerTime, onServerTime)
	network.Router.Route(protobuf.MessageID_System_Login, onLogin)
}

func onHeartbeat(player *defines.Player, reader network.PackerReader) {
	player.Send(protobuf.MessageID_System_Heartbeat, nil)
}

func onServerTime(player *defines.Player, reader network.PackerReader) {
	player.Send(protobuf.MessageID_System_ServerTime, &protobuf.Time{Unix: time.Now().Unix()})
}

func onLogin(player *defines.Player, reader network.PackerReader) {
	var params = new(protobuf.SystemLoginS)
	reader.ReadTo(params)

	events.OnLoginBeforeEvent(player)
	// Login logic ...

	defines.NotLoginPlayers.Delete(player.GetID())
	token := player.SetLogged(params.UserId)
	defines.Players.Set(player.GetID(), player)

	events.OnLoginAfterEvent(player)

	player.Send(protobuf.MessageID_System_Login, &protobuf.SystemLoginC{
		UserId: player.GetID(),
		Token:  token,
	})
}
