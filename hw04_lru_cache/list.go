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
	list  *list
}

type list struct {
	length    int
	frontItem *ListItem
	backItem  *ListItem
}

func NewList() List {
	return new(list)
}

func (l list) Len() int {
	return l.length
}

func (l list) Front() *ListItem {
	return l.frontItem
}

func (l list) Back() *ListItem {
	return l.backItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v, Next: l.frontItem, list: l}
	if l.frontItem != nil {
		l.frontItem.Prev = &item
	}
	l.frontItem = &item
	l.length++
	if l.backItem == nil {
		l.backItem = &item
	}
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v, Prev: l.backItem, list: l}
	if l.backItem != nil {
		l.backItem.Next = &item
	}
	l.backItem = &item
	l.length++
	if l.frontItem == nil {
		l.frontItem = &item
	}
	return &item
}

func (l *list) Remove(i *ListItem) {
	if i.list != l {
		return
	}
	if i == l.frontItem {
		if i == l.backItem {
			l.frontItem = nil
			l.backItem = nil
		} else {
			l.frontItem = i.Next
			l.frontItem.Prev = i.Prev
		}
		l.length--
		return
	}
	if i == l.backItem {
		l.backItem = i.Prev
		l.backItem.Next = i.Next
		l.length--
		return
	}
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.list != l || l.frontItem == i {
		return
	}
	if i == l.backItem {
		l.backItem = i.Prev
		l.backItem.Next = nil
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Next = l.frontItem
	l.frontItem.Prev = i
	l.frontItem = i
	l.frontItem.Prev = nil
}
