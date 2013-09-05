package memory

import (
	"testing"
)

type unknown float64

type dummy struct {
	Pk    string
	Data  string
	Age   int
	Thing unknown
}

func (d *dummy) Id() string            { return d.Pk }
func (d *dummy) SetId(s string)        { d.Pk = s }
func (d *dummy) Table() string         { return "dummies" }
func (d *dummy) Convert(i interface{}) { *d = *i.(*dummy) }
func (d *dummy) Validate() error       { return nil }
func (d *dummy) Equal(i interface{}) bool {
	dum := i.(*dummy)
	if dum.Data == d.Data && dum.Age == d.Age {
		return true
	}
	return false
}

func TestMemorySetsIdOnSave(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	obj := &dummy{}
	store.Save(obj)
	if obj.Id() == "" {
		t.Errorf("Store did not set the ID")
	}
}

func TestMemoryStore(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	obj := &dummy{Data: "hello"}
	store.Save(obj)
	obj2, _ := store.Get(obj.Id(), &dummy{})
	if obj != obj2 {
		t.Errorf("Got the wrong object!")
	}
	tag3 := &dummy{Data: "delightful"}
	store.Save(tag3)
	if store.Len() != 2 {
		t.Errorf("Didn't save an obj successfully")
	}
	store.Clear()
	if store.Len() != 0 {
		t.Errorf("Didn't clear successfully")
	}
}

func TestLenOfMemoryStore(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	if store.Len() != 0 {
		t.Errorf("Calculated a bad length")
	}
}

func TestMemoryTables(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	store.Save(&dummy{Data: "hello"})
	if store.Tables() != "dummies\n" {
		t.Errorf("Didn't implement Tables in the expected way.")
	}
}

func TestMemoryNoItem(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	_, err := store.Get("Hello", &dummy{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
}

func TestMemoryClear(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	store.Save(&dummy{Data: "hi"})
	store.Save(&dummy{Data: "bye"})
	if store.Len() != 2 {
		t.Errorf("Store is not returning correct length: %d", store.Len())
	}
	store.Clear()
	if store.Len() != 0 {
		t.Errorf("Store is not clearing correctly: %d", store.Len())
	}
}

var dummies = []*dummy{
	&dummy{
		Data: "one",
		Age:  1,
	},
	&dummy{
		Data: "onerous",
		Age:  2,
	},
	&dummy{
		Data: "Halifax",
		Age:  1,
	},
	&dummy{
		Data:  "Horrible",
		Age:   1,
		Thing: unknown(3.333),
	},
}

func TestMemoryFilter(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	for _, dummy := range dummies {
		store.Save(dummy)
	}
	collection, err := store.Filter("Data", "Halifax", &dummy{})
	if err != nil {
		t.Errorf("Got an error filtering data %s", err)
	}
	if len(collection) != 1 {
		t.Errorf("Filter returned too many things")
	}
	collection, err = store.Filter("Age", 1, &dummy{})
	if err != nil {
		t.Errorf("Error filtering on Age")
	}
	if len(collection) != 3 {
		t.Errorf("Filter did not return the correct thing")
	}
	collection, err = store.Filter("Thing", unknown(3.333), &dummy{})
	if err != nil {
		t.Errorf("Error filtering on unknown type")
	}
	if len(collection) != 0 {
		t.Errorf("We got something and we shouldn't have")
	}
}

func TestMemoryExists(t *testing.T) {
	store := NewMemoryStore(&dummy{})
	for _, dummy := range dummies {
		store.Save(dummy)
	}
	check := &dummy{
		Data: "Horrible",
		Age:  1,
	}
	existance := store.Exists(check)
	if !existance {
		t.Errorf("We didn't find something that exists")
	}
	notIn := &dummy{Data: "definitely new"}
	existance = store.Exists(notIn)
	if existance {
		t.Errorf("Exist is telling us something is there when it is not")
	}

}

func TestTableLen(t *testing.T) {
	table := Table{}
	table["one"] = interface{}("hello")
	if table.Len() != 1 {
		t.Errorf("Table len is not working")
	}
}
