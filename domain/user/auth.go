package auth

import (
	domain_customer "multi-finance/domain/customer"
	domain_partner "multi-finance/domain/partner"
)

type UserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	// gorm.Model
	ID           string                    `json:"id"`
	PartnerId    string                    `json:"partner_id"`
	FullName     string                    `json:"full_name"`
	Email        string                    `json:"email" `
	Role         string                    `json:"role"`
	AccesToken   string                    `json:"access_token"`
	RefreshToken string                    `json:"refresh_token"`
	Partner      domain_partner.Partner    `gorm:"foreignKey:PartnerId;references:Id" json:"partner,omitempty"`
	Customer     *domain_customer.Customer `gorm:"foreignKey:Id;" json:"customer,omitempty" `
}
type Tabler interface {
	TableName() string
}

func (UserResponse) TableName() string {
	return "users"
}
