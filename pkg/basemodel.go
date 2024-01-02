package back

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/xid"
)

// Allows to use as if it was the std
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Base struct {
	state     ModelState // Represent the actual state of the model.
	ModelName string
	id        int
}

func (b Base) GetLogs(send func([]byte) (Json, error)) (Json, error) {
	defer b.setReady()
	err := b.verifyIfReady()
	if err != nil {
		return nil, err
	}
	q := FrontEndQuery{
		Input:     nil,
		QueryType: GetLogs,
		Id:        xid.New().String(),
	}
	logs, err := b.encodeSendDecode(&q, send)
	if err != nil {
		return nil, fmt.Errorf("Encountered error while encoding, sending or decoding: %w", err)
	}
	return logs, nil
}

// Getter for the state field in model
func (b Base) GetState() ModelState {
	return b.state
}

// Start the prediction computation
func (b Base) Predict(q *FrontEndQuery, send func([]byte) (Json, error)) (Json, error) {
	defer b.setReady()
	err := b.verifyIfReady()
	if err != nil {
		return nil, err
	}
	ans, err := b.encodeSendDecode(q, send)
	if err != nil {
		return nil,
			fmt.Errorf("Encountered error while encoding, sending or decoding: %w", err)
	}
	return ans, nil
}

// Utility function for setting the state of the base at Ready
func (b *Base) setReady() {
	b.state = Ready
}

// Utility function for setting the state of the base at Down
func (b *Base) setDown() {
	b.state = Down
}

// Utility function for setting the state of the base at Loading
func (b *Base) setLoading() {
	b.state = Loading
}

// Utility function for setting the state of the base at Processing
func (b *Base) setProcessing() {
	b.state = Processing
}

// Utility function for comparing the state of the model with a given state
func (b *Base) isInState(s ModelState) bool {
	return b.GetState() == s
}

// Utility function for comparing the state of the model with a given state.
// It returns the formated error.
func (b *Base) verifyIfReady() error {
	if !b.isInState(Ready) {
		err := fmt.Errorf("Model is not ready. Model %s with id %v is %s",
			b.ModelName, b.id, b.state)
		return err
	}
	return nil
}

// Utility function for encoding a query, sending it over and than returning the
// Response. It internally modifies the state of the base
func (b *Base) encodeSendDecode(q *FrontEndQuery, send func([]byte) (Json, error)) (Json, error) {
	defer b.setReady()
	encoded, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.setProcessing()
	res, err := send(encoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.setReady()
	return res, err
}

func (b Base) Id() int {
	return int(b.id)
}

func NewBase(name string, id int) Basicmodeler {
	b := Base{Ready, name, id}
	return b
}
