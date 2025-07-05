package systems

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/deglan/horrorchain/constants"
	"github.com/deglan/horrorchain/engine/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Position struct {
	X, Y int
}

func MovePlayer(player *entities.Player, colliders []image.Rectangle) {
	player.Dx = 0
	player.Dy = 0

	playerSpeed := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.Dx = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.Dx = playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		player.Dy = -playerSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		player.Dy = playerSpeed
	}

	player.X += player.Dx
	CheckCollisionHorizontal(player.Sprite, colliders)

	player.Y += player.Dy
	CheckCollisionVertical(player.Sprite, colliders)
}

func MoveEnemies(enemies []*entities.Enemy, player *entities.Player, colliders []image.Rectangle, grid [][]bool) {
	shuffleEnemies(enemies)
	for i, enemy := range enemies {
		if !enemy.FollowsPlayer {
			continue
		}

		tempGrid := copyGrid(grid)
		for j, other := range enemies {
			if i == j {
				continue
			}
			ox, oy := worldToGrid(other.X, other.Y)
			if oy >= 0 && oy < len(tempGrid) && ox >= 0 && ox < len(tempGrid[0]) {
				tempGrid[oy][ox] = true
			}
		}

		dx := player.X - enemy.X
		dy := player.Y - enemy.Y
		dist := math.Hypot(dx, dy)
		if dist == 0 {
			continue
		}

		dirX := dx / dist
		dirY := dy / dist

		speed := 1.0
		newX := enemy.X + dirX*speed
		newY := enemy.Y + dirY*speed

		ex, ey := worldToGrid(newX, newY)
		px, py := worldToGrid(player.X, player.Y)

		if ex < 0 || ey < 0 || ex >= len(grid[0]) || ey >= len(grid) {
			continue
		}

		path := bfs(tempGrid, ex, ey, px, py)

		if len(path) > 1 {
			next := path[1]
			nextX, nextY := gridToWorld(next.X, next.Y)

			dirX := nextX - enemy.X
			dirY := nextY - enemy.Y
			dist := math.Hypot(dirX, dirY)
			if dist > 0 {
				moveX := dirX / dist
				moveY := dirY / dist
				speed := 1.0

				newX := enemy.X + moveX*speed
				newY := enemy.Y + moveY*speed

				newRectX := image.Rect(
					int(newX), int(enemy.Y),
					int(newX)+constants.Tilesize,
					int(enemy.Y)+constants.Tilesize,
				)

				newRectY := image.Rect(
					int(enemy.X), int(newY),
					int(enemy.X)+constants.Tilesize,
					int(newY)+constants.Tilesize,
				)

				if !isBlocked(newRectX, enemies, colliders, i) {
					enemy.X = newX
				} else {
					bounce := 0.5
					enemy.X -= moveX * bounce
				}

				if !isBlocked(newRectY, enemies, colliders, i) {
					enemy.Y = newY
				} else {
					bounce := 0.5
					enemy.Y -= moveY * bounce
				}

				enemy.Dx = math.Copysign(1, moveX)
				enemy.Dy = math.Copysign(1, moveY)

			}
		}

	}
}

func shuffleEnemies(enemies []*entities.Enemy) {
	rand.Shuffle(len(enemies), func(i, j int) {
		enemies[i], enemies[j] = enemies[j], enemies[i]
	})
}

func isBlocked(newRect image.Rectangle, enemies []*entities.Enemy, colliders []image.Rectangle, i int) bool {
	collides := false

	for _, col := range colliders {
		if col.Overlaps(newRect) {
			collides = true
			break
		}
	}

	for j, other := range enemies {
		if i == j {
			continue
		}
		otherRect := image.Rect(
			int(other.X), int(other.Y),
			int(other.X)+constants.Tilesize,
			int(other.Y)+constants.Tilesize,
		)
		if otherRect.Overlaps(newRect) {
			collides = true
			break
		}
	}

	return collides
}

func bfs(grid [][]bool, startX, startY, goalX, goalY int) []Position {
	type Node struct {
		X, Y int
		Prev *Node
	}

	visited := make([][]bool, len(grid))
	for i := range visited {
		visited[i] = make([]bool, len(grid[0]))
	}
	queue := []Node{{X: startX, Y: startY, Prev: nil}}
	visited[startY][startX] = true

	dirs := []Position{
		{X: 0, Y: -1}, // up
		{X: 1, Y: 0},  // right
		{X: 0, Y: 1},  // down
		{X: -1, Y: 0}, // left
	}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.X == goalX && node.Y == goalY {
			path := []Position{{X: node.X, Y: node.Y}}
			for node.Prev != nil {
				node = *node.Prev
				path = append([]Position{{X: node.X, Y: node.Y}}, path...)
			}
			return path
		}

		for _, dir := range dirs {
			newX := node.X + dir.X
			newY := node.Y + dir.Y

			if newX < 0 || newX >= len(grid[0]) || newY < 0 || newY >= len(grid) {
				continue
			}

			if grid[newY][newX] || visited[newY][newX] {
				continue
			}

			visited[newY][newX] = true
			queue = append(queue, Node{X: newX, Y: newY, Prev: &node})
		}
	}

	return nil
}

func copyGrid(grid [][]bool) [][]bool {
	copy := make([][]bool, len(grid))
	for y := range grid {
		copy[y] = make([]bool, len(grid[0]))
		for x := range grid[0] {
			copy[y][x] = grid[y][x]
		}
	}
	return copy
}

func worldToGrid(x, y float64) (int, int) {
	return int(x / constants.Tilesize), int(y / constants.Tilesize)
}

func gridToWorld(x, y int) (float64, float64) {
	return float64(x * constants.Tilesize), float64(y * constants.Tilesize)
}

func CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize,
			),
		) {
			if sprite.Dx > 0.0 {
				sprite.X = float64(collider.Min.X) - constants.Tilesize
			} else if sprite.Dx < 0.0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	for _, collider := range colliders {
		if collider.Overlaps(
			image.Rect(
				int(sprite.X),
				int(sprite.Y),
				int(sprite.X)+constants.Tilesize,
				int(sprite.Y)+constants.Tilesize,
			),
		) {
			if sprite.Dy > 0.0 {
				sprite.Y = float64(collider.Min.Y) - constants.Tilesize
			} else if sprite.Dy < 0.0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func DrawEnemyDebugRects(screen *ebiten.Image, enemies []*entities.Enemy, cameraX, cameraY float64) {
	for _, e := range enemies {
		rect := image.Rect(int(e.X), int(e.Y), int(e.X)+constants.Tilesize, int(e.Y)+constants.Tilesize)
		vector.StrokeRect(screen,
			float32(rect.Min.X)+float32(cameraX),
			float32(rect.Min.Y)+float32(cameraY),
			float32(rect.Dx()), float32(rect.Dy()),
			1.0, color.RGBA{255, 0, 0, 255}, true)
	}
}
