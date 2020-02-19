package image

import (
	model "photos/model"

	"gopkg.in/gographics/imagick.v3/imagick"
)

const mobileWidth uint = 640
const tabletWidth uint = 1280
const displayWidth uint = 1920

func getDimensions(width, height, resizeTo uint) (uint, uint) {
	if width <= resizeTo {
		return width, height
	}

	percent := float64(resizeTo) * 100 / float64(width)
	h := percent / 100 * float64(height)

	return resizeTo, uint(h)
}

// ResizeImage resize an image
func ResizeImage(image []byte, info model.File, uploadDir string) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	if err := mw.ReadImageBlob(image); err != nil {
		panic(err)
	}

	mobileW, mobileH := getDimensions(info.Width, info.Height, displayWidth)
	mw.ResizeImage(mobileW, mobileH, imagick.FILTER_LANCZOS)

	if err := mw.SetImageCompressionQuality(75); err != nil {
		panic(err)
	}

	mw.WriteImage(uploadDir + info.Hash + "_mobile")
}
