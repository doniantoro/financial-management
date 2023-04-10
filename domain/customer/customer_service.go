package customer

import "multi-finance/helper"

type CustomerUsecase interface {
	ListApplication(req *CustomerQueryParam) (*helper.Pagination, error)
	ApplyApplication(req *CustomerRequest, id string) error
	FindApplication(id string) (*Customer, error)
	UpdateStatusApplication(req *CustomerRequestUpdate) error
}
type CustomerRepository interface {
	Index(req helper.Pagination) (*helper.Pagination, error)
	Store(req *CustomerRequest) error
	Find(condition string) (*Customer, error)
	UpdateStatus(req *CustomerRequestUpdate) error
}
