package schemas

// DashboardStats overall stats for house owner
type DashboardStats struct {
	TotalHouses    int     `json:"total_houses" example:"10"`
	TotalViews     int     `json:"total_views" example:"1500"`
	TotalLikes     int     `json:"total_likes" example:"120"`
	AveragePrice   float64 `json:"average_price" example:"85000"`
	MedianPrice    float64 `json:"median_price" example:"75000"`
	AverageViews   float64 `json:"average_views" example:"150"`
	AverageLikes   float64 `json:"average_likes" example:"12"`
}

// HouseStatsItem per-house stats
type HouseStatsItem struct {
	ID        int    `json:"id" example:"1"`
	NameEN    string `json:"name_en" example:"Beach House"`
	NameKZ    string `json:"name_kz"`
	NameRU    string `json:"name_ru"`
	Slug      string `json:"slug" example:"beach-house"`
	Price     int    `json:"price" example:"50000"`
	ViewCount int    `json:"view_count" example:"250"`
	LikeCount int    `json:"like_count" example:"15"`
}

// LikesOverTime likes grouped by date for chart
type LikesOverTime struct {
	Date  string `json:"date" example:"2026-03-25"`
	Count int    `json:"count" example:"5"`
}

// ViewsOverTime views grouped by date for chart (based on current snapshot)
type TopHouseStats struct {
	ID        int    `json:"id" example:"1"`
	NameEN    string `json:"name_en" example:"Beach House"`
	NameKZ    string `json:"name_kz"`
	NameRU    string `json:"name_ru"`
	Slug      string `json:"slug"`
	Value     int    `json:"value" example:"250"`
}

// PriceDistribution price range bucket for chart
type PriceDistribution struct {
	Range string `json:"range" example:"0-50000"`
	Count int    `json:"count" example:"3"`
}

// HouseDetailStats detailed stats for a single house
type HouseDetailStats struct {
	HouseID        int              `json:"house_id" example:"1"`
	NameEN         string           `json:"name_en" example:"Beach House"`
	NameKZ         string           `json:"name_kz"`
	NameRU         string           `json:"name_ru"`
	Slug           string           `json:"slug"`
	TotalViews     int              `json:"total_views" example:"350"`
	UniqueViews    int              `json:"unique_views" example:"280"`
	TotalLikes     int              `json:"total_likes" example:"15"`
	TotalBookings  int              `json:"total_bookings" example:"5"`
	RecentViewers  []ViewerInfo     `json:"recent_viewers"`
	RecentLikers   []LikerInfo      `json:"recent_likers"`
	RecentBookings []BookingInfo    `json:"recent_bookings"`
	ViewsPerDay    []LikesOverTime  `json:"views_per_day"`
	LikesPerDay    []LikesOverTime  `json:"likes_per_day"`
}

// ViewerInfo who viewed the house
type ViewerInfo struct {
	UserID    *int   `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	IP        string `json:"ip,omitempty"`
	ViewedAt  string `json:"viewed_at"`
}

// LikerInfo who liked the house
type LikerInfo struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	LikedAt   string `json:"liked_at"`
}

// BookingInfo who booked the house
type BookingInfo struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	GuestCount int    `json:"guest_count"`
	Status     string `json:"status"`
	TotalPrice int    `json:"total_price"`
	CreatedAt  string `json:"created_at"`
}
