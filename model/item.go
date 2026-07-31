package model

type TaskStatus string

const (
	StatusNone       TaskStatus = ""
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Item struct {
	ID       string      `json:"id"`
	Text     string      `json:"text"`
	IsTask   bool        `json:"is_task"`
	Status   TaskStatus  `json:"status"`
	Folded   bool        `json:"folded"`
	Children []*Item     `json:"children,omitempty"`
	Parent   *Item       `json:"-"`
}

func NewItem(id, text string) *Item {
	return &Item{
		ID:       id,
		Text:     text,
		IsTask:   false,
		Status:   StatusNone,
		Folded:   false,
		Children: []*Item{},
	}
}

func NewTask(id, text string, status TaskStatus) *Item {
	if status == StatusNone {
		status = StatusTodo
	}
	return &Item{
		ID:       id,
		Text:     text,
		IsTask:   true,
		Status:   status,
		Folded:   false,
		Children: []*Item{},
	}
}

func (i *Item) Clone() *Item {
	if i == nil {
		return nil
	}
	newItem := &Item{
		ID:     i.ID,
		Text:   i.Text,
		IsTask: i.IsTask,
		Status: i.Status,
		Folded: i.Folded,
	}
	for _, child := range i.Children {
		childClone := child.Clone()
		childClone.Parent = newItem
		newItem.Children = append(newItem.Children, childClone)
	}
	return newItem
}

type VisibleItem struct {
	Item        *Item
	Depth       int
	Index       int // index in flat visible list
	HasChildren bool
	Parent      *Item
}
