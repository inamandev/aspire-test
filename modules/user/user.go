package user

import "errors"

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Status   bool   `json:"Status"`
}

func GetByUsername(username string) (User, error) {
	if username == "admin" {
		return User{
			Id:       1,
			Username: "admin",
			Password: "$2a$10$W90qj9W7FXEBiT/8dBjRFel3o98ZEBnLrxu/qD3kGrRN0X1sYpRaC",
			Name:     "Admin",
			Role:     "admin",
			Status:   true,
		}, nil
	}
	return User{}, errors.New("not admin")
}

func GetHashedPassword(username string) string {
	return "$2a$10$W90qj9W7FXEBiT/8dBjRFel3o98ZEBnLrxu/qD3kGrRN0X1sYpRaC"
}
