package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.CTime = now
	u.UTime = now
	// 注意要传指针
	return dao.db.WithContext(ctx).Create(&u).Error
}

// User 对应数据库表
type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间
	CTime int64
	// 更新时间
	UTime int64
}
