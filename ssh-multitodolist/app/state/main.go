package state

// Shared state of the application

type State struct {
	Username       string
	Color          string
	editTab        int
	activeTab      int
	removingTab    int
	cursor         int // which to-do list item our Cursor is pointing at
	previousCursor int // which to-do list item our Cursor is pointing at (before input is active)
	editCursor     int
	notify         func()
}

func (s *State) EditTab(v int) {
	oldV := s.editTab
	s.editTab = v

	if oldV != s.editTab {
		s.notify()
	}
}
func (s *State) GetEditTab() int {
	return s.editTab
}
func (s *State) ActiveTab(v int) {
	oldV := s.activeTab
	s.activeTab = v

	if oldV != s.activeTab {
		s.notify()
	}
}
func (s *State) GetActiveTab() int {
	return s.activeTab
}
func (s *State) RemovingTab(v int) {
	oldV := s.removingTab
	s.removingTab = v

	if oldV != s.removingTab {
		s.notify()
	}
}
func (s *State) GetRemovingTab() int {
	return s.removingTab
}

func (s *State) Cursor(v int) {
	oldV := s.cursor
	s.cursor = v

	if oldV != s.cursor {
		s.notify()
	}
}
func (s *State) GetCursor() int {
	return s.cursor
}
func (s *State) PreviousCursor(v int) {
	oldV := s.previousCursor
	s.previousCursor = v

	if oldV != s.previousCursor {
		s.notify()
	}
}
func (s *State) GetPreviousCursor() int {
	return s.previousCursor
}
func (s *State) EditCursor(v int) {
	oldV := s.editCursor
	s.editCursor = v

	if oldV != s.editCursor {
		s.notify()
	}
}
func (s *State) GetEditCursor() int {
	return s.editCursor
}

func New(username string, color string, notify func()) *State {
	return &State{
		Username:    username,
		Color:       color,
		editTab:     -1,
		removingTab: -1,
		editCursor:  -1,
		notify:      notify,
	}
}
