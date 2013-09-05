# Go Model and Datastore Interfaces

This package provides a Model and Datastore interface.
It's the plug between the two ends.

A model implements the Model interface.
The Model interface talks to Datastore interface.
The Datastore interface is implemented by datastore.

This way you can swap out the pieces that drive your backend
when it makes logical sense to do so.

## Model interface

```go
type Model interface {
	// A getter for the Model's Id.
	Id() string
	// Datastore will set the Id().
	SetId(string)
	// Return the name of the collection, sometimes known as a table name.
	Table() string
	// A method to convert an interface to this model type.
	Convert(interface{})
	// Return nil if the model passes all validation.
	Validate() error
	// Define what it means to be equal.
	// This usually starts with a type conversion.
	Equal(interface{}) bool
}
```

## Datastore interface

```go
type Datastore interface {
        // Save the model and set the Id.
        Save(Model) error
        // Return the model with the given Id.
        Get(string, Model) (interface{}, error)
        // Get a collection of models where the field name matches the value passed in.
        Filter(string, interface{}, Model) ([]interface{}, error)
        // Test for model existence in the datastore.
        Exists(Model) bool
}
```

See an implementation of both Model and Datastore in the test files.
Also memory/ is an implementation of the Datastore interface.
