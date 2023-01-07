package data

import (
	"api/features/book"

	"gorm.io/gorm"
)

type Books struct {
	gorm.Model
	Judul       string
	TahunTerbit int
	Penulis     string
	UserID      uint
}

type BookPemilik struct {
	ID          uint
	Judul       string
	TahunTerbit int
	Penulis     string
	Name        string
}

func ToCore(data Books) book.Core {
	return book.Core{
		ID:          data.ID,
		Judul:       data.Judul,
		TahunTerbit: data.TahunTerbit,
		Penulis:     data.Penulis,
	}
}

func (dataModel *BookPemilik) ModelsToCore() book.Core { //fungsi yang mengambil data dari  user gorm(model.go)  dan merubah data ke entities usercore
	return book.Core{
		ID:          dataModel.ID,
		Judul:       dataModel.Judul,
		Penulis:     dataModel.Penulis,
		TahunTerbit: dataModel.TahunTerbit,
		Pemilik:     dataModel.Name,
	}
}

func ListModelTOCore(dataModel []BookPemilik) []book.Core { //fungsi yang mengambil data dari  user gorm(model.go)  dan merubah data ke entities usercore
	var dataCore []book.Core
	for _, value := range dataModel {
		dataCore = append(dataCore, value.ModelsToCore())
	}
	return dataCore //  untuk menampilkan data ke controller
}

func CoreToData(data book.Core) Books {
	return Books{
		Model:       gorm.Model{ID: data.ID},
		Judul:       data.Judul,
		Penulis:     data.Penulis,
		TahunTerbit: data.TahunTerbit,
	}
}
