package main

import (
	"bytes"
	"encoding/json"
	"image/color"
	"log"
)

type Config struct {
	ScreenWidth       int        `json:"screenWidth"`
	ScreenHeight      int        `json:"screenHeight"`
	Title             string     `json:"title"`
	BgColor           color.RGBA `json:"bgColor"`
	ShipSpeedFactor   float64    `json:"shipSpeedFactor"`
	BulletSpeedFactor float64    `json:"bulletSpeedFactor"`
	BulletWidth       int        `json:"bulletWidth"`
	BulletHeight      int        `json:"bulletHeight"`
	BulletColor       color.RGBA `json:"bulletColor"`
	MaxBulletNum      int        `json:"maxBulletNum"`
	BulletInterval    int64      `json:"bulletInterval"`
	AlienSpeedFactor  float64    `json:"alienSpeedFactor"`
	TitleFontSize     int        `json:"titleFontSize"`
	FontSize          int        `json:"fontSize"`
	SmallFontSize     int        `json:"smallFontSize"`
}

func loadConfig() *Config {
	var cfg Config
	if err := json.NewDecoder(bytes.NewReader(ConfigJson)).Decode(&cfg); err != nil {
		log.Fatalf("json.Decode failed: %v\n", err)
	}

	return &cfg
}
