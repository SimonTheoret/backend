package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// Allows to use as if it was the std
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type base struct {
	state     ModelState // Represent the actual state of the model.
	ModelName string
	id        int
}

func (b *base) GetLogs(send func([]byte) (Json, error)) (Json, error) {
	defer b.setReady()
	err := b.verifyIfReady()
	if err != nil {
		return nil, err
	}
	q := Query{
		Input:     nil,
		QueryType: GetLogs,
		id:        b.id,
	}
	logs, err := b.encodeSendDecode(&q, send)
	if err != nil {
		return nil, fmt.Errorf("Encountered error while encoding, sending or decoding: %w", err)
	}
	return logs, nil
}

// Getter for the state field in model
func (b *base) GetState() ModelState {
	return b.state
}

// Start the prediction computation
func (b *base) Predict(q *Query, send func([]byte) (Json, error)) (Json, error) {
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
func (b *base) setReady() {
	b.state = Ready
}

// Utility function for setting the state of the base at Down
func (b *base) setDown() {
	b.state = Down
}

// Utility function for setting the state of the base at Loading
func (b *base) setLoading() {
	b.state = Loading
}

// Utility function for setting the state of the base at Processing
func (b *base) setProcessing() {
	b.state = Processing
}

// Utility function for comparing the state of the model with a given state
func (b *base) isInState(s ModelState) bool {
	return b.GetState() == s
}

// Utility function for comparing the state of the model with a given state.
// It returns the formated error.
func (b *base) verifyIfReady() error {
	if !b.isInState(Ready) {
		err := fmt.Errorf("Model is not ready. Model %s with id %v is %s",
			b.ModelName, b.id, b.state)
		return err
	}
	return nil
}

// Utility function for encoding a query, sending it over and than returning the
// Response. It internally modifies the state of the base
func (b *base) encodeSendDecode(q *Query, send func([]byte) (Json, error)) (Json, error) {
	defer b.setReady()
	encoded, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.state = Processing
	res, err := send(encoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.state = Ready
	return res, err
}
