package scraper

type dotpropertyItemEntity struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	DatePosted  string `json:"datePosted"`
	Image       string `json:"image"`
	About       struct {
		Category         string `json:"@type"`
		NumberOfBedrooms int    `json:"numberOfBedrooms"`
		ContainedInPlace struct {
			Name string `json:"name"`
		} `json:"containedInPlace"`
		Address struct {
			StreetAddress   string `json:"streetAddress"`
			AddressLocality string `json:"addressLocality"`
			AddressRegion   string `json:"addressRegion"`
		} `json:"address"`
		Geo struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"geo"`
	} `json:"about"`
	Offers struct {
		Price         interface{} `json:"price"`
		PriceCurrency string      `json:"priceCurrency"`
	} `json:"offers"`
}

type dotpropertyItemListEntity struct {
	NumberOfItems   int `json:"numberOfItems"`
	ItemListElement []struct {
		Position int                   `json:"position"`
		Item     dotpropertyItemEntity `json:"item"`
	} `json:"itemListElement"`
}
