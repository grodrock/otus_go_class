package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	size                int
	frontItem, backItem *ListItem
}

// длина списка.
func (l *list) Len() int {
	return l.size
}

// первый элемент списка.
func (l *list) Front() *ListItem {
	return l.frontItem
}

// последний элемент списка.
func (l *list) Back() *ListItem {
	return l.backItem
}

// добавить значение в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{
		Value: v,
		Next:  l.frontItem,
	}

	if l.backItem == nil {
		l.backItem = &newItem
	}
	if l.frontItem != nil {
		l.frontItem.Prev = &newItem
	}

	l.frontItem = &newItem
	l.size++

	return &newItem
}

// добавить значение в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{
		Value: v,
		Prev:  l.backItem,
	}
	if l.frontItem == nil {
		l.frontItem = &newItem
	}
	if l.backItem != nil {
		l.backItem.Next = &newItem
	}

	l.backItem = &newItem
	l.size++

	return &newItem
}

// удалить элемент.
func (l *list) Remove(i *ListItem) {
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.frontItem = i.Prev
	}
	l.size--
}

// переместить элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		// мы и так первые.
		return
	}

	i.Prev.Next = i.Next
	i.Next = l.frontItem
	i.Prev = nil
	l.frontItem.Prev = i
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
