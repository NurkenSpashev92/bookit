package schemas

// Image image response DTO
type Image struct {
	ID        int    `json:"id" example:"1"`
	Original  string `json:"original" example:"https://bucket.s3.eu-central-1.amazonaws.com/houses/original/1_123.jpg"`
	Thumbnail string `json:"thumbnail" example:"https://bucket.s3.eu-central-1.amazonaws.com/houses/thumbnail/1_123.webp"`
	MimeType  string `json:"mime_type" example:"image/jpeg"`
	Size      *int   `json:"size" example:"204800"`
	HouseID   *int   `json:"house_id" example:"1"`
}
