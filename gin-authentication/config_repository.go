package authentication

var (
	_ Repository = (*ConfigRepository)(nil)
)

// ConfigRepository implements the Repository interface using a Config object.
type ConfigRepository struct {
	cfg *Config
}

// NewConfigRepository creates a new ConfigRepository with the given configuration.
func NewConfigRepository(cfg *Config) *ConfigRepository {
	return &ConfigRepository{cfg: cfg}
}

// GetByToken retrieves a principal by token.
func (r *ConfigRepository) GetByToken(token string) (Principal, error) {
	for i := range r.cfg.Tokens {
		if err := CompareSecrets(r.cfg.Tokens[i], token); err == nil {
			return tokenPrincipal, nil
		}
	}
	return nil, ErrPrincipalNotFound
}

// GetByCredentials retrieves a principal by credentials.
func (r *ConfigRepository) GetByCredentials(id, password string) (Principal, error) {
	if passwd, ok := r.cfg.Principals[id]; ok {
		if err := CompareSecrets(passwd, password); err != nil {
			return nil, err
		}
		return &DefaultPrincipal{id: id, name: id}, nil
	}
	return nil, ErrPrincipalNotFound
}
