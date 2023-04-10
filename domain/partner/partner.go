package partner

type Partner struct {
	Id             string `json:"id,omitempty"`
	PartnerName    string `json:"partner_name,omitempty" `
	PartnerAddress string `json:"partner_address,omitempty" `
	CreatedAt      string `json:"created_at,omitempty" `
}
