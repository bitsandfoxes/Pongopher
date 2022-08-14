package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"time"
)

const (
	screenWidth  = 256
	screenHeight = 256

	gameBoundsLeft   = 16
	gameBoundsRight  = screenWidth - 16
	gameBoundsTop    = 16
	gameBoundsBottom = screenHeight - 16

	menu     = "menu"
	playing  = "playing"
	gameOver = "gameOver"
)

type Palette struct {
	blue   color.RGBA
	red    color.RGBA
	yellow color.RGBA
	green  color.RGBA
	white  color.RGBA
}

type Game struct {
	pressedKeys  []ebiten.Key
	palette      Palette
	currentState string

	titleFace font.Face
	fontFace  font.Face
}

var (
	prevUpdateTime = time.Now()
	deltaTime      float64

	//go:embed embed/images/Gopher.png
	gopherRaw           []byte
	gopher              *ebiten.Image
	gopherRotation      float64
	gopherRotationSpeed float64

	ball          GameObject
	ballDirection Vector

	playerPaddle GameObject
	aiPaddle     GameObject

	lineDot *ebiten.Image
)

func (g *Game) Update() error {
	UpdateTime()

	switch g.currentState {
	case menu:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.currentState = playing
		}
	case playing:
		g.HandleInput()
		g.HandleAi()
		g.HandleBall()
	case gameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			ball.SetPosition(Vector{screenWidth/2 - 8, screenHeight/2 - 8})
			ballDirection = Vector{-1, 0}

			playerPaddle.SetPosition(Vector{gameBoundsLeft + 8 + 8, screenHeight / 2})
			aiPaddle.SetPosition(Vector{gameBoundsRight - 8 - 8, screenHeight / 2})
			g.currentState = playing
		}
	}
	return nil
}

func UpdateTime() {
	deltaTime = time.Since(prevUpdateTime).Seconds()
	prevUpdateTime = time.Now()
}

func (g *Game) HandleInput() {
	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	for _, key := range g.pressedKeys {
		switch key {
		case ebiten.KeyW:
			playerPaddle.SetPosition(Vector{playerPaddle.position.x, playerPaddle.position.y - playerPaddle.speed*deltaTime})
		case ebiten.KeyS:
			playerPaddle.SetPosition(Vector{playerPaddle.position.x, playerPaddle.position.y + playerPaddle.speed*deltaTime})
			//case ebiten.KeyA:
			//	playerPaddle.position.x -= playerPaddle.speed.x * deltaTime
			//case ebiten.KeyD:
			//	playerPaddle.position.x += playerPaddle.speed.x * deltaTime
		}
	}

	if playerPaddle.bounds.top < gameBoundsTop {
		playerPaddle.SetPosition(Vector{playerPaddle.position.x, gameBoundsTop + playerPaddle.size.y/2})
	} else if playerPaddle.bounds.bottom > gameBoundsBottom {
		playerPaddle.SetPosition(Vector{playerPaddle.position.x, gameBoundsBottom - playerPaddle.size.y/2})
	}
}

func (g *Game) HandleBall() {
	ballVector := ballDirection.Multiply(ball.speed * deltaTime)
	ball.SetPosition(ball.position.Add(ballVector))

	if ball.bounds.right > gameBoundsRight || ball.bounds.left < gameBoundsLeft {
		g.currentState = gameOver
		return
	}

	if ball.bounds.top <= gameBoundsTop {
		ball.SetPosition(Vector{ball.position.x, gameBoundsTop + ball.size.y/2})
		ballDirection.y = math.Abs(ballDirection.y)
	} else if ball.bounds.bottom >= gameBoundsBottom {
		ball.SetPosition(Vector{ball.position.x, gameBoundsBottom - ball.size.y/2})
		ballDirection.y = math.Abs(ballDirection.y) * -1
	}

	if ball.Hit(playerPaddle) {
		ballDirection.x = math.Abs(ballDirection.x)

		reflection := ball.position.Subtract(playerPaddle.position)
		ballDirection = reflection.Normalized()
	} else if ball.Hit(aiPaddle) {
		ballDirection.x = math.Abs(ballDirection.x) * -1

		reflection := ball.position.Subtract(aiPaddle.position)
		ballDirection = reflection.Normalized()
	}
}

func (g *Game) HandleAi() {
	movement := aiPaddle.speed * deltaTime
	distanceToBall := math.Abs(ball.position.y - aiPaddle.position.y)

	movement = math.Min(movement, distanceToBall)

	if aiPaddle.position.y > ball.position.y {
		movement *= -1
	}

	aiPaddle.SetPosition(Vector{aiPaddle.position.x, aiPaddle.position.y + movement})
	if aiPaddle.bounds.top < gameBoundsTop {
		aiPaddle.SetPosition(Vector{aiPaddle.position.x, gameBoundsTop + aiPaddle.size.y/2})
	} else if aiPaddle.bounds.bottom > gameBoundsBottom {
		aiPaddle.SetPosition(Vector{aiPaddle.position.x, gameBoundsBottom - aiPaddle.size.y/2})
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.palette.blue)

	switch g.currentState {
	case menu:
		g.DrawStart(screen)
	case playing:
		g.DrawGame(screen)
	case gameOver:
		g.DrawGameOver(screen)
	}
}

func (g *Game) DrawStart(screen *ebiten.Image) {
	gopherDrawOptions := &ebiten.DrawImageOptions{}
	gopherDrawOptions.GeoM.Scale(4, 4)
	gopherDrawOptions.GeoM.Translate(screenWidth/2-16*4/2, 20)
	screen.DrawImage(gopher, gopherDrawOptions)

	DrawText(screen, "Welcome to", g.fontFace, 100)
	DrawText(screen, "Pon-Gopher", g.titleFace, 125)
	DrawText(screen, "'W' and 'S'", g.fontFace, 150)
	DrawText(screen, "to move up or down", g.fontFace, 160)
	DrawText(screen, "-- Space --", g.fontFace, 190)
	DrawText(screen, "to play", g.fontFace, 200)
}

func (g *Game) DrawGameOver(screen *ebiten.Image) {
	DrawText(screen, "lol", g.titleFace, 65)
	DrawText(screen, "u done gophering around", g.fontFace, 80)

	gopherDrawOptions := &ebiten.DrawImageOptions{}
	gopherRotation += gopherRotationSpeed * deltaTime
	gopherDrawOptions.GeoM.Translate(-8, -8)
	gopherDrawOptions.GeoM.Rotate(gopherRotation)
	gopherDrawOptions.GeoM.Translate(8, 8)

	gopherDrawOptions.GeoM.Scale(4, 4)

	gopherDrawOptions.GeoM.Translate(screenWidth/2-16*4/2, 85)
	screen.DrawImage(gopher, gopherDrawOptions)

	DrawText(screen, "-- Space --", g.fontFace, 175)
	DrawText(screen, "to gopher some more", g.fontFace, 185)
}

func DrawText(screen *ebiten.Image, stringToDraw string, face font.Face, y int) {
	bounds := text.BoundString(face, stringToDraw)
	text.Draw(screen, stringToDraw, face, screenWidth/2-bounds.Dx()/2, y, color.White)
}

func (g *Game) DrawGame(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, gameBoundsLeft, gameBoundsTop, gameBoundsRight, gameBoundsTop, color.White)
	ebitenutil.DrawLine(screen, gameBoundsLeft, gameBoundsTop, gameBoundsLeft, gameBoundsBottom, color.White)
	ebitenutil.DrawLine(screen, gameBoundsLeft, gameBoundsBottom, gameBoundsRight, gameBoundsBottom, color.White)
	ebitenutil.DrawLine(screen, gameBoundsRight, gameBoundsTop, gameBoundsRight, gameBoundsBottom, color.White)

	playerPaddleDrawOptions := &ebiten.DrawImageOptions{}
	playerPaddleDrawOptions.GeoM.Translate(playerPaddle.GetScreenPosition())
	screen.DrawImage(playerPaddle.image, playerPaddleDrawOptions)

	aiPaddleDrawOptions := &ebiten.DrawImageOptions{}
	aiPaddleDrawOptions.GeoM.Translate(aiPaddle.GetScreenPosition())
	screen.DrawImage(aiPaddle.image, aiPaddleDrawOptions)

	ballDrawOptions := &ebiten.DrawImageOptions{}
	ballDrawOptions.GeoM.Scale(2, 2)
	ballDrawOptions.GeoM.Translate(ball.GetScreenPosition())
	screen.DrawImage(gopher, ballDrawOptions)

	distance := gameBoundsBottom - gameBoundsTop
	distance /= 4

	lineDrawOptions := &ebiten.DrawImageOptions{}
	lineDrawOptions.GeoM.Translate(screenWidth/2-2, gameBoundsTop+2)
	screen.DrawImage(lineDot, lineDrawOptions)

	for i := 1; i < distance; i++ {
		if i%2 == 0 {
			lineDrawOptions.GeoM.Translate(0, float64(8))
			screen.DrawImage(lineDot, lineDrawOptions)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

func (g *Game) init() {
	g.InitFonts()
	g.InitPalette()
	g.InitGopher()

	ball = MakeGameObject(28, 24, screenWidth/2, screenHeight/2, 350, g.palette.yellow)
	ballDirection = Vector{-1, 0}

	playerPaddle = MakeGameObject(16, 64, gameBoundsLeft+8+8, screenHeight/2, 250, g.palette.yellow)
	aiPaddle = MakeGameObject(16, 64, gameBoundsRight-8-8, screenHeight/2, 250, g.palette.green)

	lineDot = ebiten.NewImage(4, 4)
	lineDot.Fill(g.palette.white)

	g.currentState = menu
}

func (g *Game) InitFonts() {
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	g.titleFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     36,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     18,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
func (g *Game) InitPalette() {
	g.palette = Palette{
		blue:   color.RGBA{R: 38, G: 84, B: 124, A: 255},
		red:    color.RGBA{R: 254, G: 95, B: 85, A: 255},
		yellow: color.RGBA{R: 255, G: 209, B: 102, A: 255},
		green:  color.RGBA{R: 6, G: 214, B: 160, A: 255},
		white:  color.RGBA{R: 254, G: 249, B: 255, A: 255},
	}
}
func (g *Game) InitGopher() {
	img, _, err := image.Decode(bytes.NewReader(gopherRaw))
	if err != nil {
		log.Fatal(err)
	}
	gopher = ebiten.NewImageFromImage(img)
	gopherRotationSpeed = 0.2
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Pon-go")

	if err := ebiten.RunGame(NewGame()); err != nil {
		panic(err)
	}
}
