package tile

import (
	"encoding/json"
	"image"
	"os"
	"path"

	"github.com/deglan/horrorchain/constants"
)

type TilemapLayerJSON struct {
	Data   []int  `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Name   string `json:"name"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
	// raw data for each tileset (path, gid)
	Tilesets []map[string]any `json:"tilesets"`
}

func (t *TilemapJSON) GenTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, 0)

	for _, tilesetData := range t.Tilesets {
		// convert map relative path to project relative path
		tilesetPath := path.Join("assets/maps/", tilesetData["source"].(string))
		tileset, err := NewTileset(tilesetPath, int(tilesetData["firstgid"].(float64)))
		if err != nil {
			return nil, err
		}

		tilesets = append(tilesets, tileset)
	}

	return tilesets, nil
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}

func (t *TilemapJSON) ExtractCollidersFromLayer(layerName string) []image.Rectangle {
	var colliders []image.Rectangle
	var layer *TilemapLayerJSON
	for _, l := range t.Layers {
		if l.Name == layerName {
			layer = &l
			break
		}
	}
	if layer == nil {
		return colliders
	}

	for i, tileID := range layer.Data {
		if tileID == 0 {
			continue
		}
		x := i % layer.Width
		y := i / layer.Width

		rect := image.Rect(
			x*constants.Tilesize,
			y*constants.Tilesize,
			(x+1)*constants.Tilesize,
			(y+1)*constants.Tilesize,
		)
		colliders = append(colliders, rect)
	}
	return colliders
}

func (t *TilemapJSON) BuildGrid(solidLayerNames ...string) [][]bool {
	var width, height int

	for _, name := range solidLayerNames {
		for _, layer := range t.Layers {
			if layer.Name == name {
				width = layer.Width
				height = layer.Height
				break
			}
		}
		if width > 0 && height > 0 {
			break
		}
	}

	if width == 0 || height == 0 {
		return nil
	}

	grid := make([][]bool, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]bool, width)
	}

	for _, name := range solidLayerNames {
		for _, layer := range t.Layers {
			if layer.Name != name {
				continue
			}
			for i, tileID := range layer.Data {
				if tileID == 0 {
					continue
				}
				x := i % layer.Width
				y := i / layer.Width
				grid[y][x] = true
			}
		}
	}

	return grid
}
