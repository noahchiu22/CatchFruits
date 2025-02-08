package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"strconv"

	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var fruits = []string{
	"apple",
	"orange",
	"peach",
	"cherry",
	"bomb",
}

var fruitsWeightMap = map[string]float64{
	"apple":  0,
	"orange": 0,
	"peach":  0,
	"cherry": 0,
	"bomb":   0,
}

type Sprite struct {
	Image                  *ebiten.Image
	name                   string
	X, Y                   float64 // position
	scaleX, scaleY         float64 // scale
	weight, speed, gravity float64 // physics
}

type Game struct {
	basket  *Sprite
	fruits  []*Sprite
	counter int
	points  int
	over    bool
	pause   bool
	level   int
}

func (g *Game) Update() error {
	width, height := ebiten.WindowSize()

	// game over or pause
	if g.over || g.pause {
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			if g.over {
				g.points = 0
				g.fruits = []*Sprite{}
			}
			g.over = false
			g.pause = false
		}
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.pause = true
	}

	// move the basket to the left
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.basket.X > 0 {
			g.basket.X -= (7 + float64(g.level))
		}
	}
	// move the basket to the right
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.basket.X+float64(g.basket.Image.Bounds().Dx()) < float64(width) {
			g.basket.X += (7 + float64(g.level))
		}
	}

	// generate a new fruit every 120 frames
	if g.counter%120 == 0 {
		cube := ebiten.NewImage(50, 50)
		fruit := &Sprite{
			Image:   cube,
			X:       rand.Float64() * float64(width-cube.Bounds().Dx()),
			Y:       -float64(cube.Bounds().Dy()),
			gravity: 3 + float64(g.level)*2,
		}
		err := fruit.genRandomFruit()
		if err != nil {
			log.Fatal(err)
			return err
		}

		g.fruits = append(g.fruits, fruit)
	}

	i := 0
	// add gravity to the fruits
	// check if the fruit is out of bounds
	for i < len(g.fruits) {
		fruit := g.fruits[i]

		intoBasket := fruit.Y > g.basket.Y &&
			fruit.X > g.basket.X &&
			fruit.X+float64(fruit.Image.Bounds().Dx()) < g.basket.X+float64(g.basket.Image.Bounds().Dx())

		fruit.gravity += fruit.weight
		fruit.Y += fruit.gravity
		if fruit.Y > float64(height) || intoBasket {
			g.fruits = append(g.fruits[:i], g.fruits[i+1:]...)
			if intoBasket {
				if fruit.name == "bomb" {
					g.over = true
					return nil
				}
				g.points++
			}
		}
		i++
	}

	g.level = g.points / 10

	// game counter
	g.counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 145, G: 209, B: 255, A: 255})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("screen size: %d*%d", screen.Bounds().Size().X, screen.Bounds().Size().Y), 0, 0)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("points: %d", g.points), 0, 14)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(g.basket.scaleX, g.basket.scaleY)
	opts.GeoM.Translate(g.basket.X, g.basket.Y)
	screen.DrawImage(g.basket.Image, opts)
	for _, fruit := range g.fruits {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(fruit.X, fruit.Y)
		screen.DrawImage(fruit.Image, opts)
	}

	// game over screen
	if g.over || g.pause {
		width, height := ebiten.WindowSize()
		gameOverScreen := ebiten.NewImage(width, height)
		vector.DrawFilledRect(gameOverScreen, 0, 0, float32(width), float32(height), color.RGBA{R: 0, G: 0, B: 0, A: 100}, true)
		if g.over {
			ebitenutil.DebugPrintAt(gameOverScreen, "Game Over", width/2-100, height/2-100)
		}
		ebitenutil.DebugPrintAt(gameOverScreen, "You got "+strconv.Itoa(g.points)+" points", width/2-100, height/2-50)
		ebitenutil.DebugPrintAt(gameOverScreen, "Press Enter to start", width/2-100, height/2)
		screen.DrawImage(gameOverScreen, nil)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowTitle("Hello, World!")
	_, height := ebiten.WindowSize()

	basketCube := ebiten.NewImage(100, 100)
	basketImg, _, err := ebitenutil.NewImageFromFile("assets/images/basket.png")
	if err != nil {
		log.Fatal(err)
	}

	basketCube.DrawImage(basketImg, shrinkIntoCube(basketCube, basketImg))
	basket := &Sprite{
		Image:  basketCube,
		scaleX: 1,
		scaleY: 1,
		X:      100,
		Y:      float64(height - basketCube.Bounds().Dy()),
	}

	if err := ebiten.RunGame(&Game{
		basket: basket,
	}); err != nil {
		log.Fatal(err)
	}
}

func shrinkIntoCube(cube *ebiten.Image, img *ebiten.Image) *ebiten.DrawImageOptions {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(float64(cube.Bounds().Dx())/float64(img.Bounds().Dx()),
		float64(cube.Bounds().Dy())/float64(img.Bounds().Dy()))
	return opts
}

func (f *Sprite) genRandomFruit() error {
	i := rand.IntN(len(fruits))

	img, _, err := ebitenutil.NewImageFromFile("assets/images/" + fruits[i] + ".png")
	if err != nil {
		log.Fatal(err)
		return err
	}
	f.Image.DrawImage(img, shrinkIntoCube(f.Image, img))
	f.weight = fruitsWeightMap[fruits[i]]
	f.name = fruits[i]
	return nil
}
