package data

type ListItem struct {
	Value   string
	Checked bool
}

type NamedList struct {
	Name  string
	Items []ListItem
}

type Repository interface {
	List() []NamedList
	ListNames() []string
	Get(index int) NamedList
	GetName(index int) string
	Create(name string) int
	Update(index int, list NamedList)
	UpdateName(index int, newName string)
	Delete(index int)
	Init() error
	Commit() error
}
