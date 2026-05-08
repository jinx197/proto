package model

type Employee struct {
	Id   int64  `gorm:"primaryKey"`
	Name string `gorm:"size:50"`
	Dept string `gorm:"size:50"`
}

func (Employee) TableName() string {
	return "employee"
}
