package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Input struct {
	msg string
}

// NewGame()函数用于初始化游戏：为Game结构体的信息赋值

func (i *Input) Update(ship *Ship, cfg *Config) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		ship.x -= cfg.ShipSpeedFactor
		if ship.x < -float64(ship.width)/2 {
			ship.x = -float64(ship.width) / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		ship.x += cfg.ShipSpeedFactor
		if ship.x > float64(cfg.ScreenWidth)-float64(ship.width)/2 {
			ship.x = float64(cfg.ScreenWidth) - float64(ship.width)/2
		}
	}
}

// 增加两个if语句实现越界检测即可
