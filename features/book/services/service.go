package services

import (
	"api/features/book"
	"api/helper"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
)

type bookSrv struct {
	data     book.BookData
	validasi *validator.Validate
}

// Delete implements book.BookService

// Update implements book.BookService

func New(d book.BookData) book.BookService {
	return &bookSrv{
		data:     d,
		validasi: validator.New(),
	}
}

func (bs *bookSrv) Add(token interface{}, newBook book.Core) (book.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return book.Core{}, errors.New("user not found")
	}

	err := bs.validasi.Struct(newBook)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
		}
		return book.Core{}, errors.New("validation error")
	}

	res, err := bs.data.Add(userID, newBook)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "user not found"
		} else {
			msg = "internal server error"
		}
		return book.Core{}, errors.New(msg)
	}

	return res, nil

}

func (bs *bookSrv) MyBook(token interface{}) ([]book.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return nil, errors.New("user not found")
	}

	res, _ := bs.data.MyBook(userID)

	return res, nil
}
func (bs *bookSrv) Update(token interface{}, bookID int, updatedData book.Core) (book.Core, error) {

	if validasieror := bs.validasi.Struct(updatedData); validasieror != nil {
		return book.Core{}, nil
	}

	updatedDatas, err := bs.data.Update(bookID, updatedData)

	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "Book not found"
		} else {
			msg = "internal server error"
		}
		return book.Core{}, errors.New(msg)
	}

	return updatedDatas, nil
}

// All implements book.BookService
func (bs *bookSrv) AllBook() ([]book.Core, error) {
	All, err := bs.data.AllBook()
	fmt.Println("ini service", All)

	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "Book not found"
		} else {
			msg = "internal server error"
		}
		return []book.Core{}, errors.New(msg)
	}

	return All, nil
}
func (bs *bookSrv) Delete(token interface{}, bookID int) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		return errors.New("user not found")
	}

	err := bs.data.Delete(userID, bookID)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "Book not found"
		} else {
			msg = "internal server error"

		}
		return errors.New(msg)
	}
	return nil
}
