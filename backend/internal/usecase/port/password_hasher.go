package port

type IPasswordHasher interface {
	CheckPasswordHash(password, hash string) (bool, error)
	HashPassword(password string) (string, error)
}
