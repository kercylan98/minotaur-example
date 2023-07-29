package events

import "github.com/kercylan98/minotaur-example/game/defines"

type (
	LoginBeforeEventHandle func(player *defines.Player)
	LoginAfterEventHandle  func(player *defines.Player)
)

var (
	loginBeforeEventHandles []LoginBeforeEventHandle
	loginAfterEventHandles  []LoginAfterEventHandle
)

// RegLoginBeforeEvent 注册登录前事件
func RegLoginBeforeEvent(handle LoginBeforeEventHandle) {
	loginBeforeEventHandles = append(loginBeforeEventHandles, handle)
}

func OnLoginBeforeEvent(player *defines.Player) {
	for _, handle := range loginBeforeEventHandles {
		handle(player)
	}
}

// RegLoginAfterEvent 注册登录后事件
func RegLoginAfterEvent(handle LoginAfterEventHandle) {
	loginAfterEventHandles = append(loginAfterEventHandles, handle)
}

func OnLoginAfterEvent(player *defines.Player) {
	for _, handle := range loginAfterEventHandles {
		handle(player)
	}
}
