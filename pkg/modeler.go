package main

// Represent the possible model states.
type ModelState int

// Enums for the type of the response. It is used to process the response.
type responseType int

// Enums for the type of the query. It is used to process the query.
type queryType int

// Structured Json
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

type Modeler interface {
	send([]byte) (Json, error) // Send data to the model.
	basicModeler
}

// Interface for every single model. Every model must be able to predict (whether
// it infers or generates is irrelevant), tell it's state (ready, loading,
// calculating, etc), send the logs informations. A model can be reached by
// HTTP, pipe, ports, FFI.
type basicModeler interface {
	Predict(*Query, func([]byte) (Json, error)) (Json, error) // Sends a query and returns a prediction
	GetState() ModelState                                     // Get the state of the model
	GetLogs(func([]byte) (Json, error)) (Json, error)         // Get the logs associated with the model
}

// Queries are given to a modeler. They contain necessary information, such as
// the inputs and the instructions given to the model itself.
type Query struct {
	Input     Json      // Information given to the model for prediction. Possibly empty
	QueryType queryType // Type of the query
	id        int       // Identifier for the query
}

// Response are a wrapper around the returned values of the model. They are
// built with a responseFormatter.
type Response struct {
	Response     Json         // Returned values of the model.
	ResponseType responseType // Type of the Response.
	id           int          // Identifier for the prediction. It is not part of the returned json
}

const (
	NoPredictions responseType = iota // Correctly formated response, but without any values in the prediction field
	Predictions
	Empty
	Logs
	Error
	Unknown
)

const (
	Predict queryType = iota
	GetLogs
)
