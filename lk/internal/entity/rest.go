package entity

import (
	"time"
)

type InfoResponse struct {
	Id          uint32           `json:"id"`
	Name        string           `json:"name"`
	Address     string           `json:"address"`
	Logo        string           `json:"logo_url"`
	Description []string         `json:"description"`
	Img         []string         `json:"img_urls"`
	Schedule    []ScheduleDTO    `json:"schedule"`
	Phone       string           `json:"phone"`
	Email       string           `json:"email"`
	SocialMedia []SocialMediaDTO `json:"media"`
}

type BaseInfoRequest struct {
	Id      uint32 `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Logo    string `json:"logo_url"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type DescripAndImgs struct {
	Description    []string
	DescripIndexes []uint8
	Img            []Img
}

type Img struct {
	Index    uint8
	Data     []byte
	Ext      string
	MimeType string
}

type BaseInfo struct {
	Id      uint32
	Name    string
	Address string
	Logo    string
	Phone   string
	Email   string
}

func (b *BaseInfoRequest) FromDTO() BaseInfo {
	return BaseInfo{
		Id:      b.Id,
		Name:    b.Name,
		Address: b.Address,
		Logo:    b.Logo,
		Phone:   b.Phone,
		Email:   b.Email,
	}
}

func (b *BaseInfo) ToDTO() BaseInfoRequest {
	return BaseInfoRequest{
		Id:      b.Id,
		Name:    b.Name,
		Address: b.Address,
		Logo:    b.Logo,
		Phone:   b.Phone,
		Email:   b.Email,
	}
}

type ScheduleDTO struct {
	Day   string `json:"day"`
	Open  string `json:"open_time"`
	Close string `json:"close_time"`
}

type SocialMediaDTO struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Rest struct {
	Id          uint32
	Name        string
	Address     string
	Logo        string
	Description []string
	Img         []string
	Schedule    []Schedule
	Phone       string
	Email       string
	SocialMedia []SocialMedia
}

func (r *Rest) ToDTO() InfoResponse {
	schedule := make([]ScheduleDTO, len(r.Schedule))
	if len(r.Schedule) != 0 {
		for i, s := range r.Schedule {
			schedule[i] = s.ToDTO()
		}
	}
	return InfoResponse{
		Id:          r.Id,
		Name:        r.Name,
		Address:     r.Address,
		Logo:        r.Logo,
		Description: r.Description,
		Img:         r.Img,
		Schedule:    schedule,
		Phone:       r.Phone,
		Email:       r.Email,
	}
}

type Schedule struct {
	Day   string
	Open  time.Time
	Close time.Time
}

func (s *Schedule) ToDTO() ScheduleDTO {
	return ScheduleDTO{
		Day:   s.Day,
		Open:  s.Open.Format("15:00"),
		Close: s.Close.Format("15:00"),
	}
}

type SocialMedia struct {
	Type string
	Url  string
}
