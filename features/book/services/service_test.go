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
		assert.Equal(t, uint(0), res.ID)
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
