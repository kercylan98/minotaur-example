package system

import (
	"github.com/kercylan98/minotaur-example/protocol/protocol"
	"github.com/kercylan98/minotaur-example/src/bin/game/application"
	"github.com/kercylan98/minotaur-example/src/pkg/packet"
	"github.com/kercylan98/minotaur-example/src/pkg/typedefines"
)

func init() {
	application.MessageRouter.Register(protocol.MessageType_MT_System, protocol.MessageSystemID_MI_Heartbeat).Bind(onHeartbeat)
}

func onHeartbeat(player *typedefines.Player, reader packet.Reader) {
	player.Push(protocol.MessageType_MT_System, protocol.MessageSystemID_MI_Heartbeat, nil)
}
