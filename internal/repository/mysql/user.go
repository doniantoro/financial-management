package repository

import (
	domain_user "multi-finance/domain/user"

	"gorm.io/gorm"
)

type MysqlUserRepository struct {
	db *gorm.DB
}

func NewMysqlUserRepository(db *gorm.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db}
}

func (q MysqlUserRepository) Find(req *domain_user.UserRequest) (*domain_user.UserResponse, error) {
	var user domain_user.UserResponse
	err := q.db.Where("email = ? AND password = ? ", req.Email, req.Password).Model(&user).Preload("Partner").First(&user)

	if err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
func (q MysqlUserRepository) FindById(id string) (*domain_user.UserResponse, error) {
	var user domain_user.UserResponse
	err := q.db.Where("id = ?  ", id).First(&user)
	// var user = User{ID: 10}
	if err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
