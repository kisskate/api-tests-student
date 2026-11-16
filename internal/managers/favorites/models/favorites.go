package models

import advertisementModels "api-tests-template/internal/managers/advertisements/models"

type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

type Favorite struct {
	FavoriteID    string                            `json:"favorite_id"`
	AddedAt       string                            `json:"added_at"`
	Advertisement advertisementModels.Advertisement `json:"advertisement"`
}

type FavoritesList struct {
	Items      []Favorite `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	Total      int        `json:"total"`
	TotalPages int        `json:"total_pages"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type FavoriteResult struct {
	StatusCode int
	Favorite   Favorite
	Error      ErrorResponse
}

type FavoritesListResult struct {
	StatusCode int
	List       FavoritesList
	Error      ErrorResponse
}
