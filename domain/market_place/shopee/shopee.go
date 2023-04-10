package marketplace_shopee

type (
	ReponseFindData struct {
		OtrPrice    int           `json:"otr_price"`
		Description string        `json:"description"`
		ProductName string        `json:"product_name"`
		Installment []Installment `json:"installment"`
	}
	RequestOrder struct {
		ProductId string `json:"product_id"`
		OrderId   string `json:"order_id"`
	}
	Installment struct {
		Period   int `json:"period"`
		Amount   int `json:"amount"`
		Interest int `json:"interest"`
	}
)
