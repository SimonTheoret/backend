package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

// This mock send function returns always return a random json
type testModeler struct {
	basicModeler
	dest string
}

// Mock function for send. Returns random Json
func (tm *testModeler) send([]byte) (Json, error) {

	val, err := gofakeit.JSON(nil)
	if err != nil {
		log.Fatal(err)
	}
	var jsonContent Json
	json.Unmarshal(val, jsonContent)
	return jsonContent, nil
}

// Verifies if testModeler implements the modeler interface
func TestTestModelerInterface(t *testing.T) {
	tested := testModeler{
		basicModeler: base{
			state:     Ready,
			ModelName: "Test interface",
			id:        123,
		},
		dest: "None",
	}
	testP := any(&tested) // Must take & because interfaces are alwaays reference
	_, ok := testP.(Modeler)
	assert.True(t, ok, "testModeler does NOT implement the modeler interface")
}

// Returns a non empty, short but reasonable Json response, similar to a model response
func setUpJsonResponse(t *testing.T) Json {
	return Json{
		"value0": 0.,
		"value1": 1.,
		"value2": 2.,
	}
}

// Returns a random Json response
func setUpRandomJsonResponse(t *testing.T) Json {

	val, err := gofakeit.JSON(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(val)
	var jsonContent Json
	json.Unmarshal(val, jsonContent)
	return jsonContent
}

// Returns a testModeler instance
func setUpTestModeler(t *testing.T) testModeler {
	return testModeler{
		basicModeler: base{state: Ready,
			ModelName: "TestModel",
			id:        10},
		dest: "Not a destination",
	}
}

// Returns an empty Json response
func setUpEmptyJsonResponse(t *testing.T) Json {
	return Json{}
}

func setUpFormatter(t *testing.T, rawResponse Json) responseFormatter {
	return responseFormatter{rawResponse: rawResponse}
}

// Asserts that the raw response is formatted correctly into a ModelResponse
func TestBasicFormattingResponse(t *testing.T) {
	model := setUpTestModeler(t)
	rawResponse := setUpJsonResponse(t)
	rf := setUpFormatter(t, rawResponse)
	response, err := rf.FormatRawResponse(rawResponse, &model)
	if err != nil {
		log.Fatal(err)
	}
	expectedResponse := ModelResponse{Response: rawResponse, ResponseType: Unknown, id: 10}
	assert.Equal(t, &expectedResponse, response, "Both responses should be equal")
}

// Asserts that the preprocess results in the right Json
func TestPreProcessingFormattingResponse(t *testing.T) {
	raw := Json{
		"value0": 0,
		"value1": 1,
		"value2": 2,
	}
	rf := setUpFormatter(t, raw)
	addOne := func(raw Json) (Json, error) {
		for k := range raw {
			raw[k] = raw[k].(int) + 1
		}
		return raw, nil
	}
	newraw, err := addOne(raw)
	if err != nil {
		log.Fatal(fmt.Errorf("Could not apply AddOne: %w", err))
	}
	processed, err := rf.preProcess(raw)
	if err != nil {
		log.Fatal(fmt.Errorf("Could not apply preProcessing: %w", err))
	}
	log.SetOutput(os.Stdout)
	assert.Equal(t, processed, newraw, "Could not preprocess the raw response")

}

// Asserts that formatter builder is OK.
func TestBuildingFormatter(t *testing.T) {
	raw := Json{
		"value0": 0,
		"value1": 1,
		"value2": 2,
	}

	addOne := func(raw Json, fargs ...any) (Json, error) { // Preprocessing func
		for k := range raw {
			raw[k] = raw[k].(int) + 1
		}
		return raw, nil
	}
	funcs := []PreProcessingFunction{addOne}

	actRF := NewFormatter().
		RawResponse(raw).
		PreProcessingFunctions(funcs).
		Build()
	expectedRF := &responseFormatter{raw, funcs, nil, ModelResponse{}}
	assert.Equal(t, expectedRF, actRF, "Should be equal")
}

func TestFormatRawResponse(t *testing.T) {

	preProcessingF := func(raw Json, fargs ...any) (Json, error) { // Preprocessing func
		for k := range raw {
			raw[k] = raw[k].(int) + 1
		}
		return raw, nil
	}
	preFuncs := []PreProcessingFunction{preProcessingF}
	postProcessingF := func(res *ModelResponse, fargs ...any) (*ModelResponse, error) { // Postprocessing func
		for k := range res.Response {
			res.Response[k] = res.Response[k].(int) + 1
		}
		return res, nil
	}
	postFuncs := []PostProcessingFunction{postProcessingF}
	raw := Json{
		"value0": 0,
		"value1": 1,
		"value2": 2,
	}
	rf := NewFormatter().
		RawResponse(raw).
		PreProcessingFunctions(preFuncs).
		PostProcessingFunctions(postFuncs).
		Build()
	model := setUpTestModeler(t)
	actRes, err := rf.FormatRawResponse(raw, &model)
	if err != nil {
		log.Fatal(err)
	}
	expRaw := Json{
		"value0": 2,
		"value1": 3,
		"value2": 4,
	}

	expRes := &ModelResponse{Response: expRaw, ResponseType: Unknown, id: 10}
	assert.Equal(t, actRes, expRes, "Both model responses should be equal")
}
