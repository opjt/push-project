package database

import (
	"fmt"
	"push/common/lib/env"
	"push/common/lib/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDB struct {
	conn *gorm.DB
}

func NewDB(env env.Env, log *logger.Logger) (*MariaDB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DB.User, env.DB.Password, env.DB.Host, env.DB.Port, env.DB.Database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.NewGormLogger(log),
	})
	if err != nil {
		return nil, fmt.Errorf("MariaDB 연결 실패: %w", err)
	}

	log.Debug("MariaDB 연결 성공")
	return &MariaDB{db}, nil
}
func (db *MariaDB) GetDB() *gorm.DB {
	return db.conn
}
func (db *MariaDB) Ping() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
