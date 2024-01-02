# Backend
HTTP server capable of serving local and external models.

## How to access a model 
To send a JSON query to a model, POST a request to
`address:port/post?id=yourmodelid`, where `yourmodelid` is the model id. 

## How do model respond
Models have to parse the POST request content coming from the server and send
back the response. 
### Why only use POST request to communicate with HTTP models
Due to the possibility of accessing both a internal and external model, it is
simpler to send *only* post request. It is not the server's role to treat every
type of request and feature of a model. Instead, these responsibilities are
delegated to the sender and/or the model.

## How to implement a model
Internally, models are represented as the `Modeler` interface. They must
implement 4 functions:

- `send(body []byte) (Json, error) // Send data to the model.`
- `Start(*responseFormatter)       /* Start the model, make it wait for input and format responses with a responseFormatter. Must be launched with a goroutine.*/`
- `QueryChannel() InputChan        // returns the channel for the incoming query to this model`
- `ResponseChannel() OutputChan    // returns the channel for sending back the response`

Finally, models must embed a `basicModeler`. The simplest choice is to embed a
`Base` `struct`.

## TODO How to add a new model
Not implemented (yet)

## What types of model are used
Currently, only models served through HTTP are used. Subprocess models (i.e.
models running as subprocesses of the server) will be added, with the ability to
specify an entry point for downloading models and running them as subprocesses.
