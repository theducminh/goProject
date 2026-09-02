package domain

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	FirstName    string `gorm:"not null"`
	LastName     string `gorm:"not null"`
	IsActive     bool   `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Roles        []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRole struct {
	UserID    uint `gorm:"primaryKey"`
	RoleID    uint `gorm:"primaryKey"`
	CreatedAt time.Time
}

type Product struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"uniqueIndex;not null"`
	Name      string `gorm:"not null"`
	Price     uint   `gorm:"not null"`
	Stock     uint   `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Supplier struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Customer struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Order struct {
	ID         uint   `gorm:"primaryKey"`
	CustomerID uint   `gorm:"not null;index"`
	Status     string `gorm:"not null;index"`
	Total      uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Items      []OrderItem
}

type OrderItem struct {
	ID        uint `gorm:"primaryKey"`
	OrderID   uint `gorm:"not null;index"`
	ProductID uint `gorm:"not null;index"`
	Quantity  uint `gorm:"not null"`
	UnitPrice uint `gorm:"not null"`
	CreatedAt time.Time
}
