package schemas

type Image struct {
	ID        int    `json:"id"`
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail"`
	MimeType  string `json:"mime_type"`
	Size      *int   `json:"size"`
	HouseID   *int   `json:"house_id"`
}
