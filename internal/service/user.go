package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// SignUp 不能把Handler层的request传进来，因为 service 在 Handler 的下层, 不能用上层的东西
// user 不用指针原因: (1) 不需要判断空指针 (2)消耗不了多少性能
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	return svc.repo.Create(ctx, u)
}
