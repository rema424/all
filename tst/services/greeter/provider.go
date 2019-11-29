package greeter

// Provider ...
type Provider struct {
	db Database
}

// NewProvider ...
func NewProvider(db Database) *Provider {
	return &Provider{db}
}
