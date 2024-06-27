package entity

type CreateShopRequest struct {
	UserId string `query:"user_id" validate:"required,uuid"`
	Name   string `json:"name" validate:"required,min=3,max=100"`
}

type UpsertShopResponse struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
