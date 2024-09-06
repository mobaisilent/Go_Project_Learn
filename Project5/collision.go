package main

import (
)

// CheckCollision 检查子弹和外星人之间是否有碰撞
func CheckCollision(bullet *Bullet, alien *Alien) bool {
	alienTop, alienLeft := alien.y, alien.x
	alienBottom, alienRight := alien.y+float64(alien.height), alien.x+float64(alien.width)
	// 左上角
	x, y := bullet.x, bullet.y
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 右上角
	x, y = bullet.x+float64(bullet.width), bullet.y
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 左下角
	x, y = bullet.x, bullet.y+float64(bullet.height)
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	// 右下角
	x, y = bullet.x+float64(bullet.width), bullet.y+float64(bullet.height)
	if y > alienTop && y < alienBottom && x > alienLeft && x < alienRight {
		return true
	}

	return false
}
