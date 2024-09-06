package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// 创建子弹结构体：其实就是个 黑色的矩形
type Bullet struct {
	image       *ebiten.Image
	width       int
	height      int
	x           float64
	y           float64
	speedFactor float64
}

// 创建子弹信息初始化的方法
func NewBullet(cfg *Config, ship *Ship) *Bullet {
	rect := image.Rect(0, 0, cfg.BulletWidth, cfg.BulletHeight)
	img := ebiten.NewImage(rect.Dx(), rect.Dy())
	img.Fill(cfg.BulletColor)

	return &Bullet{
		image:       img,
		width:       cfg.BulletWidth,
		height:      cfg.BulletHeight,
		x:           ship.x + float64(ship.width-cfg.BulletWidth)/2,
		y:           ship.y - float64(cfg.BulletHeight),
		speedFactor: cfg.BulletSpeedFactor,
	}
}

// 创建子弹绘制的方法
func (bullet *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(bullet.x, bullet.y)
	screen.DrawImage(bullet.image, op)
}

// 判断子弹是否全部离开屏幕
func (bullet *Bullet) outOfScreen() bool {
	return bullet.y < -float64(bullet.height)
}
