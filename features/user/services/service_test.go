package services

import (
	"api/features/user"
	"api/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegis(t *testing.T) {
	data := mocks.NewUserData(t)

	Input := user.Core{
		Name:     "fajar",
		Email:    "frizky861@gmail.com",
		Alamat:   "Pasuruan",
		HP:       "081223234",
		Password: "fajar1411",
	}

	Respon := user.Core{
		ID:       uint(1),
		Name:     "fajar",
		Email:    "frizky861@gmail.com",
		Alamat:   "Pasuruan",
		HP:       "081223234",
		Password: "$2a$10$AnAc63Ppl4nCdckcMV2lR.EsLezCtZujjzTLoq/XAoA/NivbDqzCO",
	}
	data.On("Register", mock.Anything).Return(Respon, nil).Once()
	svc := New(data)

	res, err := svc.Register(Input)
	assert.Nil(t, err)
	assert.Equal(t, Respon.ID, res.ID)
	assert.Equal(t, Respon.Name, res.Name)
	data.AssertExpectations(t)

}
