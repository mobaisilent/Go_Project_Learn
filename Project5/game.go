package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Mode int

// 创建全局const常量
const (
	ModeTitle Mode = iota
	ModeGame
	ModeOver
)

type Game struct {
	input     *Input
	ship      *Ship
	cfg       *Config
	bullets   map[*Bullet]struct{}
	aliens    map[*Alien]struct{}
	mode      Mode
	failCount int // 被外星人碰撞和移出屏幕的外星人数量之和
	overMsg   string
}

func (g *Game) init() {
	g.CreateAliens()
	g.CreateFonts()
}

func NewGame() *Game {
	cfg := loadConfig()
	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle(cfg.Title)

	g := &Game{
		input: &Input{
			msg: "hello world",
		},
		ship:    NewShip(cfg.ScreenWidth, cfg.ScreenHeight),
		cfg:     cfg,
		bullets: make(map[*Bullet]struct{}),
		aliens:  make(map[*Alien]struct{}),
	}
	// 调用 CreateAliens 创建一组外星人
	g.CreateAliens()
	g.init()
	return g
}

func (g *Game) addBullet(bullet *Bullet) {
	g.bullets[bullet] = struct{}{}
}

func (g *Game) addAlien(alien *Alien) {
	g.aliens[alien] = struct{}{}
}

func (g *Game) CreateAliens() {
	alien := NewAlien(g.cfg)

	availableSpaceX := g.cfg.ScreenWidth - 2*alien.width
	numAliens := availableSpaceX/(2*alien.width) + 1

	for i := 0; i < numAliens; i++ {
		alien = NewAlien(g.cfg)
		alien.x = float64(alien.width + 2*alien.width*i)
		g.addAlien(alien)
	}

	for row := 0; row < 2; row++ {
		for i := 0; i < numAliens; i++ {
			alien = NewAlien(g.cfg)
			alien.x = float64(alien.width + 2*alien.width*i)
			alien.y = float64(alien.height*row) * 1.5
			g.addAlien(alien)
		}
	}
}
