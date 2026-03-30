package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/repositories"
	"github.com/nurkenspashev92/bookit/pkg/aws"
)

const maxHouseImages = 15

type ImageService struct {
	repository *repositories.HouseImageRepository
	s3         *aws.AwsS3Client
}

func NewImageService(repo *repositories.HouseImageRepository, s3 *aws.AwsS3Client) *ImageService {
	return &ImageService{
		repository: repo,
		s3:         s3,
	}
}

type uploadedImage struct {
	key      string
	thumbKey string
	size     int
	mime     string
	width    int
	height   int
}

func (s *ImageService) UploadHouseImages(ctx context.Context, houseID int, files []*multipart.FileHeader) error {
	count, err := s.repository.CountByHouse(ctx, houseID)
	if err != nil {
		return err
	}

	if count+len(files) > maxHouseImages {
		return ErrMaxImagesExceeded
	}

	// Phase 1: Process images (CPU-bound) — parallel, no limit needed
	type processed struct {
		result *Result
		err    error
	}
	procResults := make([]processed, len(files))

	procGroup, _ := errgroup.WithContext(ctx)
	for idx, file := range files {
		procGroup.Go(func() error {
			img, err := Process(file)
			procResults[idx] = processed{result: img, err: err}
			return err
		})
	}
	if err := procGroup.Wait(); err != nil {
		return err
	}

	// Phase 2: Upload to S3 — limit concurrency to avoid overwhelming S3
	results := make([]uploadedImage, len(files))
	uploadGroup, ctx := errgroup.WithContext(ctx)
	uploadGroup.SetLimit(6)

	for idx, pr := range procResults {
		img := pr.result
		ts := time.Now().UnixNano() + int64(idx)
		originalKey := fmt.Sprintf("houses/original/%d_%d.jpg", houseID, ts)
		thumbKey := fmt.Sprintf("houses/thumbnail/%d_%d.webp", houseID, ts)

		// Upload original
		uploadGroup.Go(func() error {
			_, err := s.s3.UploadCompressed(ctx, originalKey, img.Original, "image/jpeg")
			return err
		})

		// Upload thumbnail
		uploadGroup.Go(func() error {
			_, err := s.s3.UploadCompressed(ctx, thumbKey, img.Thumbnail, "image/webp")
			return err
		})

		results[idx] = uploadedImage{
			key:      originalKey,
			thumbKey: thumbKey,
			size:     img.Size,
			mime:     img.Mime,
			width:    img.Width,
			height:   img.Height,
		}
	}

	if err := uploadGroup.Wait(); err != nil {
		for _, r := range results {
			_ = s.s3.Delete(context.Background(), r.key)
			_ = s.s3.Delete(context.Background(), r.thumbKey)
		}
		return err
	}

	// Phase 3: Batch insert into DB — single query
	images := make([]models.Image, 0, len(results))
	for _, r := range results {
		images = append(images, models.Image{
			Original:  r.key,
			Thumbnail: r.thumbKey,
			MimeType:  r.mime,
			Width:     &r.width,
			Height:    &r.height,
			Size:      &r.size,
			HouseID:   &houseID,
		})
	}

	if err := s.repository.CreateBatch(ctx, images); err != nil {
		for _, r := range results {
			_ = s.s3.Delete(context.Background(), r.key)
			_ = s.s3.Delete(context.Background(), r.thumbKey)
		}
		return fmt.Errorf("db save failed: %w", err)
	}

	return nil
}

func (s *ImageService) DeleteHouseImage(ctx context.Context, imageID int) error {
	keys, err := s.repository.DeleteReturningKeys(ctx, imageID)
	if err != nil {
		return ErrImageNotFound
	}

	// Delete S3 objects in parallel
	g, ctx := errgroup.WithContext(ctx)
	if keys.Original != "" {
		g.Go(func() error {
			return s.s3.Delete(ctx, keys.Original)
		})
	}
	if keys.Thumbnail != "" {
		g.Go(func() error {
			return s.s3.Delete(ctx, keys.Thumbnail)
		})
	}
	_ = g.Wait()

	return nil
}
