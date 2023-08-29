package domain

type User struct {
	Email    string
	Password string
	// ConfirmPassword 不需要, 因为已经在Handler层处理过了
}
