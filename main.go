package main

import (
	_ "github.com/kercylan98/minotaur-example/game/services/system"
	"github.com/kercylan98/minotaur-example/network"
)

func main() {
	network.Init(":9999")
}
