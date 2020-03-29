package revindex

import "sort"

// Empty struct needs zero bytes
type Void struct{}

// Set to store file indexes
type Set map[int]Void

func (s *Set) Put(val int) {
	(*s)[val] = Void{}
}

func (s *Set) PutAll(values []int) {
	for _, val := range values {
		s.Put(val)
	}
}

func (s *Set) SortedKeys() []int {
	keys := make([]int, 0, len(*s))
	for key := range *s {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	return keys
}
