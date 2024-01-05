package back

import (
	"fmt"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
)

// Allows to use as if it was the std
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// The base struct is an implementer of the modeler interface. The base struct allow to easily build new modeler.
type Base struct {
	ModelName string
	id        Id
	in        InputChan
	out       OutputChan
	Sender
}

func NewBase(modelName string, sender Sender) Base {
	return Base{modelName, Id(xid.New()), make(InputChan), make(OutputChan), sender}
}

// Start the prediction computation
func (b *Base) Predict(mess *message[queryType], opt Options) (Json, error) {
	content, err := mess.ByteContent()
	if err != nil {
		return nil,
			fmt.Errorf("Encountered error while predicting: %w", err)
	}
	ans, err := b.send(content, Predict, opt)
	if err != nil {
		return nil,
			fmt.Errorf("Encountered error while predicting: %w", err)
	}
	return ans, nil
}

// Get the logs from the model
func (b *Base) GetLogs(opt Options) (Json, error) {
	logs, err := b.send(nil, GetLogs, opt)
	if err != nil {
		return nil, fmt.Errorf("Encountered error while sending: %w", err)
	}
	return logs, nil
}

// Clean the logs from the model
func (b *Base) CleanLogs(opt Options) error {
	_, err := b.send(nil, CleanLogs, opt)
	if err != nil {
		return fmt.Errorf("Encountered error while sending: %w", err)
	}
	return nil
}

// Make the model connection. The model waits for inputs.
func (b *Base) start(rf *responseFormatter) error {

	for {
		in := <-b.in

		t := in.messageType //Can only be queryType
		opts := in.queryOptions

		switch t {
		case GetLogs: // Get the logs
			res, err := b.GetLogs(opts)
			b.manageErrAndSend(rf, res, err)
		case CleanLogs: // Clean the logs
			err := b.CleanLogs(opts)
			b.manageErrAndSend(rf, nil, err)
		case Predict: // Predict
			res, err := b.Predict(in, opts)
			b.manageErrAndSend(rf, res, err)
		default: //Consider unknown and prediction queries as prediction
			res, err := b.Predict(in, opts)
			if err != nil {
				errBody, _ := json.Marshal(err)
				gin.DefaultErrorWriter.Write(errBody)
			}
			b.manageErrAndSend(rf, res, err)
		}

	}

}

// Helper function. It manages the errors and sends back the *modelResponse.
func (b *Base) manageErrAndSend(rf *responseFormatter, res Json, err error) {
	var responseMessage *message[responseType]
	if err != nil {
		errBody, _ := json.Marshal(err)
		gin.DefaultErrorWriter.Write(errBody)
		jsonError := Json{"error": err.Error()}
		responseMessage, err = rf.FormatRawResponse(jsonError, b, nil)
		if err != nil {
			errBody, _ := json.Marshal(err) // Cannot accept an error when marshaling an error...
			gin.DefaultErrorWriter.Write(errBody)
		}
	} else {
		responseMessage, err = rf.FormatRawResponse(res, b, nil)
		if err != nil {
			errBody, _ := json.Marshal(err)
			gin.DefaultErrorWriter.Write(errBody)
		}
	}
	b.out <- responseMessage
}

// Returns the input channel
func (b *Base) QueryChannel() InputChan {
	return b.in
}

// Returns the output channel
func (b *Base) ResponseChannel() OutputChan {
	return b.out
}

// Return the id of the model
func (b *Base) Id() Id {
	return b.id
}
