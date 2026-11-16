package models

type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       float64
	Quantity    int
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	UserID      string  `json:"user_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Photos      []Photo `json:"photos"`
}

type AdvertisementsList struct {
	Items      []Advertisement `json:"items"`
	HasNext    bool            `json:"has_next"`
	NextCursor string          `json:"next_cursor"`
	Total      int             `json:"total"`
}

type Photo struct {
	ID              string `json:"id"`
	AdvertisementID string `json:"advertisement_id"`
	URL             string `json:"url"`
	SortOrder       int    `json:"sort_order"`
	CreatedAt       string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type AdvertisementResult struct {
	StatusCode    int
	Advertisement Advertisement
	Error         ErrorResponse
}

type PhotosResult struct {
	StatusCode int
	Photos     []Photo
	Error      ErrorResponse
}
