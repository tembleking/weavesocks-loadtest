package types

type Catalog []CatalogElement

type CatalogElement struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	ImageURL    []string `json:"imageUrl,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Count       int      `json:"count,omitempty"`
	Tag         []string `json:"tag,omitempty"`
}

type CartRequest struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}
