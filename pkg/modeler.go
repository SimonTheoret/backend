package main

// Represent the possible model states.
type ModelState int
type Json map[string]any

const (
	Ready ModelState = iota
	Loading
	Processing
	Waiting
	Down
)

// Returns the model's state to string.
func (s ModelState) String() string {
	switch s {
	case Ready:
		return "Ready"
	case Loading:
		return "Loading"
	case Processing:
		return "Processing"
	case Down:
		return "Down"
	}
	return "Undefined state"
}

type modeler interface {
	send([]byte) (*[]byte, error) // Send data to the model. Used in predict
	baseModeler
}

// Interface for any single model. Every model must be able to predict (whether
// it infers or generates is irrelevant), tell it's state (ready, loading,
// calculating, etc), send the logs informations. A model can be reached by
// HTTP, pipe, ports, FFI.
type baseModeler interface {
	Predict(*Query, func([]byte) ([]byte, error)) (*Prediction, error) // Sends a query and returns a prediction
	GetState() ModelState                                              // Get the state of the model
	GetLogs(func([]byte) ([]byte, error)) (*Json, error)               // Get the logs associated with the model
}

// Queries are given to a modeler. They contain necessary information, such as
// the inputs and the instructions given to the model itself.
type Query struct {
	Instruction string // What is the model supposed to do
	Input       *Json  // Information given to the model for prediction
	id          int    // Identifier for the query
}

// Predictions are returned by the modeler and contain the result
type Prediction struct {
	Answer *Json // Return values of the model
	id     int   // Identifier for the prediction. It is not part of the returned json
}
