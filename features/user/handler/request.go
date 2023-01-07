package handler

import "api/features/user"

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterRequest struct {
	Name     string `json:"nama" form:"nama"`
	Email    string `json:"email" form:"email"`
	Alamat   string `json:"alamat" form:"alamat"`
	HP       string `json:"hp" form:"hp"`
	Password string `json:"password" form:"password"`
}

type UpdateRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Alamat   string `json:"alamat" form:"alamat"`
	HP       string `json:"hp" form:"hp"`
	Password string `json:"password" form:"password"`
}

func ToCore(data interface{}) *user.Core {
	res := user.Core{}

	switch data.(type) {
	case LoginRequest:
		cnv := data.(LoginRequest)
		res.Email = cnv.Email
		res.Password = cnv.Password
	case RegisterRequest:
		cnv := data.(RegisterRequest)
		res.Email = cnv.Email
		res.Name = cnv.Name
		res.Alamat = cnv.Alamat
		res.HP = cnv.HP
		res.Password = cnv.Password
	case UpdateRequest:
		cnv := data.(UpdateRequest)
		res.Email = cnv.Email
		res.Name = cnv.Name
		res.Alamat = cnv.Alamat
		res.HP = cnv.HP
		res.Password = cnv.Password
	default:
		return nil
	}

	return &res
}
