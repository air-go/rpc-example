package respository

type Test struct {
	ID      uint   `gorm:"column:id" json:"id"`
	GoodsID uint64 `gorm:"column:goods_id" json:"goods_id"`
	Name    string `gorm:"column:name" json:"name"`
}

func (Test) TableName() string {
	return "test"
}
