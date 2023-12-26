package main

type HttpModel struct {
	base    baseModeler
	address string
	port    string
}

// Main implementation detail. This function takes in an array of bytes, such as
// the one resulting from calling json.Marshal() and returns the response from
// the model.
func send([]byte) (*[]byte, error) {
	//TODO: yeah
	return nil, nil
}
