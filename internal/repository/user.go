package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

// UserRepository 在 repository 层没有注册的概念，所以叫 create
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	// 有缓存就在这里操作
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
