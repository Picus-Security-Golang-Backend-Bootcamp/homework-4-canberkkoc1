package models

import (
	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	StockNumber int     `json:"stock_num"`
	PageNumber  int     `json:"page_num"`
	Price       float64 `json:"price"`
	Name        string  `json:"book_name"`
	StockCode   string  `json:"stock_code"`
	Isbn        string  `json:"isbn"`
	AuthorName  string  `json:"author_name"`
}
