package hw04lrucache

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func Test_list_PushFront(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name               string
		createList         func() List
		args               args
		expectedValue      interface{}
		expectedLength     int
		expectedFrontValue interface{}
		expectedBackValue  interface{}
	}{
		{
			name: "empty list",
			createList: func() List {
				return NewList()
			},
			args:               args{v: 10},
			expectedValue:      10,
			expectedLength:     1,
			expectedFrontValue: 10,
			expectedBackValue:  10,
		},
		{
			name: "one item list",
			createList: func() List {
				l := NewList()
				l.PushFront(20)
				return l
			},
			args:               args{v: 10},
			expectedValue:      10,
			expectedLength:     2,
			expectedFrontValue: 10,
			expectedBackValue:  20,
		},
		{
			name: "two items list",
			createList: func() List {
				l := NewList()
				l.PushFront(3)
				l.PushFront(2)
				return l
			},
			args:               args{v: 1},
			expectedValue:      1,
			expectedLength:     3,
			expectedFrontValue: 1,
			expectedBackValue:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.createList()
			item := l.PushFront(tt.args.v)
			require.NotNil(t, item)
			if !reflect.DeepEqual(item.Value, tt.expectedValue) {
				t.Errorf("PushFront() = %v, want %v", item.Value, tt.expectedValue)
			}
			if length := l.Len(); !reflect.DeepEqual(length, tt.expectedLength) {
				t.Errorf("Len() = %v, want %v", length, tt.expectedLength)
			}
			frontItem := l.Front()
			require.NotNil(t, frontItem)
			if !reflect.DeepEqual(frontItem.Value, tt.expectedFrontValue) {
				t.Errorf("Front() = %v, want %v", frontItem.Value, tt.expectedFrontValue)
			}
			backItem := l.Back()
			require.NotNil(t, backItem)
			if !reflect.DeepEqual(backItem.Value, tt.expectedBackValue) {
				t.Errorf("Back() = %v, want %v", backItem.Value, tt.expectedBackValue)
			}
		})
	}
}

func Test_list_PushBack(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name               string
		createList         func() List
		args               args
		expectedValue      interface{}
		expectedLength     int
		expectedFrontValue interface{}
		expectedBackValue  interface{}
	}{
		{
			name: "empty list",
			createList: func() List {
				return NewList()
			},
			args:               args{v: 10},
			expectedValue:      10,
			expectedLength:     1,
			expectedFrontValue: 10,
			expectedBackValue:  10,
		},
		{
			name: "one item list",
			createList: func() List {
				l := NewList()
				l.PushBack(10)
				return l
			},
			args:               args{v: 20},
			expectedValue:      20,
			expectedLength:     2,
			expectedFrontValue: 10,
			expectedBackValue:  20,
		},
		{
			name: "two items list",
			createList: func() List {
				l := NewList()
				l.PushBack(10)
				l.PushBack(20)
				return l
			},
			args:               args{v: 30},
			expectedValue:      30,
			expectedLength:     3,
			expectedFrontValue: 10,
			expectedBackValue:  30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := tt.createList()
			item := l.PushBack(tt.args.v)
			require.NotNil(t, item)
			if !reflect.DeepEqual(item.Value, tt.expectedValue) {
				t.Errorf("PushBack() = %v, want %v", item.Value, tt.expectedValue)
			}
			if length := l.Len(); !reflect.DeepEqual(length, tt.expectedLength) {
				t.Errorf("Len() = %v, want %v", length, tt.expectedLength)
			}
			frontItem := l.Front()
			require.NotNil(t, frontItem)
			if !reflect.DeepEqual(frontItem.Value, tt.expectedFrontValue) {
				t.Errorf("Front() = %v, want %v", frontItem.Value, tt.expectedFrontValue)
			}
			backItem := l.Back()
			require.NotNil(t, backItem)
			if !reflect.DeepEqual(backItem.Value, tt.expectedBackValue) {
				t.Errorf("Back() = %v, want %v", backItem.Value, tt.expectedBackValue)
			}
		})
	}
}

func Test_list_Remove(t *testing.T) {
	tests := []struct {
		name                   string
		init                   func() (List, []*ListItem)
		remove                 func(*List, []*ListItem)
		expectedLength         int
		expectedValues         []interface{}
		expectedValuesReversed []interface{}
		expectedFrontValue     interface{}
		expectedBackValue      interface{}
	}{
		{
			name: "remove an item from an one item list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10))
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[0])
			},
			expectedLength:         0,
			expectedValues:         []interface{}{},
			expectedValuesReversed: []interface{}{},
			expectedFrontValue:     nil,
			expectedBackValue:      nil,
		},
		{
			name: "remove the first item from a two items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[0])
			},
			expectedLength:         1,
			expectedValues:         []interface{}{20},
			expectedValuesReversed: []interface{}{20},
			expectedFrontValue:     20,
			expectedBackValue:      20,
		},
		{
			name: "remove the last item from a two items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[1])
			},
			expectedLength:         1,
			expectedValues:         []interface{}{10},
			expectedValuesReversed: []interface{}{10},
			expectedFrontValue:     10,
			expectedBackValue:      10,
		},
		{
			name: "remove the first item from a three items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				items = append(items, l.PushBack(30)) // [10, 20, 30]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[0])
			},
			expectedLength:         2,
			expectedValues:         []interface{}{20, 30},
			expectedValuesReversed: []interface{}{30, 20},
			expectedFrontValue:     20,
			expectedBackValue:      30,
		},
		{
			name: "remove the last item from three items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				items = append(items, l.PushBack(30)) // [10, 20, 30]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[2])
			},
			expectedLength:         2,
			expectedValues:         []interface{}{10, 20},
			expectedValuesReversed: []interface{}{20, 10},
			expectedFrontValue:     10,
			expectedBackValue:      20,
		},
		{
			name: "remove the middle item from three items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				items = append(items, l.PushBack(30)) // [10, 20, 30]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(items[1])
			},
			expectedLength:         2,
			expectedValues:         []interface{}{10, 30},
			expectedValuesReversed: []interface{}{30, 10},
			expectedFrontValue:     10,
			expectedBackValue:      30,
		},
		{
			name: "remove an item that is not in the list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				items = append(items, l.PushBack(30)) // [10, 20, 30]
				return l, items
			},
			remove: func(l *List, items []*ListItem) {
				(*l).Remove(&ListItem{Value: 20})
			},
			expectedLength:         3,
			expectedValues:         []interface{}{10, 20, 30},
			expectedValuesReversed: []interface{}{30, 20, 10},
			expectedFrontValue:     10,
			expectedBackValue:      30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, items := tt.init()

			tt.remove(&l, items)

			require.Equal(t, tt.expectedLength, l.Len())
			require.Equal(t, tt.expectedFrontValue, getItemValue(l.Front()))
			require.Equal(t, tt.expectedBackValue, getItemValue(l.Back()))
			require.Equal(t, tt.expectedValues, getListValues(&l))
			require.Equal(t, tt.expectedValuesReversed, getListValuesReversed(&l))
		})
	}
}

func Test_list_MoveToFront(t *testing.T) {
	tests := []struct {
		name                   string
		init                   func() (List, []*ListItem)
		moveToFront            func(*List, []*ListItem)
		expectedLength         int
		expectedValues         []interface{}
		expectedValuesReversed []interface{}
		expectedFrontValue     interface{}
		expectedBackValue      interface{}
	}{
		{
			name: "move to front an item from an one item list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				return l, items
			},
			moveToFront: func(l *List, items []*ListItem) {
				(*l).MoveToFront(items[0])
			},
			expectedLength:         1,
			expectedFrontValue:     10,
			expectedBackValue:      10,
			expectedValues:         []interface{}{10},
			expectedValuesReversed: []interface{}{10},
		},
		{
			name: "move to front the last item from a two items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(20)) // [20]
				items = append(items, l.PushBack(10)) // [20, 10]
				return l, items
			},
			moveToFront: func(l *List, items []*ListItem) {
				(*l).MoveToFront(items[1])
			},
			expectedLength:         2,
			expectedFrontValue:     10,
			expectedBackValue:      20,
			expectedValues:         []interface{}{10, 20},
			expectedValuesReversed: []interface{}{20, 10},
		},
		{
			name: "move to front the last item from a three items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(20)) // [20]
				items = append(items, l.PushBack(30)) // [20, 30]
				items = append(items, l.PushBack(10)) // [20, 30, 10]
				return l, items
			},
			moveToFront: func(l *List, items []*ListItem) {
				(*l).MoveToFront(items[2])
			},
			expectedLength:         3,
			expectedFrontValue:     10,
			expectedBackValue:      30,
			expectedValues:         []interface{}{10, 20, 30},
			expectedValuesReversed: []interface{}{30, 20, 10},
		},
		{
			name: "move to front the middle item from a three items list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(20)) // [20]
				items = append(items, l.PushBack(10)) // [20, 10]
				items = append(items, l.PushBack(30)) // [20, 10, 30]
				return l, items
			},
			moveToFront: func(l *List, items []*ListItem) {
				(*l).MoveToFront(items[1])
			},
			expectedLength:         3,
			expectedFrontValue:     10,
			expectedBackValue:      30,
			expectedValues:         []interface{}{10, 20, 30},
			expectedValuesReversed: []interface{}{30, 20, 10},
		},
		{
			name: "move to front an item that is not in the list",
			init: func() (List, []*ListItem) {
				items := make([]*ListItem, 0)
				l := NewList()
				items = append(items, l.PushBack(10)) // [10]
				items = append(items, l.PushBack(20)) // [10, 20]
				items = append(items, l.PushBack(30)) // [10, 20, 30]
				return l, items
			},
			moveToFront: func(l *List, items []*ListItem) {
				(*l).MoveToFront(&ListItem{Value: 30})
			},
			expectedLength:         3,
			expectedFrontValue:     10,
			expectedBackValue:      30,
			expectedValues:         []interface{}{10, 20, 30},
			expectedValuesReversed: []interface{}{30, 20, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, items := tt.init()

			tt.moveToFront(&l, items)

			require.Equal(t, tt.expectedLength, l.Len())
			require.Equal(t, tt.expectedFrontValue, getItemValue(l.Front()))
			require.Equal(t, tt.expectedBackValue, getItemValue(l.Back()))
			require.Equal(t, tt.expectedValues, getListValues(&l))
			require.Equal(t, tt.expectedValuesReversed, getListValuesReversed(&l))
		})
	}
}

func getListValues(l *List) []interface{} {
	values := make([]interface{}, 0)
	item := (*l).Front()
	for item != nil {
		values = append(values, item.Value)
		item = item.Next
	}
	return values
}

func getListValuesReversed(l *List) []interface{} {
	values := make([]interface{}, 0)
	item := (*l).Back()
	for item != nil {
		values = append(values, item.Value)
		item = item.Prev
	}
	return values
}

func getItemValue(i *ListItem) interface{} {
	if i == nil {
		return nil
	}
	return i.Value
}
