package domain

type User struct {
	Id       int64
	Email    string
	Password string
	// ConfirmPassword 不需要, 因为已经在Handler层处理过了
}
