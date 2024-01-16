# AI Backend
HTTP server capable of serving local and external models
## Motivation
Typical model deployement requires to build and maintain many backends, even possibly a
single backend for each model. The server runs each request asynchronously and returns
the model's output. Models can be reached by stdout/stdin or http requests.

## How to build a request and access a model
Models are accessed by building a JSON request to the server. The implemented models
must respond with JSON. To send a JSON query to a model, POST or GET to
`address:port/{REQUESTTYPE}/{OPERATION}?id={yourmodelid}`, where `yourmodelid` is the
model id and {REQUESTTYPE} can be either `post` or `get`.

An example of a prediction request would be:

    127.0.0.1:8080/post/predict?id=myverygoodmodel

The supported operations for each model are explained down below.

## How do models respond
Models have to parse the POST or GET request content (if any) coming from the server and send back the
response to the server.

## Supported operations
There are currently 3 supported operations for the models:
- predict: Used at inference/generation time
- getlogs: Returns the log in json format
- cleanlogs: Cleans the log.
- addmodel: Add a new model to the server
These operations are specified in the `operation` argument of the request.

## How to implement a model and extend the AI backend
The server is simple to extend. To use a new type of model, it is enough to implement
the `Sender` interface:

    type Sender interface {
        send(body []byte, qt queryType, options ...any) (Json, error) // Send data to the model.
    }
This interface is used to communicate with the model. The arguments `qt` and `options`
can be left untouched, they are used .

## How to add a new model during runtime
This feature is not implemented (yet!).
To register during runtime a new model, a POST request must be sent with the model id
and explicit type (HttpModel, SubprocessModel).

An example of a new model request would be:

    127.0.0.1:8080/post/addmodel?id=myverygoodmodel?type=httpmodel?dest="127.0.0.1:8888"

## What types of model are used
Currently, only models served through HTTP are used. Subprocess models (i.e. models
running as subprocesses of the server) will be added.

## Current limitations
As said previously, the current limitations are no new model during runtime and at the
moment only the HttpModel are available. Both of these restrictions should be lifted
soon.
