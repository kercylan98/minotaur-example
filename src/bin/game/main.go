package main

import (
	"github.com/kercylan98/minotaur-example/src/bin/game/application"
	_ "github.com/kercylan98/minotaur-example/src/bin/game/services"
)

func main() {
	application.Start(":9999")
}
