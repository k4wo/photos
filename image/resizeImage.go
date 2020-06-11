package image

import (
	model "photos/model"

	"gopkg.in/gographics/imagick.v3/imagick"
)

const mobileWidth int64 = 640
const tabletWidth int64 = 1280
const displayWidth int64 = 1920

func getDimensions(width, height, resizeTo int64) (int64, int64) {
	if width <= resizeTo {
		return width, height
	}

	percent := float64(resizeTo) * 100 / float64(width)
	h := percent / 100 * float64(height)

	return resizeTo, int64(h)
}

// ResizeImage resize an image
func ResizeImage(image []byte, info model.File, uploadDir string) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	if err := mw.ReadImageBlob(image); err != nil {
		panic(err)
	}

	mobileW, mobileH := getDimensions(info.Width.Int64, info.Height.Int64, displayWidth)
	mw.ResizeImage(uint(mobileW), uint(mobileH), imagick.FILTER_LANCZOS)

	if err := mw.SetImageCompressionQuality(75); err != nil {
		panic(err)
	}

	mw.WriteImage(uploadDir + info.Hash.String + "_mobile")
}
