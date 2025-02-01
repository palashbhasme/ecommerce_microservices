package response

type ProductVariantResponse struct {
	ID         string  `json:"id"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	StockLevel int     `json:"stock_level"`
	Size       string  `json:"size"`
	Color      string  `json:"color"`
}

type ProductResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	SKU         string                   `json:"sku"`
	Category    string                   `json:"category"`
	StockLevel  int                      `json:"stock_level"`
	Variants    []ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
}
