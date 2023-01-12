package data

import (
	"api/features/book"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type bookData struct {
	db *gorm.DB
}

// Delete implements book.BookData

func New(db *gorm.DB) book.BookData {
	return &bookData{
		db: db,
	}
}

func (bd *bookData) Add(userID int, newBook book.Core) (book.Core, error) {
	cnv := CoreToData(newBook)
	cnv.UserID = uint(userID)
	err := bd.db.Create(&cnv).Error
	if err != nil {
		return book.Core{}, err
	}

	newBook.ID = cnv.ID

	return newBook, nil
}
func (bd *bookData) Update(userID int, bookID int, updatedData book.Core) (book.Core, error) {
	cnv := CoreToData(updatedData)

	// DB Update(value)
	tx := bd.db.Where("id = ? && user_id = ?", bookID, userID).Updates(&cnv)
	if tx.Error != nil {
		log.Println("update book query error :", tx.Error)
		return book.Core{}, tx.Error
	}

	// Rows affected checking
	if tx.RowsAffected <= 0 {
		log.Println("update book query error : data not found")
		return book.Core{}, errors.New("not found")
	}

	// return result converting cnv to book.Core
	return ToCore(cnv), nil
}

//	func (bd *bookData) Delete(bookID int, userID int) error {
//		return nil
//	}
func (bd *bookData) MyBook(userID int) ([]book.Core, error) {
	var myBooks []BookPemilik
	err := bd.db.Raw("SELECT books.id, books.judul, books.tahun_terbit, books.penulis, users.name FROM books JOIN users ON users.id = books.user_id WHERE books.user_id = ?", userID).Find(&myBooks).Error
	if err != nil {
		return nil, err
	}

	var dataCore = ListModelTOCore(myBooks)

	return dataCore, nil
}

// All implements book.BookData
func (bd *bookData) AllBook() ([]book.Core, error) {
	var buku []BookPemilik
	fmt.Println("ini query", buku)
	tx := bd.db.Raw("SELECT books.id, books.judul, books.tahun_terbit, books.penulis, users.name FROM books JOIN users ON users.id = books.user_id WHERE books.deleted_at IS NULL").Find(&buku)

	fmt.Println("ini TX", tx)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var dataCore = ListModelTOCore(buku)

	return dataCore, nil
}
func (bd *bookData) Delete(userID int, bookID int) error {
	buku := Books{}
	del := bd.db.Where("id = ? AND user_id = ?", bookID, userID).Delete(&buku, bookID)
	if del.Error != nil {
		log.Println("delete book query error :", del.Error)
		return del.Error
	}
	if del.RowsAffected <= 0 {
		log.Println("delete book query error : data not found")
		return errors.New("not found")
	}

	return nil
}
