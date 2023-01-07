package handler

import (
	"api/features/book"
	"api/helper"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type bookHandle struct {
	srv book.BookService
}

// All implements book.BookHandler

func New(bs book.BookService) book.BookHandler {
	return &bookHandle{
		srv: bs,
	}
}

func (bh *bookHandle) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := AddUpdateBookRequest{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "format inputan salah")
		}

		cnv := ToCore(input)

		res, err := bh.srv.Add(c.Get("user"), *cnv)
		if err != nil {
			log.Println("trouble :  ", err.Error())
			return c.JSON(helper.PrintErrorResponse(err.Error()))
		}

		book := ToResponse("add", res)

		return c.JSON(helper.PrintSuccessReponse(http.StatusCreated, "sukses menambahkan buku", book))
	}
}

func (bh *bookHandle) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := AddUpdateBookRequest{}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, "format inputan salah")
		}
		cnv := *ToCore(input)

		id, _ := strconv.Atoi(c.Param("id"))

		res, err := bh.srv.Update(c.Get("user"), id, cnv)
		if err != nil {
			log.Println("trouble :  ", err.Error())
			return c.JSON(helper.PrintErrorResponse(err.Error()))
		}

		book := ToResponse("update", res)

		return c.JSON(helper.PrintSuccessReponse(http.StatusCreated, "sukses mnegupdate buku", book))
	}
}

func (bh *bookHandle) MyBook() echo.HandlerFunc {
	return func(c echo.Context) error {

		res, _ := bh.srv.MyBook(c.Get("user"))

		listRes := ListBookCoreToBooksRespon(res)

		return c.JSON(http.StatusOK, listRes)
	}
}
func (bh *bookHandle) All() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, _ := bh.srv.All()

		listRes := ListBookCoreToBooksRespon(result)
		fmt.Println("ini handler", listRes)
		return c.JSON(http.StatusOK, listRes)
	}

}
