package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	input   *Input
	ship    *Ship
	cfg     *Config
	bullets map[*Bullet]struct{}
}

func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(cfg.Title)

	return &Game{
		input: &Input{
			msg: "hello world",
		},
		ship:    NewShip(cfg.ScreenWidth, cfg.ScreenHeight),
		cfg:     cfg,
		bullets: make(map[*Bullet]struct{}),
	}
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}
