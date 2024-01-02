package back

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

// TODO: rework modeler and basic modeler hierachy ?

// Implemented models. Any structure implementing the modeler inerface is able
// to communicate with a model with the send function. It also obtains all the
// methods from the basicModeler interfac.
type Modeler interface {
	send(body []byte) (Json, error) // Send data to the model.
	Start(*responseFormatter)       // Start the model, make it wait for input and format responses with a responseFormatter
	QueryChannel() InputChan        // returns the channel for the incoming query to this model
	ResponseChannel() OutputChan    // returns the channel for sending back the response
	basicModeler                    // base for every model
}

// Interface for every single model. Every model must be able to predict (whether
// it infers or generates is irrelevant), tell it's state (ready, loading,
// calculating, etc), send the logs informations. A model can be reached by
// HTTP, pipe, ports, FFI.
type basicModeler interface {
	Predict(*FrontEndQuery, func([]byte) (Json, error)) (Json, error) // Sends a query and returns a prediction
	GetState() ModelState                                             // Get the state of the model
	GetLogs(func([]byte) (Json, error)) (Json, error)                 // Get the logs associated with the model
	Id() int
}

// Queries are given to a modeler. They contain necessary information, such as
// the inputs and the instructions given to the model itself.
type FrontEndQuery struct {
	Input     Json      // Information given to the model for prediction. Possibly empty
	QueryType queryType // Type of the query
	id        int       // Identifier for the query
}

// ModelResponse are a wrapper around the returned values of the model. They are
// built with a responseFormatter.
type ModelResponse struct {
	Response     Json         // Returned values of the model.
	ResponseType responseType // Type of the Response.
	id           int          // Identifier for the prediction. It is not part of the returned json
}

const (
	Unknown     responseType = iota // Unknown response type. Default type
	Logs                            // Response containing the model's logs
	Error                           // Response contains an error message
	Predictions                     // Correctly formated response, but without any values in the prediction field
)

const (
	Predict queryType = iota // Query with the goal to predict
	GetLogs                  // Query with the goal of getting the logs
	UnknownQuery
)
