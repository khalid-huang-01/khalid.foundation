package sort

import (
	"sort"
	"testing"
)

// sortSlice
func TestSortSlice(t *testing.T) {
	array := []int{3, 1, 5, 2, 7, 33}
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	t.Log(array)
}

// sort for []person
type Person struct {
	name string
	age int
}

func TestSortPersonSlice(t *testing.T) {
	persons := []Person{
		{age: 4, name: "ke"},
		{age: 2, name: "ke"},
		{age: 11, name: "ke"},
		{age: 90, name: "ke"},
		{age: 55, name: "ke"},
	}
	sort.Slice(persons, func(i, j int) bool {
		return persons[i].age < persons[j].age
	})
	t.Log(persons)
}

// sort for persons struct
type Persons []Person

func (ps Persons) Len() int {
	return len(ps)
}
// Less reports whether the element with
// index i should sort before the element with index j.
func (ps Persons) Less(i, j int) bool {
	return ps[i].age < ps[j].age
}
// Swap swaps the elements with indexes i and j.
func (ps Persons) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func TestPersonsSort(t *testing.T)  {
	persons := Persons{
		{age: 4, name: "ke"},
		{age: 2, name: "ke"},
		{age: 11, name: "ke"},
		{age: 90, name: "ke"},
		{age: 55, name: "ke"},
	}
	sort.Sort(persons)
	t.Log(persons)
}