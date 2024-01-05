package back

import (
	"time"

	"github.com/rs/xid"
)

// Enums for the type of the response. It is used to process the response.
type responseType int

// Enums for the type of the query. It is used to process the query.
type queryType int

// Wrapper around a libarary id type.
type Id xid.ID

// Json type as a map
type Json map[string]any

// Any structure implementing the Sender interface is able
// to communicate with an actual model with the send function.
type Sender interface {
	send(body []byte, qt queryType, options ...any) (Json, error) // Send data to the model.
}

// Option type are optional arguments for Predict, GetLogs and CleanLogs methods
// of modeler.
type Options interface {
	BuildOptions() Options
}

// Interface for every single model. Every model must be able to predict
// (whether it infers or generates data is irrelevant), send the logs
// informations. A model can be reached by HTTP, pipe, ports, FFI.
type modeler interface {
	Predict(*message[queryType], Options) (Json, error) // Sends a query and returns a prediction
	GetLogs(Options) (Json, error)                      // Get the logs associated with the model
	CleanLogs(Options) error                            // Cleans the log of the associated model
	start(*responseFormatter) error                     // Start the model, make it wait for input and format responses with a responseFormatter
	queryChannel() InputChan                            // returns the channel for the incoming query to this model
	responseChannel() OutputChan                        // returns the channel for sending back the response
	id() Id                                             // Returns the model id
	Sender
}
type message[T typer] struct {
	content         Json      // JSON content of the message.
	messageType     T         // Type of the message. Should only be responseType or queryType
	id              Id        // Unique id for a message
	creationTime    time.Time // Time when this message was built
	receiverId      Id        // Receiving model's Id
	sender          string    // sender description/IP
	queryOptions    Options   // Options given when sending the message to the model
	responseOptions Options   // Options given when sending back the message
}

func (m *message[T]) ByteContent() ([]byte, error) {
	return json.Marshal(m.content)
}

// Interface for responseType and queryType
type typer interface {
	String() string
	typeString() string
	~int
}

const (
	Unknown     responseType = iota // Unknown response type. Default type
	Logs                            // Response containing the model's logs
	Error                           // Response contains an error message
	Predictions                     // Correctly formated response, but without any values in the prediction field
)

const (
	Predict   queryType = iota // Query with the goal to predict
	GetLogs                    // Query with the goal of getting the logs
	CleanLogs                  //Query with the goal of cleaning the logGetLogss
	UnknownQuery
)

// Convert the response's type (responseType) to string.
func (r responseType) String() string {
	switch r {
	case Unknown:
		return "Unknown"
	case Logs:
		return "Logs"
	case Error:
		return "Error"
	case Predictions:
		return "Predictions"
	}
	return "Undefined response type"
}

// Convert the query's type (queryType) to string.
func (q queryType) String() string {
	switch q {
	case Predict:
		return "Predict"
	case GetLogs:
		return "GetLogs"
	case CleanLogs:
		return "CleanLogs"
	case UnknownQuery:
		return "UnknownQuery"
	}
	return "Undefined query type"
}

// Returns the type to string
func (q queryType) typeString() string {
	return "queryType"
}

// Returns the type to string
func (r responseType) typeString() string {
	return "responseType"
}
