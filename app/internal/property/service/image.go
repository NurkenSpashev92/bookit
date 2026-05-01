package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/nurkenspashev92/bookit/internal/property/model"
	"github.com/nurkenspashev92/bookit/internal/property/repository"
	"github.com/nurkenspashev92/bookit/pkg/aws"
	"github.com/nurkenspashev92/bookit/pkg/cache"
	"github.com/nurkenspashev92/bookit/pkg/imageproc"
)

const maxHouseImages = 15

// HouseImageRepository describes the persistence contract ImageService depends on.
type HouseImageRepository interface {
	GetHouseIDBySlug(ctx context.Context, slug string) (int, error)
	CountByHouse(ctx context.Context, houseID int) (int, error)
	CreateBatch(ctx context.Context, images []model.Image) error
	DeleteReturningKeys(ctx context.Context, imageID int) (*repository.ImageKeys, error)
}

type ImageService struct {
	repository HouseImageRepository
	s3         *aws.AwsS3Client
	cache      *cache.Cache
}

func NewImageService(repo HouseImageRepository, s3 *aws.AwsS3Client, c *cache.Cache) *ImageService {
	return &ImageService{
		repository: repo,
		cache:      c,
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

func (s *ImageService) UploadHouseImages(ctx context.Context, slug string, files []*multipart.FileHeader) error {
	houseID, err := s.repository.GetHouseIDBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("house not found")
	}

	count, err := s.repository.CountByHouse(ctx, houseID)
	if err != nil {
		return err
	}

	if count+len(files) > maxHouseImages {
		return ErrMaxImagesExceeded
	}

	type processed struct {
		result *imageproc.Result
		err    error
	}
	procResults := make([]processed, len(files))

	procGroup, _ := errgroup.WithContext(ctx)
	for idx, file := range files {
		procGroup.Go(func() error {
			img, err := imageproc.Process(file)
			procResults[idx] = processed{result: img, err: err}
			return err
		})
	}
	if err := procGroup.Wait(); err != nil {
		return err
	}

	results := make([]uploadedImage, len(files))
	uploadGroup, uploadCtx := errgroup.WithContext(ctx)
	uploadGroup.SetLimit(6)

	for idx, pr := range procResults {
		img := pr.result
		ts := time.Now().UnixNano() + int64(idx)
		originalKey := fmt.Sprintf("houses/original/%d_%d.jpg", houseID, ts)
		thumbKey := fmt.Sprintf("houses/thumbnail/%d_%d.webp", houseID, ts)

		uploadGroup.Go(func() error {
			_, err := s.s3.UploadCompressed(uploadCtx, originalKey, img.Original, "image/jpeg")
			return err
		})

		uploadGroup.Go(func() error {
			_, err := s.s3.UploadCompressed(uploadCtx, thumbKey, img.Thumbnail, "image/webp")
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

	images := make([]model.Image, 0, len(results))
	for _, r := range results {
		images = append(images, model.Image{
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

	s.cache.DeleteByPrefix("houses:")
	return nil
}

func (s *ImageService) DeleteHouseImage(ctx context.Context, imageID int) error {
	keys, err := s.repository.DeleteReturningKeys(ctx, imageID)
	if err != nil {
		return ErrImageNotFound
	}

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

	s.cache.DeleteByPrefix("houses:")
	return nil
}
