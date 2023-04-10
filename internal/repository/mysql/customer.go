package repository

import (
	domain_customer "multi-finance/domain/customer"

	helper "multi-finance/helper"

	"gorm.io/gorm"
)

type MysqlCustomerRepository struct {
	db *gorm.DB
}

func NewMysqlCustomerRepository(db *gorm.DB) *MysqlCustomerRepository {
	return &MysqlCustomerRepository{db}
}

func (q MysqlCustomerRepository) Find(condition string) (*domain_customer.Customer, error) {
	var user domain_customer.Customer
	err := q.db.Debug().Where(condition).Order("created_at desc").Last(&user)
	if err.Error != nil {
		return nil, err.Error
	}
	return &user, nil
}
func (q MysqlCustomerRepository) Store(req *domain_customer.CustomerRequest) error {
	err := q.db.Create(&req)

	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (q MysqlCustomerRepository) UpdateStatus(req *domain_customer.CustomerRequestUpdate) error {
	err := q.db.Model(req).Updates(req)
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (q MysqlCustomerRepository) Index(pagination helper.Pagination) (*helper.Pagination, error) {
	var customer []*domain_customer.Customer
	q.db.Scopes(helper.Paginate(customer, &pagination, q.db)).Find(&customer)

	pagination.Rows = customer
	return &pagination, nil
}
