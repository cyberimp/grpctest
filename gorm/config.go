package gorm

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (s *Conn) ConnectDB() error {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=5432"
	var err error
	s.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return s.migrate()
}

func (s *Conn) migrate() error {
	return s.db.AutoMigrate(&Post{})
}
