package entity

import "time"

type CreateShopRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`
	Name   string `json:"name" validate:"required,min=3,max=100"`
}

type UpdateShopRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`
	Id     string `params:"id" validate:"required,uuid"`
	Name   string `json:"name" validate:"required,min=3,max=100"`
}

type UpsertShopResponse struct {
	Id        string `json:"id" db:"id"`
	UserId    string `json:"user_id" db:"user_id"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

type DeleteShopRequest struct {
	Id     string `query:"id" validate:"required,uuid"`
	UserId string `query:"user_id" validate:"required,uuid"`
}

type GetShopsRequest struct {
	// UserId string `query:"user_id" validate:"required,uuid"`

	Page     int    `query:"page" validate:"required"`
	Limit    int    `query:"limit" validate:"required"`
	ShopName string `query:"shop_name"`
}

func (g *GetShopsRequest) SetDefaults() {
	if g.Page < 1 {
		g.Page = 1
	}

	if g.Limit < 1 {
		g.Limit = 10
	}
}

type GetShopsResponse struct {
	Items []ShopItem `json:"items"`
	Meta  Meta       `json:"meta"`
}

type ShopItem struct {
	Id        string     `json:"id" db:"id"`
	UserId    string     `json:"user_id" db:"user_id"`
	Name      string     `json:"name" db:"name"`
	CretedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

type Meta struct {
	TotalData int `json:"total_data"`
	TotalPage int `json:"total_page"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}

func (m *Meta) CountTotalPage() {
	if m.TotalData == 0 {
		m.TotalPage = 0
		return
	}

	m.TotalPage = m.TotalData / m.Limit
	if m.TotalData%m.Limit > 0 {
		m.TotalPage++
	}
}
