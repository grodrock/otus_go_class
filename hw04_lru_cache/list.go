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

	// указываем конец списка на элемент, если конца пока нет: size == 0
	if l.backItem == nil {
		l.backItem = &newItem
	}
	// передвигаем первый элемент и смещаем указатель на новый
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
	// указываем начало списка на элемент, если начала пока нет: size == 0
	if l.frontItem == nil {
		l.frontItem = &newItem
	}
	// передвигаем последний элемент и смещаем указатель на новый
	if l.backItem != nil {
		l.backItem.Next = &newItem
	}
	l.backItem = &newItem

	l.size++

	return &newItem
}

// удалить элемент.
func (l *list) Remove(i *ListItem) {
	// смещаем указатель следующего...
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.backItem = i.Prev
	}
	// ... и предыдущего элемента
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

	// соединяем соседей
	i.Prev.Next = i.Next
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	// перемещаем элемент в начало списка
	i.Next = l.frontItem
	i.Prev = nil
	l.frontItem.Prev = i
	l.frontItem = i
}

func NewList() List {
	return new(list)
}
