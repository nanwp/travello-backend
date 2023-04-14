package categories

import "time"

type Category struct {
	ID int `gorm:"AUTO_INCREMENT"`
	Name string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Category) TableName() string{
	return "tbl_category"
}
