package repo

type IUserRepo interface {
	CreateUser() error
	UpdateUser() error
	DeleteUser() error
	GetUser() error
}

type userRepo struct{}

// DeleteUser implements IUserRepo.
func (u *userRepo) DeleteUser() error {
	panic("unimplemented")
}

// GetUser implements IUserRepo.
func (u *userRepo) GetUser() error {
	panic("unimplemented")
}

// UpdateUser implements IUserRepo.
func (u *userRepo) UpdateUser() error {
	panic("unimplemented")
}

// CreateUser implements IUserRepo.
func (u *userRepo) CreateUser() error {
	panic("unimplemented")
}

func NewUserRepo() IUserRepo {
	return &userRepo{}
}
