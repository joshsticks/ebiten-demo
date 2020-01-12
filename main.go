package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

var (
	mapWidth  int = 24
	mapHeight int = 24

	screenWidth  int = 640
	screenHeight int = 480

	worldMap = [][]int{
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}

	//x and y start position
	posX, posY float64 = 22, 12

	//initial direction vector
	dirX, dirY float64 = -1, 0

	//the 2d raycaster version of camera plane
	planeX, planeY float64 = 0, 0.66

	//time of current frame
	time float64 = 0

	//time of previous frame
	oldTime float64 = 0
)

func update(screen *ebiten.Image) error {

	handleMovement()

	for x := 0; x < screenWidth; x++ {
		//calculate ray position and direction
		var cameraX float64 = 2*float64(x)/float64(screenWidth) - 1 //x-coordinate in camera space
		var rayDirX float64 = dirX + planeX*cameraX
		var rayDirY float64 = dirY + planeY*cameraX

		//which box of the map we're in
		var mapX int = int(posX)
		var mapY int = int(posY)

		//length of ray from current position to next x or y-side
		var sideDistX float64
		var sideDistY float64

		//length of ray from one x or y-side to next x or y-side
		var deltaDistX float64 = math.Abs(1 / rayDirX)
		var deltaDistY float64 = math.Abs(1 / rayDirY)
		var perpWallDist float64

		//what direction to step in x or y-direction (either +1 or -1)
		var stepX int
		var stepY int

		var hit int = 0 //was there a wall hit?
		var side int    //was a NS or a EW wall hit?

		//calculate step and initial sideDist
		if rayDirX < 0 {
			stepX = -1
			sideDistX = (posX - float64(mapX)) * deltaDistX
		} else {
			stepX = 1
			sideDistX = (float64(mapX) + 1.0 - posX) * deltaDistX
		}

		if rayDirY < 0 {
			stepY = -1
			sideDistY = (posY - float64(mapY)) * deltaDistY
		} else {
			stepY = 1
			sideDistY = (float64(mapY) + 1.0 - posY) * deltaDistY
		}

		//perform DDA
		for hit == 0 {
			//jump to next map square, OR in x-direction, OR in y-direction
			if sideDistX < sideDistY {
				sideDistX += deltaDistX
				mapX += stepX
				side = 0
			} else {
				sideDistY += deltaDistY
				mapY += stepY
				side = 1
			}

			//Check if ray has hit a wall
			if worldMap[mapX][mapY] > 0 {
				hit = 1
			}
		}

		//Calculate distance projected on camera direction (Euclidean distance will give fisheye effect!)
		if side == 0 {
			perpWallDist = (float64(mapX) - posX + float64(1-stepX)/2) / rayDirX
		} else {
			perpWallDist = (float64(mapY) - posY + float64(1-stepY)/2) / rayDirY
		}

		//Calculate height of line to draw on screen
		var lineHeight int = int(float64(screenHeight) / perpWallDist)

		//calculate lowest and highest pixel to fill in current stripe
		var drawStart int = -lineHeight/2 + screenHeight/2
		if drawStart < 0 {
			drawStart = 0
		}

		var drawEnd int = lineHeight/2 + screenHeight/2
		if drawEnd >= screenHeight {
			drawEnd = screenHeight - 1
		}

		//choose wall color
		var wallColor color.RGBA

		switch worldMap[mapX][mapY] {
		case 1:
			wallColor = color.RGBA{255, 0, 0, 255}
		case 2:
			wallColor = color.RGBA{0, 255, 0, 255}
		case 3:
			wallColor = color.RGBA{0, 0, 255, 255}
		case 4:
			wallColor = color.RGBA{255, 255, 255, 255}
		default:
			wallColor = color.RGBA{255, 255, 0, 255}
		}

		//give x and y sides different brightness
		if side == 1 {
			wallColor = color.RGBA{wallColor.R / 2, wallColor.G / 2, wallColor.B / 2, wallColor.A}
		}

		//draw the pixels of the stripe as a vertical line
		ebitenutil.DrawLine(screen, float64(x), float64(drawStart), float64(x), float64(drawEnd), wallColor)
	}

	return nil
}

func handleMovement() {
	//speed modifiers
	// var moveSpeed float64 = ebiten.CurrentFPS() * 5.0 //the constant value is in squares/second
	// var rotSpeed float64 = ebiten.CurrentFPS() * 3.0  //the constant value is in radians/second

	var moveSpeed float64 = 0.75
	var rotSpeed float64 = 0.1

	//move forward if no wall in front of you
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		if worldMap[int(posX+dirX*moveSpeed)][int(posY)] == 0 {
			posX += dirX * moveSpeed
		}
		if worldMap[int(posX)][int(posY+dirY*moveSpeed)] == 0 {
			posY += dirY * moveSpeed
		}
	}

	//move backwards if no wall behind you
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {

		if worldMap[int(posX-dirX*moveSpeed)][int(posY)] == 0 {
			posX -= dirX * moveSpeed
		}
		if worldMap[int(posX)][int(posY-dirY*moveSpeed)] == 0 {
			posY -= dirY * moveSpeed
		}
	}

	//rotate to the right
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		//both camera direction and camera plane must be rotated
		var oldDirX float64 = dirX
		dirX = dirX*math.Cos(-rotSpeed) - dirY*math.Sin(-rotSpeed)
		dirY = oldDirX*math.Sin(-rotSpeed) + dirY*math.Cos(-rotSpeed)
		var oldPlaneX float64 = planeX
		planeX = planeX*math.Cos(-rotSpeed) - planeY*math.Sin(-rotSpeed)
		planeY = oldPlaneX*math.Sin(-rotSpeed) + planeY*math.Cos(-rotSpeed)
	}

	//rotate to the left
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		//both camera direction and camera plane must be rotated
		var oldDirX float64 = dirX
		dirX = dirX*math.Cos(rotSpeed) - dirY*math.Sin(rotSpeed)
		dirY = oldDirX*math.Sin(rotSpeed) + dirY*math.Cos(rotSpeed)
		var oldPlaneX float64 = planeX
		planeX = planeX*math.Cos(rotSpeed) - planeY*math.Sin(rotSpeed)
		planeY = oldPlaneX*math.Sin(rotSpeed) + planeY*math.Cos(rotSpeed)
	}
}

func main() {
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "raycasting"); err != nil {
		panic(err)
	}
}
