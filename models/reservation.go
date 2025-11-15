package models

import "time"
type Reservation struct {
	ID					 uint   `gorm:"primaryKey" json:"id"`
	Name				 string `json:"name" gorm:"size:100;not null"`
	IsPresent		 bool   `json:"is_present" gorm:"not null"`
	Email			 string `json:"email" gorm:"size:100;not null" binding:"email,required"`
	Code 			 string `json:"code" gorm:"size:50;not null;unique"`
	TotalGuests			 int    `json:"total_guests" gorm:"size:3:not null;default:1; max:3"`
	Status			 string `json:"status" gorm:"size:50;not null;default:'tidak_datang'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}