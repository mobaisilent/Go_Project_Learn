package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Input struct {
	msg            string
	lastBulletTime time.Time
}

// NewGame()函数用于初始化游戏：为Game结构体的信息赋值

func (i *Input) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.ship.x -= g.cfg.ShipSpeedFactor
		if g.ship.x < -float64(g.ship.width)/2 {
			g.ship.x = -float64(g.ship.width) / 2
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.ship.x += g.cfg.ShipSpeedFactor
		if g.ship.x > float64(g.cfg.ScreenWidth)-float64(g.ship.width)/2 {
			g.ship.x = float64(g.cfg.ScreenWidth) - float64(g.ship.width)/2
		}
	}
	// 修改这里为else 也就是飞船发射子弹的逻辑和左右移动是无关的
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(g.bullets) < g.cfg.MaxBulletNum &&
			time.Since(i.lastBulletTime).Milliseconds() > g.cfg.BulletInterval {
			bullet := NewBullet(g.cfg, g.ship)
			g.addBullet(bullet)
			i.lastBulletTime = time.Now()
		}
	}
}

// 增加两个if语句实现越界检测即可
