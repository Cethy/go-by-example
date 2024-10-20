package data

type Repository struct {
	sourcePath string
	data       []NamedList
}

/*type Repository interface {
	List() []NamedList
	Get(index int) NamedList
	Create(name string) NamedList
	Update(index int, list NamedList) NamedList
	Delete(index int)
	Init() error
	Commit() error
}*/

func NewRepository(sourcePath string) *Repository {
	p := &Repository{sourcePath, []NamedList{}}
	return p
}

func (p *Repository) List() []NamedList {
	return p.data
}
func (p *Repository) ListNames() []string {
	raw := p.List()
	output := make([]string, len(raw))
	for i, r := range raw {
		output[i] = r.Name
	}
	return output
}

func (p *Repository) Get(index int) NamedList {
	return p.data[index]
}
func (p *Repository) GetName(index int) string {
	return p.data[index].Name
}

func (p *Repository) Create(name string) int {
	newList := NamedList{name, []ListItem{}}
	p.data = append(p.data, newList)
	return len(p.data)
}

func (p *Repository) Update(index int, newList NamedList) {
	p.data[index] = newList
}
func (p *Repository) UpdateName(index int, newName string) {
	p.data[index].Name = newName
}

func (p *Repository) Delete(index int) {
	p.data = append(p.data[:index], p.data[index+1:]...)
}

func (p *Repository) Init() error {
	var err error
	p.data, err = ReadData(p.sourcePath)
	return err
}

func (p *Repository) Commit() error {
	return WriteData(p.data, p.sourcePath)
}
