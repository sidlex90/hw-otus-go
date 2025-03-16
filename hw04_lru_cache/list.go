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
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.len == 0 {
		l.front = newItem
		l.back = newItem
	} else {
		newItem.Next = l.front
		l.front.Prev = newItem
		l.front = newItem
	}

	l.len++

	return newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}

	if l.len == 0 {
		l.front = newItem
		l.back = newItem
	} else {
		newItem.Prev = l.back
		l.back.Next = newItem
		l.back = newItem
	}

	l.len++

	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		if l.back == i {
			l.back = i.Prev
		}

		i.Prev.Next = i.Next
		i.Prev = nil
	}

	if i.Next != nil {
		if l.front == i {
			l.front = i.Next
		}

		i.Next.Prev = i.Prev
		i.Next = nil
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i == l.front {
		return
	}

	l.Remove(i)
	*i = *l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
