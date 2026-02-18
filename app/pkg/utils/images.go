package utils

import (
	"github.com/nurkenspashev92/bookit/configs"
	"github.com/nurkenspashev92/bookit/internal/models"
	"github.com/nurkenspashev92/bookit/internal/schemas"
)

func FillHouseImagesURL(cfg *configs.AwsConfig, images []models.Image) {
	for i := range images {
		if images[i].Original != "" {
			images[i].Original = cfg.AwsS3URL(images[i].Original)
		}
		if images[i].Thumbnail != "" {
			images[i].Thumbnail = cfg.AwsS3URL(images[i].Thumbnail)
		}
	}
}

func FillHouseListImagesURL(cfg *configs.AwsConfig, houses []schemas.HouseListItem) {
	for i := range houses {
		for j := range houses[i].Images {
			if houses[i].Images[j].Original != "" {
				houses[i].Images[j].Original = cfg.AwsS3URL(houses[i].Images[j].Original)
			}
			if houses[i].Images[j].Thumbnail != "" {
				houses[i].Images[j].Thumbnail = cfg.AwsS3URL(houses[i].Images[j].Thumbnail)
			}
		}
	}
}
