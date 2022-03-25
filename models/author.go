package models

import "gorm.io/gorm"

type Author struct {
	gorm.Model
	AuthorName string `json:"author_name"`
}
