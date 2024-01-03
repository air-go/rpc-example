package respository

type Test struct {
	ID      int64  `gorm:"column:id" json:"id"`
	GoodsID int64  `gorm:"column:goods_id" json:"goods_id"`
	Name    string `gorm:"column:name" json:"name"`
}

func (Test) TableName() string {
	return "test"
}
