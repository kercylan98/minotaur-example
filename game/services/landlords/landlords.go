// Package landlords 包中仅提供了部分有针对性的逻辑以供参考，并非完整的游戏实现
package landlords

import (
	"github.com/kercylan98/minotaur-example/game/defines"
	"github.com/kercylan98/minotaur-example/network"
	"github.com/kercylan98/minotaur-example/protobuf/protobuf"
	"github.com/kercylan98/minotaur/game/fight"
	"github.com/kercylan98/minotaur/game/poker"
	"github.com/kercylan98/minotaur/game/room"
	"github.com/kercylan98/minotaur/utils/concurrent"
	"github.com/kercylan98/minotaur/utils/hash"
	"google.golang.org/protobuf/proto"
)

var (
	serviceData    *defines.LandlordsData
	getPlayerSeat  func(helper *room.Helper[string, *defines.Player, *defines.LandlordsRoom]) map[string]int32
	rule           *poker.Rule[int32, int32, *protobuf.Card]
	calcCardsScore func(cards []*protobuf.Card) int
)

func init() {
	serviceData = &defines.LandlordsData{
		RoomManager: room.NewManager[string, *defines.Player, *defines.LandlordsRoom](),
		InRoom:      concurrent.NewBalanceMap[string, int64](),
	}
	getPlayerSeat = func(helper *room.Helper[string, *defines.Player, *defines.LandlordsRoom]) map[string]int32 {
		players := make(map[string]int32)
		for seat, playerId := range helper.GetSeatInfoMap() {
			players[playerId] = int32(seat)
		}
		return players
	}
	rule = poker.NewRule[int32, int32, *protobuf.Card](
		poker.WithHand[int32, int32, *protobuf.Card]("三带一", 1, poker.HandThreeOfKindWithOne[int32, int32, *protobuf.Card]()),
	)
	calcCardsScore = func(cards []*protobuf.Card) int {
		var cardScore = func(card *protobuf.Card) int {
			switch card.Point {
			case 1:
				return 12
			case 2:
				return 13
			default:
				return int(card.Point) - 2
			}
		}
		name, _ := rule.PokerHand(cards...)
		var score int
		switch name {
		case "三带一":
			group := poker.GroupByPoint[int32, int32, *protobuf.Card](cards...)
			for _, cards := range group {
				if len(cards) == 3 {
					score += cardScore(cards[0]) * 100
				} else {
					score += cardScore(cards[0])
				}
			}
		default:
			for _, card := range cards {
				score += cardScore(card)
			}
		}
		return score
	}
}

func init() {
	network.Router.Route(protobuf.MessageID_Landlords_CreateRoom, onCreateRoom) // 创建房间
	network.Router.Route(protobuf.MessageID_Landlords_JoinRoom, onJoinRoom)     // 加入房间
	network.Router.Route(protobuf.MessageID_Landlords_LeaveRoom, onLeaveRoom)   // 离开房间
	network.Router.Route(protobuf.MessageID_Landlords_Ready, onReady)           // 准备
	network.Router.Route(protobuf.MessageID_Landlords_Start, onStart)           // 开始
	network.Router.Route(protobuf.MessageID_Landlords_Call, onCall)             // 叫地主
	network.Router.Route(protobuf.MessageID_Landlords_CallPass, onCallPass)     // 不叫
	network.Router.Route(protobuf.MessageID_Landlords_Grab, onGrab)             // 抢地主
	network.Router.Route(protobuf.MessageID_Landlords_GrabPass, onGrabPass)     // 不抢
	network.Router.Route(protobuf.MessageID_Landlords_Play, onPlay)             // 出牌
	network.Router.Route(protobuf.MessageID_Landlords_Pass, onPass)             // 不出
	network.Router.Route(protobuf.MessageID_Landlords_End, onEnd)               // 结束

	serviceData.RoomManager.RegPlayerJoinRoomEvent(onPlayerJoinRoom)   // 玩家加入房间
	serviceData.RoomManager.RegPlayerLeaveRoomEvent(onPlayerLeaveRoom) // 玩家离开房间
}

func onCreateRoom(player *defines.Player, reader network.PackerReader) {

	if serviceData.InRoom.Exist(player.GetID()) {
		return // 已经在房间中
	}

	r := &defines.LandlordsRoom{
		Guid: serviceData.NextGuid(),
	}
	serviceData.RoomManager.CreateRoom(r, room.WithPlayerLimit[string, *defines.Player, *defines.LandlordsRoom](3))

	helper := serviceData.RoomManager.GetHelper(r)
	_ = helper.Join(player)

	player.Send(protobuf.MessageID_Landlords_CreateRoom, &protobuf.LandlordsCreateRoomC{RoomId: r.Guid})
}

func onJoinRoom(player *defines.Player, reader network.PackerReader) {
	var params = new(protobuf.LandlordsJoinRoomS)
	reader.ReadTo(params)

	r := serviceData.RoomManager.GetRoom(params.RoomId)
	if serviceData.InRoom.Exist(player.GetID()) || r == nil {
		return // 已经在房间中 || 房间不存在
	}

	helper := serviceData.RoomManager.GetHelper(r)
	if err := helper.Join(player); err != nil {
		return // 人数已满
	}

	player.Send(protobuf.MessageID_Landlords_JoinRoom, &protobuf.LandlordsJoinRoomC{
		RoomId:       params.RoomId,
		UserIds:      hash.KeyToSlice(r.Players),
		ReadyUserIds: hash.KeyToSlice(r.Readies),
		UserSeat:     getPlayerSeat(helper),
	})

	message := &protobuf.LandlordsJoinRoomNotifyC{
		RoomId: params.RoomId,
		UserId: player.GetID(),
	}
	for _, p := range helper.GetPlayers() {
		if p.GetID() == player.GetID() {
			continue
		}
		p.Send(protobuf.MessageID_Landlords_Notify_JoinRoom, message)
	}
}

func onLeaveRoom(player *defines.Player, reader network.PackerReader) {

	if !serviceData.InRoom.Exist(player.GetID()) {
		return // 不在房间中
	}

	r := serviceData.RoomManager.GetRoom(serviceData.InRoom.Get(player.GetID()))
	if r.Round != nil {
		return // 已经开始游戏
	}

	helper := serviceData.RoomManager.GetHelper(r)
	helper.Leave(player)

	player.Send(protobuf.MessageID_Landlords_LeaveRoom, &protobuf.LandlordsLeaveRoomC{})

	message := &protobuf.LandlordsLeaveRoomNotifyC{
		RoomId: r.Guid,
		UserId: player.GetID(),
	}
	for _, p := range helper.GetPlayers() {
		if p.GetID() == player.GetID() {
			continue
		}
		p.Send(protobuf.MessageID_Landlords_Notify_LeaveRoom, message)
	}
}

func onReady(player *defines.Player, reader network.PackerReader) {
	// 常规逻辑不过多赘述，参考 onCreateRoom onJoinRoom onLeaveRoom
}

func onStart(player *defines.Player, reader network.PackerReader) {
	var params = new(protobuf.LandlordsStartGameS)
	reader.ReadTo(params)

	r := serviceData.RoomManager.GetRoom(params.RoomId)
	helper := serviceData.RoomManager.GetHelper(r)
	if r == nil || helper.GetPlayerCount() < helper.GetPlayerLimit() {
		return // 房间不存在 || 人数不足
	}

	// 阵容
	var camp = make([]*fight.RoundCamp, 3)
	for seat := range helper.GetSeatInfoMap() {
		camp[seat] = fight.NewRoundCamp(seat, seat)
	}

	r.Round = fight.NewRound[*defines.LandlordsRoom](r, camp, func(round *fight.Round[*defines.LandlordsRoom]) bool {
		// 结束逻辑校验
		return true
	})

	r.Round.Run()

}

func onCall(player *defines.Player, reader network.PackerReader) {
	var params = new(protobuf.LandlordsCallS)
	reader.ReadTo(params)

	r := serviceData.RoomManager.GetRoom(params.RoomId)
	helper := serviceData.RoomManager.GetHelper(r)
	if r == nil || r.Round == nil {
		return // 房间不存在 || 游戏未开始
	}

	seat := helper.GetSeat(player.GetID())
	if !r.Round.AllowAction(seat, seat) {
		return // 不是自己的操作回合
	}

	// ... 叫地主
}

func onCallPass(player *defines.Player, reader network.PackerReader) {
	// ... 不叫
}

func onGrab(player *defines.Player, reader network.PackerReader) {
	// ... 抢地主
}

func onGrabPass(player *defines.Player, reader network.PackerReader) {
	// ... 不抢
}

func onPlay(player *defines.Player, reader network.PackerReader) {
	// 假设玩家出牌
	reader = func(caseTo proto.Message) any {
		params := caseTo.(*protobuf.LandlordsPlayS)
		params.RoomId = 1
		params.Cards = []*protobuf.Card{{Point: 1, Color: 1}, {Point: 1, Color: 2}, {Point: 1, Color: 3}, {Point: 10, Color: 2}}
		reader.ReadTo(params)
		return params
	}
	var params = new(protobuf.LandlordsPlayS)
	reader.ReadTo(params)

	// 假设对手的牌
	targetCards := []*protobuf.Card{{Point: 10, Color: 1}, {Point: 10, Color: 2}, {Point: 10, Color: 3}, {Point: 2, Color: 2}}

	tName, _ := rule.PokerHand(targetCards...)
	name, _ := rule.PokerHand(params.Cards...)
	if tName != name {
		return // 牌型不一致
	}

	if calcCardsScore(params.Cards) < calcCardsScore(targetCards) {
		return // 牌型不够大
	}

	// ... 其他逻辑
}

func onPass(player *defines.Player, reader network.PackerReader) {

}

func onEnd(player *defines.Player, reader network.PackerReader) {

}

func onPlayerJoinRoom(landlordsRoom *defines.LandlordsRoom, player *defines.Player) {
	serviceData.InRoom.Set(player.GetID(), landlordsRoom.Guid)
}

func onPlayerLeaveRoom(landlordsRoom *defines.LandlordsRoom, player *defines.Player) {
	serviceData.InRoom.Delete(player.GetID())
}
