package defines

import (
	"github.com/kercylan98/minotaur/game/fight"
	"github.com/kercylan98/minotaur/game/room"
	"github.com/kercylan98/minotaur/utils/concurrent"
	"sync/atomic"
)

// LandlordsData 斗地主数据
type LandlordsData struct {
	Guid        atomic.Int64
	RoomManager *room.Manager[string, *Player, *LandlordsRoom]
	InRoom      *concurrent.BalanceMap[string, int64] // 玩家所在房间
}

func (slf *LandlordsData) NextGuid() int64 {
	return slf.Guid.Add(1)
}

// LandlordsRoom 斗地主房间
type LandlordsRoom struct {
	Guid  int64
	Round *fight.Round[*LandlordsRoom]

	Players map[string]struct{} // 玩家
	Readies map[string]struct{} // 准备
}

func (slf *LandlordsRoom) GetGuid() int64 {
	return slf.Guid
}
