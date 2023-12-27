package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var jsonify = jsoniter.ConfigCompatibleWithStandardLibrary

type base struct {
	state     ModelState // Represent the actual state of the model.
	ModelName string
	id        int
}

// TODO: Reformat
func (b *base) GetLogs(send func([]byte) (*[]byte, error)) (Json, error) {
	if b.state != Ready {
		err := fmt.Errorf("Model is not ready. Model %s with id %v is %s",
			b.ModelName, b.id, b.state)
		return nil, err
	}
	q := Query{
		Instruction: "get logs",
		Input:       nil,
		id:          b.id,
	}
	encoded, err := jsonify.Marshal(q)
	if err != nil {
		b.state = Ready
		return nil, fmt.Errorf("Failed to get the logs: %w", err)
	}
	b.state = Processing
	res, err := send(encoded)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the logs: %w", err)
	}
	var logs Json
	err = jsonify.Unmarshal(*res, logs)
	if err != nil {
		return nil, fmt.Errorf("Failed to get the logs: %w", err)
	}
	return logs, nil
}

// Getter for the state field in model
func (b *base) GetState() ModelState {
	return b.state
}

// Start the prediction computation
// TODO: Reformat
func (b *base) Predict(q *Query, send func([]byte) ([]byte, error)) (*Response, error) {
	if b.state != Ready {
		err := fmt.Errorf("Model is not ready. Model %s with id %v is %s",
			b.ModelName, b.id, b.state)
		return nil, err
	}
	encoded, err := jsonify.Marshal(q)
	if err != nil {
		b.state = Ready
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.state = Processing
	res, err := send(encoded)
	if err != nil {
		b.state = Ready
		//TODO: Better error management in the case of down model.
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	b.state = Ready
	var ans Json
	err = jsonify.Unmarshal(res, ans)
	if err != nil {
		return nil, fmt.Errorf("Failed to predict: %w", err)
	}
	pred := &Response{Response: ans, id: q.id}
	return pred, nil
}

func (b *base) setReady() {
	b.state = Ready
}
func (b *base) setDown() {
	b.state = Down
}
func (b *base) setLoading() {
	b.state = Down
}
func (b *base) setProcessing() {
	b.state = Down
}
