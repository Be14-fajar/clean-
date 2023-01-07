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
func (bd *bookData) Update(bookID int, updatedData book.Core) (book.Core, error) {
	BooksModel := CoreToData(updatedData)
	BooksModel.ID = uint(bookID)

	Input := bd.db.Where("id = ?", bookID).Updates(&BooksModel)
	if Input.Error != nil {
		log.Println("Get By ID query error", Input.Error.Error())
		return book.Core{}, Input.Error
	}
	if Input.RowsAffected <= 0 {
		return book.Core{}, errors.New("Not found")
	}

	return ToCore(BooksModel), nil
}

//	func (bd *bookData) Delete(bookID int, userID int) error {
//		return nil
//	}
func (bd *bookData) MyBook(userID int) ([]book.Core, error) {
	myBooks := []BookPemilik{}
	err := bd.db.Raw("SELECT books.id, books.judul, books.tahun_terbit, books.penulis, users.name FROM books JOIN users ON users.id = books.user_id WHERE books.user_id = ?", userID).Find(&myBooks).Error
	if err != nil {
		return nil, err
	}

	var dataCore = ListModelTOCore(myBooks)

	return dataCore, nil
}

// All implements book.BookData
func (bd *bookData) All() ([]book.Core, error) {
	var buku []BookPemilik
	fmt.Println("ini query", buku)
	// tx := bd.db.Model(&buku).Select("books.id, books.judul, books.tahun_terbit, books.penulis, users.name").Joins("users on users.id = books.user_id").Where("books.user_id").Find(buku)
	tx := bd.db.Select("books.id, books.judul, books.tahun_terbit, books.penulis, users.name FROM books JOIN users ON users.id = books.user_id WHERE books.user_id ").Find(buku)
	if tx != nil {
		return nil, tx.Error
	}
	var dataCore = ListModelTOCore(buku)
	fmt.Println("ini data core", dataCore)

	return dataCore, nil
}
