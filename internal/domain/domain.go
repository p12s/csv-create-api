package domain

// Product - product
type Product struct {
	Id    string `csv:"-" json:"id" db:"id"`
	Name  string `csv:"PRODUCT NAME" json:"name" db:"name" binding:"required"`
	Price int    `csv:"PRICE" json:"price" db:"price" binding:"required"`
}

// UpdateProductInput - product update data
type UpdateProductInput struct {
	Name  *string `json:"name" db:"name"`
	Price *int    `json:"price" db:"price"`
}
