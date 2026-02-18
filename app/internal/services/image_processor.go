package services

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

type Result struct {
	Original  []byte
	Thumbnail []byte
	Width     int
	Height    int
	Mime      string
	Size      int
}

func Process(file *multipart.FileHeader) (*Result, error) {

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	img, format, err := image.Decode(src)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	// --- THUMB ---
	thumb := imaging.Fit(img, 400, 400, imaging.Lanczos)

	var thumbBuf bytes.Buffer
	err = webp.Encode(&thumbBuf, thumb, &webp.Options{
		Lossless: false,
		Quality:  75,
	})
	if err != nil {
		return nil, err
	}

	// --- ORIGINAL COMPRESS ---
	var originalBuf bytes.Buffer

	switch format {
	case "png":
		err = png.Encode(&originalBuf, img)
	default:
		err = jpeg.Encode(&originalBuf, img, &jpeg.Options{Quality: 85})
	}
	if err != nil {
		return nil, err
	}

	return &Result{
		Original:  originalBuf.Bytes(),
		Thumbnail: thumbBuf.Bytes(),
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		Mime:      "image/webp",
		Size:      originalBuf.Len(),
	}, nil
}
