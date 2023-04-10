package transaction

import "time"

type (
	PartnerTransactionRequest struct {
		Id                     string `json:"id" `
		TransactionId          string `json:"transaction_id" validate:"required" `
		PartnerId              string `json:"partner_id" validate:"required" `
		CustomerId             string `json:"customer_id,omitempty" validate:"required" gorm:"-" `
		ProductName            string `json:"product_name" validate:"required" gorm:"-"`
		OtrPrice               int    `json:"otr_price" validate:"required" gorm:"-"`
		IdCardImage            string `json:"id_card_image" validate:"required"`
		SelfieImage            string `json:"selfie_image" validate:"required"`
		TransactionLetterImage string `json:"transaction_letter_image" validate:"required"`
		SalesId                string `json:"sales_id" validate:"required"`
		InstallmentPeriod      int    `json:"installment_period" validate:"required,gte=1,lte=4" gorm:"-"`
	}
	TransactionRequest struct {
		Id                string `json:"id" `
		CustomerId        string `json:"customer_id" `
		ProductId         string `json:"product_id" validate:"required"`
		ProductName       string `json:"product_name" `
		InstallmentPeriod int    `json:"installment_period" validate:"required" `
		InstallmentAmount int    `json:"installment_amount" `
		OtrPrice          int    `json:"partner_address" `
		Fee               int    `json:"fee" `
		Interest          int    `json:"Interest" `
		TransactionCome   string `json:"transaction_come"`
		Status            string `json:"status"`
	}
	TransactioQueryParam struct {
		Id        string `json:"id" `
		Limit     string `json:"limit"`
		Page      string `json:"page"`
		Status    string `json:"status"`
		Role      string `json:"name"`
		PartnerId string `json:"partner_id"`
	}
	TransactionResponse struct {
		Id                string        `json:"id" `
		CustomerId        string        `json:"customer_id" `
		ProductId         string        `json:"product_id"`
		ProductName       string        `json:"product_name"`
		InstallmentPeriod int           `json:"installment_period"`
		InstallmentAmount int           `json:"installment_amount" `
		OtrPrice          int           `json:"partner_address" `
		Fee               int           `json:"fee" `
		Interest          int           `json:"Interest" `
		TransactionCome   string        `json:"transaction_come"`
		Status            string        `json:"status"`
		Installment       []Installment `json:"Installment" gorm:"foreignKey:transaction_id;" `
		CreatedAt         time.Time     `json:"created_at"`
	}
	Installment struct {
		Id            string    `json:"id" `
		TransactionId string    `json:"transaction_id"  `
		Pinalty       int       `json:"pinalty"  `
		Period        time.Time `json:"period"  `
		Status        string    `json:"status"  `
	}

	Tabler interface {
		TableName() string
	}
)

func (TransactionRequest) TableName() string {
	return "transactions"
}

func (PartnerTransactionRequest) TableName() string {
	return "transactions_document"
}
func (TransactionResponse) TableName() string {
	return "transactions"
}

// type TransactionDocumentRepositoryRequest struct {
// 	Id                string `json:"id" `
// 	CustomerId        string `json:"customer_id" `
// 	ProductId         string `json:"product_id" validate:"required"`
// 	ProductName       string `json:"product_name" `
// 	InstallMentPeriod string `json:"installment_period" validate:"required" `
// 	InstallAmount     string `json:"installment_amount" `
// 	OtrPrice          int    `json:"partner_address" `
// 	Fee               int    `json:"fee" `
// 	Interest          int    `json:"Interest" validate:"max=2,min=1"`
// 	TransactionCome   string `json:"transaction_come"`
// }
