package data

import (
	"fmt"
	"time"

	"github.com/adamelfsborg-code/food/user/lib"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"

	"github.com/google/uuid"
)

//lint:ignore U1000 Ignore unused function temporarily for debugging
type UserDto struct {
	tableName struct{}  `pg:"app_user,alias:au"`
	Id        uuid.UUID `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Name      string    `json:"name" db:"name" validate:"max=20,min=3"`
	Password  string    `json:"password" db:"password" validate:"max=50,min=3"`
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
type UserTable struct {
	tableName struct{}  `pg:"app_user,alias:au"`
	Id        uuid.UUID `json:"id" db:"id"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Name      string    `json:"name" db:"name" validate:"max=20,min=3"`
}

func NewUserDto(name, password string) (*UserDto, error) {
	validate := validator.New()

	user := &UserDto{
		Name:     name,
		Password: password,
	}

	errs := validate.Struct(user)
	if errs != nil {
		return nil, errs
	}

	hash, err := lib.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Password = hash
	return user, nil
}

func (d *DataConn) Register(user UserDto) error {
	_, err := d.DB.Model(&user).Insert()
	return err
}

func (d *DataConn) Login(name, password string) (string, error) {
	user := d.getUserWithPasswordByName(name)

	valid := lib.CheckPasswordHash(password, user.Password)
	if !valid {
		return "", fmt.Errorf("user does not exist")
	}

	claims := jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(d.Env.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (d *DataConn) Ping(userId uuid.UUID) (UserTable, error) {
	var user UserTable

	err := d.DB.Model(&user).Where("id = ?", userId).Select()
	if err != nil {
		return user, err
	}

	return user, nil
}

func (d *DataConn) getUserWithPasswordByName(name string) UserDto {
	var user UserDto

	d.DB.Model(&user).Where("name = ?", name).Select()
	return user
}
