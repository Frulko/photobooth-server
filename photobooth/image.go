package PhotoBooth

import (
	"gopkg.in/gographics/imagick.v2/imagick"
)

type (
	Image struct {

	}
)


func (c *Image) Prepare(blob []byte) ([]byte){
	imagick.Initialize()

	// Schedule cleanup
	defer imagick.Terminate()
	var err error

	mw := imagick.NewMagickWand()

	err = mw.ReadImageBlob(blob)
	if err != nil {
		panic(err)
	}

	// Get original logo size
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	// Calculate half the size
	hWidth := uint(width / 4)
	hHeight := uint(height / 4)


	target := imagick.NewPixelWand()
	target.SetColor("white")

	/*err = mw.RotateImage(target, -90.0)
	if err != nil {
		panic(err)
	}*/
	// Resize the image using the Lanczos filter
	// The blur factor is a float, where > 1 is blurry, < 1 is sharp
	err = mw.ResizeImage(hWidth, hHeight, imagick.FILTER_LANCZOS, 1)
	if err != nil {
		panic(err)
	}
	

	// Set the compression quality to 95 (high quality = low compression)
	err = mw.SetImageCompressionQuality(70)
	if err != nil {
		panic(err)
	}

	return mw.GetImagesBlob()
}