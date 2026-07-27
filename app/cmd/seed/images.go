package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/imageproc"
)

type preparedImage struct {
	original  []byte
	thumbnail []byte
	mime      string
	width     int
	height    int
	size      int
}

type imageRow struct {
	houseID   int
	original  string
	thumbnail string
	mime      string
	width     int
	height    int
	size      int
}

func prepareImages(dir string) ([]preparedImage, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read images dir %q: %w", dir, err)
	}

	var prepared []preparedImage

	for _, e := range entries {
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			continue
		}

		path := filepath.Join(dir, e.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}

		result, err := imageproc.ProcessBytes(data, ext)
		if err != nil {
			return nil, fmt.Errorf("process %q: %w", path, err)
		}

		prepared = append(prepared, preparedImage{
			original:  result.Original,
			thumbnail: result.Thumbnail,
			mime:      result.Mime,
			width:     result.Width,
			height:    result.Height,
			size:      result.Size,
		})
	}

	if len(prepared) == 0 {
		return nil, fmt.Errorf("no image files found in %s", dir)
	}

	return prepared, nil
}

func uploadHouseImages(
	ctx context.Context, s3client *aws.AwsS3Client, houseIDs []int, images []preparedImage,
) []imageRow {
	type slot struct {
		row imageRow
		ok  bool
	}

	slots := make([]slot, len(houseIDs)*imagesPerHouse)

	var failed int64

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(uploadConcurrency)

	for houseIdx, houseID := range houseIDs {
		for j := range imagesPerHouse {
			slotIdx := houseIdx*imagesPerHouse + j
			img := images[j%len(images)]

			g.Go(func() error {
				originalKey := fmt.Sprintf("houses/original/%d_%d.jpg", houseID, j)
				thumbKey := fmt.Sprintf("houses/thumbnail/%d_%d.webp", houseID, j)

				if _, err := s3client.UploadCompressed(gctx, originalKey, img.original, "image/jpeg"); err != nil {
					log.Printf("[house %d] upload original %d: %v", houseID, j+1, err)
					atomic.AddInt64(&failed, 1)
					return nil
				}

				if _, err := s3client.UploadCompressed(gctx, thumbKey, img.thumbnail, "image/webp"); err != nil {
					log.Printf("[house %d] upload thumbnail %d: %v", houseID, j+1, err)
					atomic.AddInt64(&failed, 1)
					return nil
				}

				slots[slotIdx] = slot{
					row: imageRow{
						houseID:   houseID,
						original:  originalKey,
						thumbnail: thumbKey,
						mime:      img.mime,
						width:     img.width,
						height:    img.height,
						size:      img.size,
					},
					ok: true,
				}

				return nil
			})
		}
	}

	_ = g.Wait()

	rows := make([]imageRow, 0, len(slots))
	for _, s := range slots {
		if s.ok {
			rows = append(rows, s.row)
		}
	}

	if failed > 0 {
		log.Printf("Uploads failed: %d of %d", failed, len(slots))
	}

	return rows
}
