package model

type WebsiteAcmeAccount struct {
	BaseModel
	Email      string `gorm:"not null" json:"email"`
	URL        string `gorm:"not null" json:"url"`
	PrivateKey string `gorm:"not null" json:"-"`
	Type       string `gorm:"not null;default:letsencrypt" json:"type"`
	EabKid     string `json:"eabKid"`
	EabHmacKey string `json:"eabHmacKey"`
	KeyType    string `gorm:"not null;default:2048" json:"keyType"`
	UseProxy   bool   `gorm:"default:false" json:"useProxy"`
	CaDirURL   string `json:"caDirURL"`
}

func (w WebsiteAcmeAccount) TableName() string {
	return "website_acme_accounts"
}
