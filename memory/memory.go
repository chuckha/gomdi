package memory

import (
	"fmt"
	"github.com/ChuckHa/gomdi"
	"reflect"
	"strconv"
)

type Table map[string]interface{}

func (t Table) Len() (i int) {
	for _ = range t {
		i++
	}
	return i
}

type MemoryStore map[string]Table

func NewMemoryStore(models ...gomdi.Model) MemoryStore {
	m := make(MemoryStore)
	for _, model := range models {
		m.createTable(model)
	}
	return m
}

// Internal helper to create a table given a model
func (m MemoryStore) createTable(model gomdi.Model) {
	table := make(Table)
	m[model.Table()] = table
}

func (m MemoryStore) Save(model gomdi.Model) error {
	if model.Id() == "" {
		m.setId(model)
	}
	table := m[model.Table()]
	table[model.Id()] = model
	return nil
}

func (m MemoryStore) Get(id string, model gomdi.Model) (interface{}, error) {
	if item, ok := m[model.Table()][id]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("Model from table %s with id %s does not exist", model.Table(), id)
}

// Iterate through each model and check if the field
// we passed is the same as the value.
func (m MemoryStore) Filter(field string, value interface{}, model gomdi.Model) ([]interface{}, error) {
	models := make([]interface{}, 0)
	for _, item := range m[model.Table()] {
		var found bool
		v := reflect.ValueOf(item).Elem()
		// Get the interface value so we can do a type switch
		finterface := v.FieldByName(field).Interface()
		switch finterface.(type) {
		default:
			found = false
		case string:
			found = finterface.(string) == value.(string)
		case int:
			found = finterface.(int) == value.(int)
		}
		if found {
			models = append(models, item)
		}
	}
	return models, nil
}

func (m MemoryStore) Exists(model gomdi.Model) bool {
	for _, item := range m[model.Table()] {
		if model.Equal(item) {
			return true
		}
	}
	return false
}

func (m MemoryStore) Clear() {
	for tableName := range m {
		m[tableName] = make(Table)
	}
	m = make(map[string]Table)
}

func (m MemoryStore) Len() (count int) {
	for tableName := range m {
		for _ = range m[tableName] {
			count++
		}
	}
	return count
}

// Causing a race condition. don't use concurrently, but it's a map so don't do that anyway.
func (m MemoryStore) setId(model gomdi.Model) {
	total := m[model.Table()].Len()
	model.SetId(strconv.Itoa(total))
}

// Print a list of Tables
func (m MemoryStore) Tables() string {
	x := ""
	for k := range m {
		x += fmt.Sprintf("%s\n", k)
	}
	return x
}
