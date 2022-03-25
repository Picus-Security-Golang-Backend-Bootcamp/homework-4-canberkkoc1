package handlers

import (
	"ck/helper"
	"ck/migration"
	"ck/models"
	"errors"
)

func GetAllBook() ([]models.Books, error) {

	var book []models.Books

	result := migration.DB.Find(&book)

	if result.Error != nil {
		return book, result.Error
	}

	return book, nil

}

func GetBookByName(name string) ([]models.Books, error) {

	var books []models.Books

	result := migration.DB.Where(" Name LIKE ?", "%"+name+"%").Find(&books)
	if result.Error != nil {
		return books, result.Error
	}

	return books, nil
}

func UpdateStock(id, stock int) ([]models.Books, error) {
	var books []models.Books
	var stoc_num int

	migration.DB.Find(&books)

	result := migration.DB.Table("books").Select("stock_number").Where("id = ?", id).Scan(&stoc_num)

	if stoc_num <= 0 || stoc_num < stock {
		return books, result.Error
	}

	if result.Error != nil {
		return books, result.Error
	}

	newStock := stoc_num - stock

	migration.DB.Model(&books).Where("id = ?", id).Update("stock_number", newStock)

	return books, nil
}

func DeleteBookById(id int) ([]models.Books, error) {
	var books []models.Books

	var n []int

	migration.DB.Model(&books).Pluck("id", &n)

	migration.DB.Unscoped().Delete(&books, id)

	isDeleted := helper.CheckSlice(n, id)

	if !isDeleted {

		return books, errors.New("id not found")
	}

	migration.DB.Find(&books)

	return books, nil
}
