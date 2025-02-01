package request

type ProductRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Description string                  `json:"description" binding:"required"`
	CategoryID  *string                 `json:"category_id"`
	SKU         string                  `json:"sku" binding:"required"`
	StockLevel  int                     `json:"stock_quantity"`
	Brand       string                  `json:"brand"`
	Variants    []ProductVariantRequest `json:"variants" binding:"required"`
}

type ProductVariantRequest struct {
	Color      string   `json:"color"`
	Size       string   `json:"size"`
	ProductID  string   `json:"product_id" binding:"required"`
	Price      *float64 `json:"price" binding:"required"`
	StockLevel int      `json:"stock_quantity"`
	SKU        string   `json:"sku" binding:"required"`
}
type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required"`
}
