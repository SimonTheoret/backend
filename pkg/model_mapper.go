package back


type InputChan chan *message[queryType]     // These channels send and receive only message containing query
type OutputChan chan *message[responseType] // these channels send and reveive only message containing response

type modelMapper struct {
	InputChannels  map[Id]InputChan  // Access to the modeler's receiving channel by their id.
	OutputChannels map[Id]OutputChan // Access to the modeler's sending channel by their id.
	rF             responseFormatter // Response formatter. Must be build separately with a builder.
}

type ModelMapperBuilder struct {
	mm *modelMapper
}

// Adds the model to the modelmapper
func (b *modelMapper) addNewModel( model modeler, id Id ) {

    b.InputChannels[id] = model.queryChannel()
    b.OutputChannels[id] = model.responseChannel()
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
func SetUpModels(models []modeler, rfs ...*responseFormatter) *modelMapper {
	var new_rfs []*responseFormatter
	if len(rfs) == 1 {
		new_rfs = make([]*responseFormatter, 0, len(models))
	} else {

	}
	inChan := make([]InputChan, 0)
	outChan := make([]OutputChan, 0)
	ids := make([]Id, 0)
	for i, m := range models {
		inChan = append(inChan, m.queryChannel())
		outChan = append(outChan, m.responseChannel())
		ids = append(ids, m.id())
		go m.start(new_rfs[i])                      // Launch every model in their own goroutine
	}
	return setUpModelMapper(outChan, inChan, ids)
}
