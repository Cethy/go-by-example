package data

type ListItem struct {
	Value   string `json:"value"`
	Checked bool   `json:"checked"`
}

type NamedList struct {
	Name  string     `json:"name"`
	Items []ListItem `json:"items"`
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
}

type Storage interface {
	Init() ([]NamedList, error)
	Commit(data []NamedList) error
}

type LocalCopyRepository struct {
	storage        Storage
	data           []NamedList
	notifyOnCommit func()
	notifyOnRemove func()
}

func New(storage Storage, notifyOnCommit func(), notifyOnRemove func()) (*LocalCopyRepository, error) {
	r := &LocalCopyRepository{storage, []NamedList{}, notifyOnCommit, notifyOnRemove}
	data, err := r.storage.Init()
	r.data = data

	return r, err
}

func (p *LocalCopyRepository) List() []NamedList {
	return p.data
}
func (p *LocalCopyRepository) ListNames() []string {
	raw := p.List()
	output := make([]string, len(raw))
	for i, r := range raw {
		output[i] = r.Name
	}
	return output
}

func (p *LocalCopyRepository) Get(index int) NamedList {
	if index < 0 || index >= len(p.data) {
		return NamedList{}
	}
	return p.data[index]
}
func (p *LocalCopyRepository) GetName(index int) string {
	return p.data[index].Name
}

func (p *LocalCopyRepository) Create(name string) int {
	newList := NamedList{Name: name, Items: []ListItem{}}
	p.data = append(p.data, newList)
	p.commit()
	return len(p.data) - 1
}

func (p *LocalCopyRepository) Update(index int, newList NamedList) {
	p.data[index] = newList
	p.commit()
}
func (p *LocalCopyRepository) UpdateName(index int, newName string) {
	p.data[index].Name = newName
	p.commit()
}

func (p *LocalCopyRepository) Delete(index int) {
	p.data = append(p.data[:index], p.data[index+1:]...)
	p.commit()
	p.notifyOnRemove()
}

func (p *LocalCopyRepository) commit() error {
	err := p.storage.Commit(p.data)
	if err != nil {
		p.notifyOnCommit()
	}
	return err
}
