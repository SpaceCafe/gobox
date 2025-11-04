package authentication

// Repository defines a repository for retrieving principals by token or credentials.
type Repository interface {
	GetByToken(token string) (Principal, error)
	GetByCredentials(username, password string) (Principal, error)
}
