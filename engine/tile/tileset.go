package tile

import (
	"encoding/json"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// every tileset must be able to give an image given an id
type Tileset interface {
	Img(id int) *ebiten.Image
	Contains(id int) bool
}

// the tileset data deserialized from a standard, single-image tileset
type UniformTilesetJSON struct {
	Path       string `json:"image"`
	TileWidth  int    `json:"tilewidth"`
	TileHeight int    `json:"tileheight"`
}

// the front-facing tileset object used for single-image tilesets
type UniformTileset struct {
	img        *ebiten.Image
	gid        int
	tileWidth  int
	tileHeight int
}

//TODO make dinamic tileset

func (u *UniformTileset) Img(id int) *ebiten.Image {
	id -= u.gid

	tilesPerRow := u.img.Bounds().Dx() / u.tileWidth
	srcX := (id % tilesPerRow) * u.tileWidth
	srcY := (id / tilesPerRow) * u.tileHeight

	return u.img.SubImage(
		image.Rect(
			srcX, srcY, srcX+u.tileWidth, srcY+u.tileHeight,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynTileset) Img(id int) *ebiten.Image {
	id -= d.gid

	return d.imgs[id]
}

func (u *UniformTileset) Contains(id int) bool {
	id -= u.gid
	tilesPerRow := u.img.Bounds().Dx() / u.tileWidth
	tilesPerCol := u.img.Bounds().Dy() / u.tileHeight
	return id >= 0 && id < tilesPerRow*tilesPerCol
}

func (d *DynTileset) Contains(id int) bool {
	id -= d.gid
	return id >= 0 && id < len(d.imgs)
}

func NewTileset(path string, gid int) (Tileset, error) {
	// read file contents
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		// return dyn tileset
		var dynTilesetJSON DynTilesetJSON
		err = json.Unmarshal(contents, &dynTilesetJSON)
		if err != nil {
			return nil, err
		}

		// create the tileset
		dynTileset := DynTileset{}
		dynTileset.gid = gid
		dynTileset.imgs = make([]*ebiten.Image, 0)

		// loop over tile data and load image for each
		for _, tileJSON := range dynTilesetJSON.Tiles {

			// convert tileset relative path to root relative path
			tileJSONPath := tileJSON.Path
			tileJSONPath = filepath.Clean(tileJSONPath)
			tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
			tileJSONPath = filepath.Join("assets/", tileJSONPath)

			img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
			if err != nil {
				return nil, err
			}

			dynTileset.imgs = append(dynTileset.imgs, img)
		}

		return &dynTileset, nil
	}
	// return uniform tileset
	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(contents, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	// convert tileset relative path to root relative path
	tileJSONPath := uniformTilesetJSON.Path
	tileJSONPath = filepath.Clean(tileJSONPath)
	tileJSONPath = strings.ReplaceAll(tileJSONPath, "\\", "/")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = strings.TrimPrefix(tileJSONPath, "../")
	tileJSONPath = filepath.Join("assets/", tileJSONPath)

	img, _, err := ebitenutil.NewImageFromFile(tileJSONPath)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{
		img:        img,
		gid:        gid,
		tileWidth:  uniformTilesetJSON.TileWidth,
		tileHeight: uniformTilesetJSON.TileHeight,
	}

	return &uniformTileset, nil
}
