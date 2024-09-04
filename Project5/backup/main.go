// 改代码有copilot生成，检索图片的生成失效的问题

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Ship struct {
	image  *ebiten.Image
	width  int
	height int
}

func NewShip() *Ship {
	img, _, err := ebitenutil.NewImageFromFile("resource/ship.png")
	if err != nil {
		log.Fatal(err)
	}

	width, height := img.Size()
	ship := &Ship{
		image:  img,
		width:  width,
		height: height,
	}

	return ship
}

func (ship *Ship) Draw(screen *ebiten.Image, cfg *Config) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(cfg.ScreenWidth-ship.width)/2, float64(cfg.ScreenHeight-ship.height))
	screen.DrawImage(ship.image, op)
}

type Game struct {
	ship *Ship
	cfg  *Config
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ship.Draw(screen, g.cfg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cfg.ScreenWidth, g.cfg.ScreenHeight
}

type Config struct {
	ScreenWidth  int
	ScreenHeight int
}

func main() {
	cfg := &Config{
		ScreenWidth:  800,
		ScreenHeight: 600,
	}

	ship := NewShip()

	game := &Game{
		ship: ship,
		cfg:  cfg,
	}

	ebiten.SetWindowSize(cfg.ScreenWidth, cfg.ScreenHeight)
	ebiten.SetWindowTitle("Ship Example")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
