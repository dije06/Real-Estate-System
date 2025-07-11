package models

type Listing struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int    `json:"user_id"`
	Price       int    `json:"price"`
	ListingType string `json:"listing_type"` // rent or sale
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
