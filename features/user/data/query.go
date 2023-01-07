package data

import (
	"api/features/user"
	"errors"
	"log"

	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserData {
	return &userQuery{
		db: db,
	}
}

func (uq *userQuery) Login(email string) (user.Core, error) {
	res := User{}

	if err := uq.db.Where("email = ?", email).First(&res).Error; err != nil {
		log.Println("login query error", err.Error())
		return user.Core{}, errors.New("data not found")
	}

	return ToCore(res), nil
}
func (uq *userQuery) Register(newUser user.Core) (user.Core, error) {
	cnv := CoreToData(newUser)
	err := uq.db.Create(&cnv).Error
	if err != nil {
		return user.Core{}, err
	}

	newUser.ID = cnv.ID

	return newUser, nil
}
func (uq *userQuery) Profile(id uint) (user.Core, error) {
	res := User{}
	if err := uq.db.Where("id = ?", id).First(&res).Error; err != nil {
		log.Println("Get By ID query error", err.Error())
		return user.Core{}, err
	}

	return ToCore(res), nil
}
func (uq *userQuery) Update(id uint, updateData user.Core) (user.Core, error) {
	userModel := CoreToData(updateData)
	userModel.ID = id

	Input := uq.db.Where("id = ?", id).Updates(&userModel)
	if Input.Error != nil {
		log.Println("Get By ID query error", Input.Error.Error())
		return user.Core{}, Input.Error
	}
	if Input.RowsAffected <= 0 {
		return user.Core{}, errors.New("Not found")
	}

	return ToCore(userModel), nil
}

// Deactive implements user.UserData
func (uq *userQuery) Deactive(id uint) (user.Core, error) {
	users := User{}

	delete := uq.db.Delete(&users, id)

	if delete.Error != nil {
		log.Println("Get By ID query error", delete.Error.Error())
		return user.Core{}, delete.Error
	}

	return ToCore(users), nil
}
