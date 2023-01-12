package services

import (
	"api/features/book"
	"api/helper"
	"api/mocks"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	data := mocks.NewBookData(t)

	t.Run("Berhasil Add Buku", func(t *testing.T) {

		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}
		Input := book.Core{
			Judul:       "Naruto",
			TahunTerbit: 2009,
			Penulis:     "masashi",
			Pemilik:     sample.Name,
		}

		Respon := book.Core{
			ID:          1,
			Judul:       Input.Judul,
			TahunTerbit: Input.TahunTerbit,
			Penulis:     Input.Penulis,
			Pemilik:     Input.Pemilik,
		}

		_, token := helper.GenerateJWT(sample.ID)
		useToken := token.(*jwt.Token)
		useToken.Valid = true

		data.On("Add", sample.ID, Input).Return(Respon, nil).Once()
		svc := New(data)

		res, err := svc.Add(useToken, Input)
		assert.Nil(t, err)
		assert.Equal(t, Respon.ID, res.ID)
		assert.Equal(t, Respon.Pemilik, res.Pemilik)
		data.AssertExpectations(t)
	})
	t.Run("jwt tidak valid", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}
		Input := book.Core{
			Judul:       "Naruto",
			TahunTerbit: 2009,
			Penulis:     "masashi",
			Pemilik:     sample.Name,
		}

		srv := New(data)

		_, token := helper.GenerateJWT(sample.ID)

		res, err := srv.Add(token, Input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user not found")
		assert.Equal(t, uint(0), res.ID) //perbandingan
	})
	t.Run("validation error", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}
		Input := book.Core{
			Judul:       "",
			TahunTerbit: 2009,
			Penulis:     "",
			Pemilik:     sample.Name,
		}

		//yang di coba testing

		srv := New(data)

		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, Input)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "validation error")
		assert.Equal(t, uint(0), res.ID)
	})
	t.Run("data not found", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   4,
			Name: "fajar1411",
		}

		Input := book.Core{
			Judul:       "Naruto",
			TahunTerbit: 2009,
			Penulis:     "masashi",
			Pemilik:     sample.Name,
		}
		data.On("Add", sample.ID, Input).Return(book.Core{}, errors.New("data not found")).Once() ///data yang akan di testinng//once data yang di pakai saat add buku

		srv := New(data)

		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, Input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "Book not found")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}

		Input := book.Core{
			Judul:       "Naruto",
			TahunTerbit: 2009,
			Penulis:     "masashi",
			Pemilik:     sample.Name,
		}
		data.On("Add", sample.ID, Input).Return(book.Core{}, errors.New("internal server error")).Once() ///data yang akan di testinng//once data yang di pakai saat add buku
		srv := New(data)                                                                                 //new service

		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, Input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		data.AssertExpectations(t)
	})
}

func TestAllBook(t *testing.T) {
	data := mocks.NewBookData(t)
	svc := New(data)
	t.Run("Berhasil Melihat semua Buku", func(t *testing.T) {

		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}
		Respon := []book.Core{
			{
				ID:          1,
				Judul:       "Naruto",
				Penulis:     "Masashi Kishimoto",
				TahunTerbit: 2000,
				Pemilik:     sample.Name,
			},
			{
				ID:          2,
				Judul:       "Boruto",
				Penulis:     "Masashi Kishimoto",
				TahunTerbit: 20020,
				Pemilik:     sample.Name,
			},
			{
				ID:          3,
				Judul:       "One piece",
				Penulis:     "Oda sensei",
				TahunTerbit: 1999,
				Pemilik:     sample.Name,
			},
		}
		data.On("AllBook").Return(Respon, nil).Once()
		svc := New(data)
		actual, err := svc.AllBook()
		assert.Nil(t, err)
		assert.Equal(t, Respon[0].ID, actual[0].ID)
		assert.Equal(t, Respon[0].Judul, actual[0].Judul)
		assert.Equal(t, Respon[0].Pemilik, actual[0].Pemilik)
		assert.Equal(t, Respon[1].ID, actual[1].ID)
		assert.Equal(t, Respon[1].Pemilik, actual[1].Pemilik)
		assert.Equal(t, Respon[2].ID, actual[2].ID)
		assert.Equal(t, Respon[2].Pemilik, actual[2].Pemilik)
	})

	// Case: user ingin melihat list buku yang, tetapi buku tidak ada buku yang ditemukan
	t.Run(" all book not found", func(t *testing.T) {
		// Programming input and return repo
		data.On("AllBook").Return(nil, errors.New("Book not found")).Once()

		// Program service
		actual, err := svc.AllBook()

		// Test
		assert.NotNil(t, err)
		assert.EqualError(t, err, "Book not found")
		assert.Nil(t, actual)

	})
	t.Run("Get all book error server", func(t *testing.T) {
		// Programming input and return repo
		data.On("AllBook").Return([]book.Core{}, errors.New("internal server error")).Once()

		// Program service
		actual, err := svc.AllBook()

		// Test
		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal server error")
		assert.Nil(t, actual)

	})
}
func TestUpdateBook(t *testing.T) {
	input := book.Core{Judul: "One Piece"}
	resData := book.Core{
		ID:      1,
		Judul:   "Naruto",
		Penulis: "Masashi Kishimoto",
		Pemilik: "fajar",
	}

	repo := mocks.NewBookData(t)

	srv := New(repo)

	t.Run("Update successfully", func(t *testing.T) {
		repo.On("Update", 1, 1, input).Return(resData, nil).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		actual, err := srv.Update(token, 1, input)

		assert.Nil(t, err)
		assert.Equal(t, resData.Judul, actual.Judul)
		assert.Equal(t, resData.ID, actual.ID)
		assert.Equal(t, resData.Pemilik, actual.Pemilik)

		repo.AssertExpectations(t)
	})

	t.Run("Update error user not found", func(t *testing.T) {

		token := jwt.New(jwt.SigningMethodHS256)
		actual, err := srv.Update(token, 1, input)

		// Test
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "id user not found")
		assert.Empty(t, actual)
	})

	t.Run("Update error invalid", func(t *testing.T) {
		input := book.Core{
			Judul:       "nar",
			Penulis:     "mas",
			TahunTerbit: 2000,
		}

		// Program service
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		actual, err := srv.Update(token, 1, input)

		// Test
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "input update book invalid")
		assert.Empty(t, actual)
	})

	t.Run("Update error book not found", func(t *testing.T) {
		// Programming input and return repo
		repo.On("Update", 1, 1, input).Return(book.Core{}, errors.New("not found")).Once()

		// Program service
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		actual, err := srv.Update(token, 1, input)

		// Test
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "Book not found")
		assert.Empty(t, actual)
		repo.AssertExpectations(t)
	})

	t.Run("Update error internal server", func(t *testing.T) {
		// Programming input and return repo
		repo.On("Update", 1, 1, input).Return(book.Core{}, errors.New("internal server error")).Once()

		// Program service
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		actual, err := srv.Update(token, 1, input)

		// Test
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Empty(t, actual)
		repo.AssertExpectations(t)
	})
}
func TestDeleteBook(t *testing.T) {
	repo := mocks.NewBookData(t)

	srv := New(repo)
	t.Run("Delete Success", func(t *testing.T) {
		repo.On("Delete", 1, 1).Return(nil).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(token, 1)

		assert.Nil(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("Delete Error", func(t *testing.T) {
		repo.On("Delete", 1, 1).Return(errors.New("user id not found")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(token, 1)

		assert.NotNil(t, err)

		repo.AssertExpectations(t)
	})
	t.Run("Delete Error", func(t *testing.T) {
		repo.On("Delete", 1, 1).Return(errors.New("Book not found")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(token, 1)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "Book not found")
		repo.AssertExpectations(t)
	})
	t.Run("Delete server error", func(t *testing.T) {
		repo.On("Delete", 1, 1).Return(errors.New("internal server error")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(token, 1)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})
}

func TestMyBook(t *testing.T) {
	repo := mocks.NewBookData(t)

	srv := New(repo)

	// Case: user ingin melihat list buku yang dimilikinya
	t.Run("MyBook list succesfully", func(t *testing.T) {
		resData := []book.Core{
			{
				ID:          1,
				Judul:       "Naruto",
				Penulis:     "Masashi Kishimoto",
				TahunTerbit: 1999,
			},
			{
				ID:          2,
				Judul:       "Dragon ball",
				Penulis:     "Akira toriyama",
				TahunTerbit: 1998,
			},
		}

		// Programming input and return repo
		repo.On("MyBook", 1).Return(resData, nil).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		actual, err := srv.MyBook(token)

		// Test
		assert.Nil(t, err)
		assert.Equal(t, resData[0].ID, actual[0].ID)
		assert.Equal(t, resData[0].Judul, actual[0].Judul)
		assert.Equal(t, resData[1].ID, actual[1].ID)
		assert.Equal(t, resData[1].Judul, actual[1].Judul)
	})
}
