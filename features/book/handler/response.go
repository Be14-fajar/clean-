package handler

import "api/features/book"

type BookResponse struct {
	ID          uint   `json:"id"`
	Judul       string `json:"judul"`
	TahunTerbit int    `json:"tahun_terbit"`
	Penulis     string `json:"penulis"`
	Pemilik     string `json:"pemilik"`
}
type AddBookResponse struct {
	Judul       string `json:"judul"`
	TahunTerbit int    `json:"tahun_terbit"`
	Penulis     string `json:"penulis"`
}
type updateBookResponse struct {
	Judul       string `json:"judul"`
	TahunTerbit int    `json:"tahun_terbit"`
	Penulis     string `json:"penulis"`
}

func ToResponse(feature string, book book.Core) interface{} {
	switch feature {
	case "add":
		return AddBookResponse{
			Judul:       book.Judul,
			TahunTerbit: book.TahunTerbit,
			Penulis:     book.Penulis,
		}
	case "update":
		return updateBookResponse{
			Judul:       book.Judul,
			TahunTerbit: book.TahunTerbit,
			Penulis:     book.Penulis,
		}
	default:
		return BookResponse{
			ID:          book.ID,
			Judul:       book.Judul,
			TahunTerbit: book.TahunTerbit,
			Penulis:     book.Penulis,
			Pemilik:     book.Pemilik,
		}
	}
}

func ListBookCoreToBookRespon(dataCore book.Core) BookResponse { // data user core yang ada di controller yang memanggil user repository
	return BookResponse{
		ID:          dataCore.ID,
		Judul:       dataCore.Judul,
		TahunTerbit: dataCore.TahunTerbit,
		Penulis:     dataCore.Penulis,
		Pemilik:     dataCore.Pemilik,
	}
}
func ListBookCoreToBooksRespon(dataCore []book.Core) []BookResponse { //data user.core data yang diambil dari entities ke respon struct
	var ResponData []BookResponse

	for _, value := range dataCore { //memanggil paramete data core yang berisi data user core
		ResponData = append(ResponData, ListBookCoreToBookRespon(value)) // mengambil data mapping dari user core to respon
	}
	return ResponData
}
