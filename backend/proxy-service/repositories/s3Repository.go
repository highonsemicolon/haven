package repositories

type ProxyRepository interface {
	GetObjectURL(string) string
}

type proxyRepository struct {
	base_path string
}

func NewProxyRepository(base_path string) (ProxyRepository, error) {

	return &proxyRepository{
		base_path: base_path,
	}, nil
}

func (p *proxyRepository) GetObjectURL(key string) string {
	return p.base_path + "/" + key
}
