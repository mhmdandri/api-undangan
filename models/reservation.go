package models

import "time"

type Reservation struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name" gorm:"size:100;not null"`
	IsPresent   bool       `json:"is_present" gorm:"not null;default:false"`
	Email       string     `json:"email" gorm:"size:100;default:''"`
	Phone       string     `json:"phone" gorm:"size:30;default:''"`
	Code        string     `json:"code" gorm:"size:50;not null;unique"`
	QuotaGuests int        `json:"quota_guests" gorm:"default:2"`
	TotalGuests int        `json:"total_guests" gorm:"default:1"`
	Status      string     `json:"status" gorm:"size:50;not null;default:'belum_konfirmasi'"`
	Source      string     `json:"source" gorm:"size:50;not null;default:'self_registered'"`
	AttendedAt  *time.Time `json:"attended_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
