package back

import "strconv"

type Id string
type InputChan chan FrontEndQuery
type OutputChan chan ModelResponse

type modelMapper struct {
	InputChannels  map[Id]InputChan  // Access to the modeler's receiving channel by their id.
	OutputChannels map[Id]OutputChan // Access to the modeler's sending channel by their id.
	rF             responseFormatter // Response formatter. Must be build separately with a builder.
}

type ModelMapperBuilder struct {
	mm *modelMapper
}

// Builds the modelMapper.
func (b *ModelMapperBuilder) Build() *modelMapper {
	return b.mm
}

// channels and ids must be the same dimension and have the same order.
func (b *ModelMapperBuilder) InputChan(channels []InputChan, ids []Id) *ModelMapperBuilder {
	var chans map[Id]InputChan
	for i, v := range channels {
		chans[ids[i]] = v
	}
	b.mm.InputChannels = chans
	return b
}

// channels and ids must be the same dimension and have the same order.
func (b *ModelMapperBuilder) OutputChan(channels []OutputChan, ids []Id) *ModelMapperBuilder {
	var chans map[Id]OutputChan
	for i, v := range channels {
		chans[ids[i]] = v
	}
	b.mm.OutputChannels = chans
	return b
}

func setUpModelMapper(outChannels []OutputChan, inChannels []InputChan, ids []Id) *modelMapper {
	mmb := &ModelMapperBuilder{}
	return mmb.InputChan(inChannels, ids).OutputChan(outChannels, ids).Build()
}

// Builds a model mapper based of a slice of modelers and start them. rfs can
// either be a slice of len == 1 or have a responseFormatter for every model
// (i.e. len(models) == len(rfs))
func SetUpModels(models []Sender, rfs ...*responseFormatter) *modelMapper {
	var new_rfs []*responseFormatter
	if len(rfs) == 1 {
		new_rfs = make([]*responseFormatter, 0, len(models))
	} else {

	}
	inChan := make([]InputChan, 0)
	outChan := make([]OutputChan, 0)
	ids := make([]Id, 0)
	for i, m := range models {
		inChan = append(inChan, m.QueryChannel())
		outChan = append(outChan, m.ResponseChannel())
		ids = append(ids, Id(strconv.Itoa(m.Id()))) // ugly conversions
		go m.Start(new_rfs[i])                      // Launch every model in their own goroutine
	}
	return setUpModelMapper(outChan, inChan, ids)
}
