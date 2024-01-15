# AIBackend
HTTP server capable of serving local and external models. The server runs each request
asynchronously and returns the model's output. Models can be reached by stdout/stdin or
http requests.

## How to build a request and access a model
Models are accessed by building a JSON request to the server. The implemented models
must respond with JSON. To send a JSON query to a model, POST or GET to
`address:port/{REQUESTTYPE}/{OPERATION}?id={yourmodelid}`, where `yourmodelid` is the
model id and {REQUESTTYPE} can be either `post` or `get`. The supported operations are
explained down below.

## How do models respond
Models have to parse the POST or GET request content (if any) coming from the server and send back the
response to the server.

## Supported operations
There are currently 3 supported operations for the models:
- predict: Used at inference/generation time
- getlogs: Returns the log in json format
- cleanlogs: Cleans the log.
These operations are specified in the `operation` argument of the request.

## How to implement a model
The server is simple to extend. To use a new type of model, it is enough to implement
the `Sender` interface:

    type Sender interface {
        send(body []byte, qt queryType, options ...any) (Json, error) // Send data to the model.
    }

## TODO How to add a new model during runtime
Not implemented (yet!).
To register during runtime a new model, a POST request must be sent with the model id
and explicit type (HttpModel, SubprocessModel).

## What types of model are used
Currently, only models served through HTTP are used. Subprocess models (i.e. models
running as subprocesses of the server) will be added.
