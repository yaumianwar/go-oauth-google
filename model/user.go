package model

type User struct {
	ID            uint64      `json:"id" gorm:"column:id"`
	Sub           string `json:"sub" gorm:"column:sub"`
	Name          string `json:"name" gorm:"column:name"`
	GivenName     string `json:"given_name" gorm:"column:given_name"`
	FamilyName    string `json:"family_name" gorm:"column:family_name"`
	Profile       string `json:"profile" gorm:"column:profile"`
	Picture       string `json:"picture" gorm:"column:picture"`
	Email         string `json:"email" gorm:"column:email"`
	EmailVerified bool   `json:"email_verified" gorm:"column:email_verified"`
	Gender        string `json:"gender" gorm:"column:gender"`
}

