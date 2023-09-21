package data

import (
	"errors"
	"github.com/kercylan98/minotaur-example/src/bin/game/application"
	"github.com/kercylan98/minotaur-example/src/pkg/typedefines"
	"github.com/kercylan98/minotaur/utils/memory"
	"github.com/kercylan98/minotaur/utils/super"
	"gorm.io/gorm"
	"time"
)

var (
	GetPlayer = memory.BindAction("player", func(id int64) *typedefines.Player {
		var player = new(typedefines.Player)
		if err := application.MySQL.First(player, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			} else {
				panic(err)
			}
		}
		return player
	})

	GetPlayerByAccount = memory.BindAction("playerByAccount", func(account string) *typedefines.Player {
		var player = new(typedefines.Player)
		if err := application.MySQL.Where("account = ?", account).First(player).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			} else {
				panic(err)
			}
		}
		return player
	})
)

var (
	GetPlayerPersist = memory.BindPersistCacheProgram("player", func(player *typedefines.Player) {
		super.RetryForever(time.Millisecond*100, func() error {
			return application.MySQL.Save(player).Error
		})
	})
)
