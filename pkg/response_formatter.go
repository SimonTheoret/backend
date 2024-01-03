package back

import "fmt"

type PreProcessingFunction func(raw Json, fargs ...any) (Json, error)
type PostProcessingFunction func(res *ModelResponse, fargs ...any) (*ModelResponse, error)

// Formats the raw response, possibly applies some pre and post processing
// transformations.
type responseFormatter struct {
	rawResponse         Json                     // The raw json, without any modification
	preProcessingFuncs  []PreProcessingFunction  // Transformations applied as preprocessing
	postProcessingFuncs []PostProcessingFunction // Transformations applied as postprocessing
	processedResponse   ModelResponse            // The model response after all processing is completed
}

// This interface is used to specify the way the response is formatted.
type Formatter interface {
	preProcess(Json, ...any) (Json, error)                          //Preprocess the untouched raw response
	postProcess(*ModelResponse, ...any) (*ModelResponse, error)     //PostProcess the newly created Response
	FormatRawResponse(Json, Sender, ...any) (*ModelResponse, error) /* Format the raw response (JSON) into a
	   Response. It is equivalent to chaining preProcess and postProcess*/
}

// Apply in order the preprocessing functions to the raw json. Returns the modified json.
// If any function returns an error, the processs is immediately stopped and returns the error.
// fargs is passed to every function as is.
func (rf *responseFormatter) preProcess(raw Json, fargs ...any) (Json, error) {
	for _, f := range rf.preProcessingFuncs {
		newRaw, err := f(raw, fargs)
		if err != nil {
			return nil, fmt.Errorf("Could not format the response, preprocessing function failed: %w", err)
		}
		raw = newRaw

	}
	return raw, nil
}

// Apply in order the postprocessing functions to the ModelResponse. Returns the
// modified ModelResponse.  If any function returns an error, the processs is
// immediately stopped and returns the error.  fargs is passed to every function
// as is.
func (rf *responseFormatter) postProcess(response *ModelResponse, fargs ...any) (*ModelResponse, error) {
	for _, f := range rf.postProcessingFuncs {
		newResponse, err := f(response, fargs)
		if err != nil {
			return nil, fmt.Errorf("Could not format the response, postprocessing function failed: %w", err)
		}
		response = newResponse

	}
	return response, nil
}

// Returns the expected
func (rf *responseFormatter) findResponseType(content Json) responseType {
	size := len(content)
	if size < 1 {
		return Error
	}
	helper := func(w string) bool {
		_, ok := content[w]
		return ok
	}
	switch {
	case helper("predictions"):
		return Predictions
	case helper("output"):
		return Predictions
	case helper("out"):
		return Predictions
	case helper("predict"):
		return Predictions
	case helper("values"):
		return Predictions
	case helper("logs"):
		return Logs
	case helper("errors"):
		return Error
	case helper("Errors"):
		return Error
	case helper("error"):
		return Error
	case helper("Error"):
		return Error
	default:
		return Unknown
	}
}

// Chains preProcess and postProcess. It returns a pointer to the newly created
// ModelResponse.  If any function returns an error, the processs is immediately
// stopped and returns the error.  fargs is passed to every function as is.
func (rf *responseFormatter) FormatRawResponse(rawResponse Json, model Sender, fargs ...any) (*ModelResponse, error) {
	var preProcessed Json
	var err error
	resType := rf.findResponseType(rawResponse)
	if rf.preProcessingFuncs != nil {
		preProcessed, err = rf.preProcess(rawResponse, fargs)
		if err != nil {
			return nil, fmt.Errorf("Could not format the raw response: %w", err)
		}
	} else {
		preProcessed = rawResponse
		err = nil
	}
	modelID := model.Id()
	modelResponse := ModelResponse{Response: preProcessed, ResponseType: resType, Id: modelID}
	var postProcessed *ModelResponse
	if rf.postProcessingFuncs != nil {
		postProcessed, err = rf.postProcess(&modelResponse, fargs)
		if err != nil {
			return nil, fmt.Errorf("Could not format the raw response: %w", err)
		}
	} else {
		postProcessed = &modelResponse
		err = nil
	}
	rf.processedResponse = *postProcessed
	return postProcessed, nil
}

// Utility struct made to chain methods for building formatters
type FormatterBuilder struct {
	rf responseFormatter
}

// Starts the chaining of options, with the end goal of building a responseFormatter.
func NewFormatter() *FormatterBuilder {
	rf := responseFormatter{rawResponse: nil,
		preProcessingFuncs:  nil,
		postProcessingFuncs: nil,
		processedResponse:   ModelResponse{}}
	return &FormatterBuilder{rf: rf}
}

// Sets the raw response of the formatter. Used to chain the options.
func (fb *FormatterBuilder) RawResponse(rawResponse Json) *FormatterBuilder {
	fb.rf.rawResponse = rawResponse
	return fb
}

// Sets the preprocessing functions of the formatter. Used to chain the options.
func (fb *FormatterBuilder) PreProcessingFunctions(funcs []PreProcessingFunction) *FormatterBuilder {
	fb.rf.preProcessingFuncs = funcs
	return fb
}

// Sets the postprocessing functions of the formatter. Used to chain the options.
func (fb *FormatterBuilder) PostProcessingFunctions(funcs []PostProcessingFunction) *FormatterBuilder {
	fb.rf.postProcessingFuncs = funcs
	return fb
}

// Build the responseFormatter. Used to chain the options.
func (fb *FormatterBuilder) Build() *responseFormatter {
	rf := fb.rf
	return &rf
}

func DefaultFormatter() *responseFormatter {
	return NewFormatter().Build()
}
