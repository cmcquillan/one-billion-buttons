package main

import (
	"image"
	"image/color"
)

func ColorizeImage(img image.Image, hex string) (image.Image, error) {
	bounds := img.Bounds()

	new := image.NewRGBA(bounds)
	newRgb, err := HexToBytes(hex)

	if err != nil {
		return nil, err
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r >>= 8
			g >>= 8
			b >>= 8

			// if r > 0 && b > 0 && g > 0 {
			// 	log.Printf("(%d, %d) -> %d %d %d", x, y, r, g, b)
			// }

			if r > 0 || b > 0 || g > 0 {
				newColor := color.NRGBA{
					R: newRgb[0],
					G: newRgb[1],
					B: newRgb[2],
					A: uint8(a),
				}

				new.Set(x, y, newColor)
			} else {
				new.Set(x, y, color.NRGBA{
					R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a),
				})
			}
		}
	}

	sub := new.SubImage(bounds)
	return sub, nil
}
