package repo

type IAuthRepo interface {
	GetPing() string
}

type authRepo struct{}

func NewAuthRepo() IAuthRepo {
	return &authRepo{}
}

func (ar *authRepo) GetPing() string {
	return "pong"
}
