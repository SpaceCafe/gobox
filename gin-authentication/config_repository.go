package authentication

var _ Repository = (*ConfigRepository)(nil)

// ConfigRepository implements the Repository interface using a Config object.
type ConfigRepository struct {
	cfg *Config
}

// NewConfigRepository creates a new ConfigRepository with the given configuration.
func NewConfigRepository(cfg *Config) *ConfigRepository {
	return &ConfigRepository{cfg: cfg}
}

// GetByCredentials retrieves a principal by credentials.
//
//nolint:ireturn // Principal is implemented by the repository.
func (r *ConfigRepository) GetByCredentials(id, password string) (Principal, error) {
	if passwd, ok := r.cfg.Principals[id]; ok {
		err := CompareSecrets(passwd, password)
		if err != nil {
			return nil, err
		}

		return &DefaultPrincipal{id: id, name: id}, nil
	}

	return nil, ErrPrincipalNotFound
}

// GetByToken retrieves a principal by token.
//
//nolint:ireturn // Principal is implemented by the repository.
func (r *ConfigRepository) GetByToken(token string) (Principal, error) {
	for i := range r.cfg.Tokens {
		err := CompareSecrets(r.cfg.Tokens[i], token)
		if err == nil {
			return tokenPrincipal, nil
		}
	}

	return nil, ErrPrincipalNotFound
}
