package registry

type Registry interface {
	Register(s Service) error
	Deregister(s Service) error
	List() ([]*Service, error)
	GetService(name string) ([]*Service, error)
}
