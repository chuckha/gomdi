package gomdi

import ()

var (
	Store Datastore
)

// An interface that knows how to talk to the Datastore interface
type Model interface {
	// A getter for the Model's Id.
	Id() string
	// Datastore will set the Id()
	SetId(string)
	// Return the name of the collection, sometimes known as a table name.
	Table() string
	// A method to convert an interface to this model type.
	Convert(interface{})
	// Return nil if the model passes all validation
	Validate() error
	// Define what it means to be equal.
	// This usually starts with a type conversion.
	Equal(interface{}) bool
}

// This datstore knows how to do stuff with types that implement Model
type Datastore interface {
	// Save the model and set the Id
	Save(Model) error
	// Return the model with the given Id
	Get(string, Model) (interface{}, error)
	// Get a collection of models where the field name matches the value passed in
	Filter(string, interface{}, Model) ([]interface{}, error)
	// Test for model existence in the datastore
	Exists(Model) bool
}

// Get the model from the datastore based on Id.
func Get(id string, m Model) error {
	model, err := Store.Get(id, m)
	m.Convert(model)
	return err
}

// Save the model to the datastore.
// Calls Validate on the model before saving.
func Save(m Model) error {
	err := m.Validate()
	if err != nil {
		return err
	}
	return Store.Save(m)
}

// Return an array of interfaces that match a criteria
func Filter(field string, value interface{}, m Model) ([]interface{}, error) {
	return Store.Filter(field, value, m)
}

// Test to see if the model already exists in the datastore
func Exists(m Model) bool {
	return Store.Exists(m)
}
