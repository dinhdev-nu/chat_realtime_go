package repo

type AuthRepo struct{}

func NewAuthRepo() *AuthRepo {
	return &AuthRepo{}
}

func (ar *AuthRepo) GetPing() string {
	return "pong"
}