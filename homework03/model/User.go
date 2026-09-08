package model

type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string `gorm:"uniqueIndex"`
	Age   int
	Posts []Post `gorm:"constraint:OnDelete:CASCADE;"`
}
