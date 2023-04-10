package customer

import "time"

type (
	CustomerRequest struct {
		Id           string `json:"id"`
		IdNumber     string `json:"id_number" validate:"required,max=16,min=16"`
		UserId       string `json:"user_id" `
		Msisdn       string `json:"msisdn" validate:"required"`
		FullName     string `json:"full_name" validate:"required"`
		LegalName    string `json:"legal_name" validate:"required"`
		PlaceOfBirth string `json:"place_of_birth" validate:"required"`
		DateOfBirth  string `json:"date_of_birth" validate:"required"`
		Sallary      int    `json:"sallary" validate:"required"`
		IdCardImage  string `json:"id_card_image" validate:"required"`
		SelfieImage  string `json:"selfie_image" validate:"required"`
	}
	CustomerRequestUpdate struct {
		Id             string    `json:"id" validate:"required"`
		Limit          int       `json:"limit" `
		AcceptedUserId string    `json:"accepted_user_id" `
		RemainLimit    int       `json:"remain_limit" `
		Status         string    `json:"status" validate:"required"`
		UpdatedAt      time.Time `json:"updated_at"`
	}

	CustomerQueryParam struct {
		Id       string `json:"id" `
		Limit    string `json:"limit"`
		Page     string `json:"page"`
		Status   string `json:"status"`
		Name     string `json:"name"`
		IdNumber string `json:"id_number"`
	}
	Customer struct {
		Id             string    `json:"id,omitempty" validate:"required"`
		IdNumber       string    `json:"id_number,omitempty" validate:"required"`
		AcceptedUserId string    `json:"accepted_user_id,omitempty" `
		UserId         string    `json:"user_id,omitempty"`
		Msisdn         string    `json:"msisdn,omitempty" validate:"required"`
		FullName       string    `json:"full_name,omitempty" validate:"required"`
		LegalName      string    `json:"legal_name,omitempty" validate:"required"`
		PlaceOfBirth   string    `json:"place_of_birth,omitempty" validate:"required"`
		DateOfBirth    string    `json:"date_of_birth,omitempty" validate:"required"`
		Sallary        int       `json:"sallary,omitempty" validate:"required"`
		IdCardImage    string    `json:"id_card_image,omitempty" validate:"required"`
		SelfieImage    string    `json:"selfie_image,omitempty" validate:"required"`
		Limit          int       `json:"limit,omitempty" validate:"required"`
		RemainLimit    int       `json:"remain_limit,omitempty" validate:"required"`
		Status         string    `json:"status,omitempty" validate:"required"`
		CreatedAt      time.Time `json:"created_at,omitempty"`
		UpdatedAt      time.Time `json:"updated_at,omitempty"`
	}
)

func (CustomerRequest) TableName() string {
	return "customers"
}

// TableName overrides the table name used by User to `profiles`
func (CustomerRequestUpdate) TableName() string {
	return "customers"
}

// type Tabler interface {
// 	TableName() string
// }
