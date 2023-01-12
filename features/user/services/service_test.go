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
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	repo := mocks.NewUserData(t) // mock data

	t.Run("Berhasil login", func(t *testing.T) {
		// input dan respond untuk mock data
		inputEmail := "jerry@alterra.id"
		// res dari data akan mengembalik password yang sudah di hash
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Name: "jerry", Email: "jerry@alterra.id", HP: "08123456", Password: hashed}

		repo.On("Login", inputEmail).Return(resData, nil).Once() // simulasi method login pada layer data

		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Tidak ditemukan", func(t *testing.T) {
		inputEmail := "putra@alterra.id"
		repo.On("Login", inputEmail).Return(user.Core{}, errors.New("data not found")).Once()

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

		repo.On("Login", inputEmail).Return(user.Core{}, errors.New("terdapat masalah pada server")).Once()
		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "be1422")

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "terdapat masalah pada server")
		assert.Empty(t, token)
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("Login error password doesnt match", func(t *testing.T) {
		inputEmail := "jerry@alterra.id"
		hashed, _ := helper.GeneratePassword("be1422")
		resData := user.Core{ID: uint(1), Name: "helmi", Email: "jerry@alterra.id", HP: "08123456", Password: hashed}
		repo.On("Login", inputEmail).Return(resData, nil).Once()
		srv := New(repo)
		token, res, err := srv.Login(inputEmail, "asal")

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
	service := New(repo)

	// Case: user mengganti nama
	t.Run("Update successfully", func(t *testing.T) {
		input := user.Core{Name: "dendy", Email: "fajar@gmail.com", Alamat: "jakart", HP: "081222222"}
		hashed, _ := helper.GeneratePassword("be1422")

		resData := user.Core{ID: uint(1), Name: "dendy", Email: "fajar@gmail.com", Alamat: "jakart", HP: "081222222", Password: hashed}
		repo.On("Update", uint(1), input).Return(resData, nil).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := service.Update(token, input)

		assert.NoError(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	// Case: id user tidak valid atau tidak ditemukan
	t.Run("Update error invalid id", func(t *testing.T) {
		input := user.Core{Name: "dendy", Email: "fajar@gmail.com", Alamat: "jakart", HP: "081222222"}
		token := jwt.New(jwt.SigningMethodHS256)
		res, err := service.Update(token, input)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "invalid user id")
		assert.Empty(t, res.Name)
		assert.Equal(t, input.ID, res.ID)
		repo.AssertExpectations(t)
	})

	// Case: user mengganti nama tetapi id tidak ditemukan??
	t.Run("Update error data not found", func(t *testing.T) {
		input := user.Core{Name: "dendy"}

		resData := user.Core{}
		repo.On("Update", uint(1), input).Return(resData, errors.New("not found")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := service.Update(token, input)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data user not found")
		assert.Empty(t, res.Name)
		assert.NotEqual(t, input.Name, res.Name)
		repo.AssertExpectations(t)
	})

	// // Case: database tidak dapat mengelola permintaan update
	t.Run("Update error internal server", func(t *testing.T) {
		input := user.Core{Name: "dendy"}

		resData := user.Core{}
		repo.On("Update", uint(1), input).Return(resData, errors.New("internal server error")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := service.Update(token, input)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "internal server error")
		assert.Empty(t, res.Name)
		assert.NotEqual(t, input.Name, res.Name)
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
	t.Run("Deactive error user not found", func(t *testing.T) {
		repo.On("Deactive", uint(1)).Return(user.Core{}, errors.New("not found")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		srv := New(repo)
		res, err := srv.Deactive(token)

		assert.NotNil(t, err)
		assert.EqualError(t, err, "data tidak ditemukan")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})

	// Case: user melakukan deactive account tetapi terjadi masalah pada database
	t.Run("Deactive  server error", func(t *testing.T) {
		repo.On("Deactive", uint(1)).Return(user.Core{}, errors.New("internal server error")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		srv := New(repo)
		res, err := srv.Deactive(token)

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, uint(0), res.ID)
		repo.AssertExpectations(t)
	})
}
func TestRegister(t *testing.T) {
	repo := mocks.NewUserData(t)

	srv := New(repo)

	// Case: user melakukan pendaftaran akun baru
	t.Run("Register successfully", func(t *testing.T) {
		// Prgramming input and return repo

		type SampleUsers struct {
			ID       int
			Name     string
			Email    string
			Password string
			Alamat   string
			HP       string
		}
		sample := SampleUsers{
			ID:       1,
			Name:     "fajar1411",
			Email:    "frizky@gmail.com",
			Password: "12345",
			Alamat:   "jakart",
			HP:       "0132334343",
		}
		input := user.Core{
			Name:     sample.Name,
			Email:    sample.Alamat,
			Password: sample.Password,
			Alamat:   sample.Alamat,
			HP:       sample.HP,
		}

		// Program service

		hashed, _ := helper.GeneratePassword(input.Password)
		resData := user.Core{
			Name:     input.Name,
			Email:    input.Email,
			Password: hashed,
			Alamat:   sample.Alamat,
			HP:       sample.HP,
		}
		repo.On("Register", mock.Anything).Return(resData, nil).Once()
		data, err := srv.Register(input)

		assert.Nil(t, err)
		errCompare := bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(input.Password))
		assert.NoError(t, errCompare)
		assert.Equal(t, data.ID, resData.ID)
		repo.AssertExpectations(t)
	})
	t.Run("Register error data duplicate", func(t *testing.T) {
		type SampleUsers struct {
			ID       int
			Name     string
			Email    string
			Password string
			Alamat   string
			HP       string
		}
		sample := SampleUsers{
			ID:       1,
			Name:     "fajar1411",
			Email:    "frizky@gmail.com",
			Password: "12345",
			Alamat:   "jakart",
			HP:       "0132334343",
		}
		input := user.Core{
			Name:     sample.Name,
			Email:    sample.Alamat,
			Password: sample.Password,
			Alamat:   sample.Alamat,
			HP:       sample.HP,
		}

		// Programming input and return repo
		repo.On("Register", mock.Anything).Return(user.Core{}, errors.New("duplicated")).Once()

		// Program service
		data, err := srv.Register(input)

		// Test
		assert.NotNil(t, err)
		assert.EqualError(t, err, "data sudah terdaftar")
		assert.Empty(t, data)
		repo.AssertExpectations(t)
	})
	t.Run("Register error data duplicate", func(t *testing.T) {
		type SampleUsers struct {
			ID       int
			Name     string
			Email    string
			Password string
			Alamat   string
			HP       string
		}
		sample := SampleUsers{
			ID:       1,
			Name:     "fajar1411",
			Email:    "frizky@gmail.com",
			Password: "12345",
			Alamat:   "jakart",
			HP:       "0132334343",
		}
		input := user.Core{
			Name:     sample.Name,
			Email:    sample.Alamat,
			Password: sample.Password,
			Alamat:   sample.Alamat,
			HP:       sample.HP,
		}

		// Programming input and return repo
		repo.On("Register", mock.Anything).Return(user.Core{}, errors.New("internal server error")).Once()

		// Program service
		data, err := srv.Register(input)

		// Test
		assert.NotNil(t, err)
		assert.EqualError(t, err, "terdapat masalah pada server")
		assert.Empty(t, data)
		repo.AssertExpectations(t)
	})
}
