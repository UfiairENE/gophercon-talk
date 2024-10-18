package db

type Data struct {
	ID      uint `gorm:"primaryKey"`
	Column1 string
	Column2 int  
	//`gorm:"index"`
}
