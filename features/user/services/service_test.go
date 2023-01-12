package services

import (
	"api/features/user"
	"api/helper"
	"api/mocks"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	repo := mocks.NewUserData(t) // mock data

	t.Run("Berhasil login", func(t *testing.T) {
		// input dan respond untuk mock data
		inputEmail := "jerry@alterra.id"
		// res dari data akan mengembalik password yang sudah di hash
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Name: "jerry", Email: "jerry@alterra.id", HP: "08123456", Password: hashed}

		repo.On("Login", inputEmail).Return(resData, nil) // simulasi method login pada layer data

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Tidak ditemukan", func(t *testing.T) {
		inputEmail := "putra@alterra.id"
		repo.On("Login", inputEmail).Return(user.Core{}, errors.New("data not found"))

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
	t.Run("server error", func(t *testing.T) {
		inputEmail := "jerry@alterra.id"

		repo.On("Login", inputEmail).Return(user.Core{}, errors.New("not found")).Once()
		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terdapat masalah pada server")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
	t.Run("Salah password", func(t *testing.T) {
		inputEmail := "jerry@alterra.id"
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Name: "jerry", Email: "jerry@alterra.id", HP: "08123456", Password: hashed}
		repo.On("Login", inputEmail).Return(resData, nil)

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1423")
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password tidak sesuai")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

}

func TestProfile(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Sukses lihat profile", func(t *testing.T) {
		resData := user.Core{ID: uint(1), Name: "jerry", Email: "jerry@alterra.id", HP: "08123456"}

		repo.On("Profile", uint(1)).Return(resData, nil).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Profile(pToken)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Profile(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data tidak ditemukan", func(t *testing.T) {
		repo.On("Profile", uint(4)).Return(user.Core{}, errors.New("data not found")).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(4)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("masalah di server", func(t *testing.T) {
		repo.On("Profile", mock.Anything).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Profile(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewUserData(t)
	srv := New(repo)

	// Case: user mengganti nama
	t.Run("Update success", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}
		Input := user.Core{
			ID:     1,
			Name:   "fajar1411",
			Email:  "frizky861@gmail.com",
			Alamat: "jakarta",
			HP:     "08122323232",
		}
		Respon := user.Core{
			ID:     Input.ID,
			Name:   Input.Name,
			Email:  Input.Email,
			Alamat: Input.Email,
			HP:     Input.HP,
		}

		repo.On("Update", uint(sample.ID), Input).Return(Respon, nil).Once()

		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, Input)

		assert.NoError(t, err)
		assert.Equal(t, Respon.ID, res.ID)
		repo.AssertExpectations(t)
	})
	t.Run("jwt tidak valid", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   4,
			Name: "fajar1411",
		}
		Input := user.Core{
			ID:     1,
			Name:   "fajar1411",
			Email:  "frizky861@gmail.com",
			Alamat: "jakarta",
			HP:     "08122323232",
		}

		srv := New(repo)

		_, token := helper.GenerateJWT(sample.ID)

		res, err := srv.Update(token, Input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "invalid user id")
		assert.Equal(t, uint(0), res.ID) //perbandingan
	})

	t.Run("Update error internal server", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}

		Input := user.Core{

			Name: "fajar1411",
		}
		repo.On("Update", uint(sample.ID), Input).Return(user.Core{}, errors.New("internal server error")).Once() ///data yang akan di testinng//once data yang di pakai saat add buku
		srv := New(repo)                                                                                          //new service

		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, Input)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Empty(t, res.Name)
		assert.NotEqual(t, Input.Name, res.Name)
		repo.AssertExpectations(t)
	})

}
func TestDeactive(t *testing.T) {
	repo := mocks.NewUserData(t)

	t.Run("Deactive succesfully", func(t *testing.T) {

		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   1,
			Name: "fajar1411",
		}

		Respon := user.Core{
			ID:     0,
			Name:   "",
			Email:  "",
			Alamat: "",
			HP:     "",
		}

		repo.On("Deactive", uint(sample.ID)).Return(Respon, nil).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(sample.ID)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Deactive(pToken)

		assert.NoError(t, err)
		assert.Equal(t, Respon.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("jwt tidak valid", func(t *testing.T) {
		type SampleUsers struct {
			ID   int
			Name string
		}
		sample := SampleUsers{
			ID:   4,
			Name: "fajar1411",
		}

		srv := New(repo)

		_, token := helper.GenerateJWT(sample.ID)

		res, err := srv.Deactive(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "id user not found")
		assert.Equal(t, uint(0), res.ID) //perbandingan
	})

}
