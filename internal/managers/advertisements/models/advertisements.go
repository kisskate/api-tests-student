package models

type Advertisement struct {
	Title string `json:"title"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}
