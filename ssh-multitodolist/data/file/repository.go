package file

import "ssh-multitodolist/data"

type Repository struct {
	sourcePath     string
	data           []data.NamedList
	notifyOnCommit func()
	notifyOnRemove func()
}

func New(roomName string, notifyOnCommit func(), notifyOnRemove func()) *Repository {
	return &Repository{"./" + roomName + ".md", []data.NamedList{}, notifyOnCommit, notifyOnRemove}
}

func (p *Repository) List() []data.NamedList {
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

func (p *Repository) Get(index int) data.NamedList {
	if index < 0 || index >= len(p.data) {
		return data.NamedList{}
	}
	return p.data[index]
}
func (p *Repository) GetName(index int) string {
	return p.data[index].Name
}

func (p *Repository) Create(name string) int {
	newList := data.NamedList{Name: name, Items: []data.ListItem{}}
	p.data = append(p.data, newList)
	p.Commit()
	return len(p.data) - 1
}

func (p *Repository) Update(index int, newList data.NamedList) {
	p.data[index] = newList
	p.Commit()
}
func (p *Repository) UpdateName(index int, newName string) {
	p.data[index].Name = newName
	p.Commit()
}

func (p *Repository) Delete(index int) {
	p.data = append(p.data[:index], p.data[index+1:]...)
	p.Commit()
	p.notifyOnRemove()
}

func (p *Repository) Init() error {
	var err error
	p.data, err = readData(p.sourcePath)
	return err
}

func (p *Repository) Commit() error {
	err := writeData(p.data, p.sourcePath)
	p.notifyOnCommit()
	return err
}
