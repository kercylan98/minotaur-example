package defines

import (
	"github.com/kercylan98/minotaur/game"
	"github.com/kercylan98/minotaur/utils/concurrent"
)

var (
	// Players 玩家列表
	Players = concurrent.NewBalanceMap[string, *Player]()

	// NotLoginPlayers 未登录玩家列表
	NotLoginPlayers = concurrent.NewBalanceMap[string, *Player]()
)

// IsLogin 是否登录
func IsLogin(player game.Player[string]) bool {
	return Players.Exist(player.GetID())
}
