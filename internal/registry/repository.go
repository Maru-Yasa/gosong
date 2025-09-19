package registry

type Repository interface {
	Save(app AppState) error
	Find(name string) (*AppState, error)
	FindAll() ([]AppState, error)
	Delete(name string) error
}
