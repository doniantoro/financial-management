package auth

type MysqlUserRepository interface {
	Find(req *UserRequest) (*UserResponse, error)
	FindById(id string) (*UserResponse, error)
}
type UserUsecase interface {
	Get(req *UserRequest) (*UserResponse, error)
	Logout(id string) error
	Refresh(id string) (*UserResponse, error)
}
