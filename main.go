package main

import (
	"image/color"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type player struct {
	x, y  float32
	w, h  float32
	speed float32
	clr   color.RGBA
}

type ennemi struct {
	x, y         float32
	w, h         float32
	ShootCadence float32
	speed        float32
	clr          color.RGBA
}

type bullet struct {
	x, y             float32 // position actuelle
	TargetX, TargetY float32 // destination finale

	radius float32
	speed  float32
	clr    color.RGBA
}

type Game struct {
	player player
	ennemi ennemi
	bullet []bullet
}

func DirigePointToPoint(speed float32, point1 *ennemi, point2 player) {
	if point1.x-float32(point2.x) < 0 {
		point1.x += speed
	}
	if point1.x-float32(point2.x) > 0 {
		point1.x -= speed
	}

	if math.Abs(float64(point1.x)-float64(point2.x)) <= float64(speed) {
		point1.x = float32(point2.x)
	}
	// y
	if point1.y-float32(point2.y) < 0 {
		point1.y += speed
	}
	if point1.y-float32(point2.y) > 0 {
		point1.y -= speed
	}

	if math.Abs(float64(point1.y)-float64(point2.y)) <= float64(speed) {
		point1.y = float32(point2.y)
	}
}

func (b *bullet) DirigeToPlayer() {
	if b.x-float32(b.TargetX) < 0 {
		b.x += b.speed
	}
	if b.x-float32(b.TargetX) > 0 {
		b.x -= b.speed
	}

	if math.Abs(float64(b.x)-float64(b.TargetX)) <= float64(b.speed) {
		b.x = float32(b.TargetX)
	}
	// y
	if b.y-float32(b.TargetY) < 0 {
		b.y += b.speed
	}
	if b.y-float32(b.TargetY) > 0 {
		b.y -= b.speed
	}

	if math.Abs(float64(b.y)-float64(b.TargetY)) <= float64(b.speed) {
		b.y = float32(b.TargetY)
	}
}

func ShootBullet(ennemi *ennemi, bulletS []bullet, player player) []bullet {
	if math.Abs(float64(ennemi.x)-float64(player.x)) > 250 {
		return bulletS
	}

	if ennemi.ShootCadence > 0 {
		ennemi.ShootCadence--
		return bulletS
	}

	// Reinitialise la cadence
	ennemi.ShootCadence = 70

	bulletS = append(bulletS, bullet{
		x:       ennemi.x,
		y:       ennemi.y + ennemi.h/2,
		TargetX: player.x + player.w/2,
		TargetY: player.y + player.h/2,
		radius:  10,
		speed:   10,
		clr:     color.RGBA{255, 0, 255, 255},
	})

	return bulletS
}

// Deletes the bullet if arrived at the end
func DeleteBullets(bullet []bullet) []bullet {
	for i := range bullet {
		if bullet[i].TargetX == bullet[i].x && bullet[i].TargetY == bullet[i].y {
			return slices.Delete(bullet, 0, len(bullet))
		}
	}
	return bullet
}

func (p *player) Mouvement() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.x -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.x += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.y -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.y += p.speed
	}
}

func (g *Game) Update() error {
	// ennemi mouv.
	DirigePointToPoint(g.ennemi.speed, &g.ennemi, g.player)

	// player mouv.
	g.player.Mouvement()

	// create bullet
	g.bullet = ShootBullet(&g.ennemi, g.bullet, g.player)

	// dirige bullet
	for i := range g.bullet {
		g.bullet[i].DirigeToPlayer()
	}

	// delete bullets for the game not lagging
	g.bullet = DeleteBullets(g.bullet)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, g.player.x, g.player.y, g.player.w, g.player.h, g.player.clr, true)
	vector.DrawFilledRect(screen, g.ennemi.x, g.ennemi.y, g.ennemi.w, g.ennemi.h, g.ennemi.clr, true)
	for _, b := range g.bullet {
		vector.DrawFilledCircle(screen, b.x, b.y, b.radius, b.clr, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	g := &Game{
		player: player{
			x:     100,
			y:     100,
			w:     50,
			h:     50,
			speed: 5,
			clr:   color.RGBA{255, 255, 255, 255},
		},
		ennemi: ennemi{
			x:     screenWidth - 100,
			y:     screenHeight - 100,
			w:     50,
			h:     50,
			speed: 1,
			clr:   color.RGBA{0, 255, 255, 255},
		},
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
