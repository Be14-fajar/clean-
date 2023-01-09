package services

import (
	"api/config"
	"api/features/book"
	"api/mocks"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAdd(d *testing.T) {
	data := mocks.NewBookData(d)

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
		Judul:       "Naruto",
		TahunTerbit: 2009,
		Penulis:     "masashi",
		Pemilik:     sample.Name,
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userID"] = sample.ID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	useToken, _ := token.SignedString([]byte(config.JWT_KEY))
	token, _ = jwt.Parse(useToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_KEY), nil
	})
	data.On("Add", sample.ID, Input).Return(Respon, nil).Once()
	svc := New(data)

	res, err := svc.Add(token, Input)
	assert.Nil(d, err)
	assert.Equal(d, Respon.ID, res.ID)
	assert.Equal(d, Respon.Pemilik, res.Pemilik)
	data.AssertExpectations(d)

}
