package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sprite struct {
	Image      *ebiten.Image
	X          int
	Y          int
	forceRight float64
	forceDown  float64
	scaleX     float64
	scaleY     float64
}

type Game struct {
	ball *Sprite
}

func (g *Game) Update() error {
	width, height := ebiten.WindowSize()
	// on the ground
	isOnTheGround := g.ball.Y+g.ball.Image.Bounds().Dy() >= height
	// hit the wall
	isHitTheWall := g.ball.X+g.ball.Image.Bounds().Dx() >= width || g.ball.X <= 0

	fmt.Println("ball bottom", g.ball.Y+g.ball.Image.Bounds().Dy(), "forceDown", g.ball.forceDown)
	// jump
	if isOnTheGround && ebiten.IsKeyPressed(ebiten.KeySpace) {
		if g.ball.forceDown > 0 {
			g.ball.forceDown = 0
		}
		if g.ball.forceDown > -28 {
			g.ball.forceDown -= 2
		}
		g.ball.scaleY = (56 + g.ball.forceDown) / 56
		return nil
	}
	g.ball.scaleY = 1

	// gravity
	if !isOnTheGround {
		g.ball.forceDown += 1
	} else {
		if g.ball.forceDown > 0 {
			g.ball.forceDown = -0.6 * g.ball.forceDown
		}
	}

	if isHitTheWall {
		g.ball.forceRight = -g.ball.forceRight
	}

	// move
	g.ball.X += int(g.ball.forceRight)
	g.ball.Y += int(g.ball.forceDown)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 145, G: 209, B: 255, A: 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("screen size: %d*%d", screen.Bounds().Size().X, screen.Bounds().Size().Y), 0, 0)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(g.ball.scaleX, g.ball.scaleY)
	opts.GeoM.Translate(float64(g.ball.X), float64(g.ball.Y+int(float64(g.ball.Image.Bounds().Dy())*(1-opts.GeoM.Element(1, 1)))))
	screen.DrawImage(g.ball.Image, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")
	ball := ebiten.NewImage(100, 100)
	vector.DrawFilledCircle(ball, 50, 50, 50, color.White, true)
	width, height := ebiten.WindowSize()
	fmt.Println(width, height)
	if err := ebiten.RunGame(&Game{
		ball: &Sprite{
			Image:      ball,
			X:          10,
			Y:          height - ball.Bounds().Dy(),
			forceRight: 5,
			scaleX:     1,
			scaleY:     1,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
