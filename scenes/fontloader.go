package scenes

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var DefaultFontFace text.Face

func LoadFont() {
	fontBytes, err := os.ReadFile("assets/font/SuperMystery-BLG4x.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	DefaultFontFace = text.NewGoXFace(face)
}
