package gomdi

import (
	"fmt"
	"testing"
)

// Modify which store we're using since we want to test with a fast datastore
func init() {
	// Implicitly create the tables.
	Store = NewTestStore(&testModel{})
}

// Even more basic than the memory implementation
type testStore map[string]map[string]interface{}

func NewTestStore(models ...Model) testStore {
	m := make(map[string]map[string]interface{})
	for _, model := range models {
		t := make(map[string]interface{})
		m[model.Table()] = t
	}
	return m
}
func (m testStore) CreateTable(model Model) {}
func (m testStore) Save(model Model) error {
	model.SetId("1")
	m[model.Table()][model.Id()] = interface{}(model)
	return nil
}
func (m testStore) Get(id string, model Model) (interface{}, error) {
	return m[model.Table()][id], nil
}
func (m testStore) Filter(field string, value interface{}, model Model) ([]interface{}, error) {
	return []interface{}{}, nil
}
func (m testStore) Exists(model Model) bool {
	if _, ok := m[model.Table()][model.Id()]; ok {
		return true
	}
	return false
}
func (m testStore) Clear()   {}
func (m testStore) Len() int { return 0 }

// Make a fake Model
type testModel struct {
	Pk   string
	Data string
}

func newTestModel(data string) *testModel {
	return &testModel{
		Data: data,
	}
}

func (t *testModel) Id() string            { return t.Pk }
func (t *testModel) SetId(s string)        { t.Pk = s }
func (t *testModel) Convert(i interface{}) { *t = *i.(*testModel) }
func (t *testModel) Table() string         { return "testModels" }
func (t *testModel) Validate() error {
	if t.Data == "ALWAYS FAIL" {
		return fmt.Errorf("Failed.")
	}
	return nil
}
func (t *testModel) Equal(i interface{}) bool {
	test := i.(*testModel)
	return test.Data == t.Data
}

// The Save function is responsible for setting the Id on a model
func TestSave(t *testing.T) {
	model := newTestModel("data")
	err := Save(model)
	if model.Id() == "" {
		t.Errorf("Id is not set correctly")
	}
	if err != nil {
		t.Errorf("Should not be an error saving")
	}
	model = newTestModel("ALWAYS FAIL")
	err = Save(model)
	if err == nil {
		t.Errorf("We should have had an error here")
	}
}

func TestGet(t *testing.T) {
	model := newTestModel("cookie")
	Save(model)
	newModel := &testModel{}
	Get(model.Id(), newModel)
	if newModel.Data != "cookie" {
		t.Errorf("Get is not working properly")
	}
}

func TestFilter(t *testing.T) {
	model := newTestModel("banana")
	Save(model)
	models, err := Filter("Data", "banana", &testModel{})
	if err != nil {
		t.Errorf("Generic filter function is broken")
	}
	if models == nil {
		t.Errorf("Generic filter did not filter properly")
	}
}

func TestExists(t *testing.T) {
	model := newTestModel("banana")
	Save(model)
	if !Exists(model) {
		t.Errorf("Existance check failure")
	}
}
