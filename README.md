# Backend
Single HTTP server capable of serving local and external models.
## How to access a model 
To send a JSON query to a model, POST a request to
`address:port/post?id=yourmodelid`, where `yourmodelid` is the model id. 
## How do model respond
Models have to parse the POST request content coming from the server and send
back the response. 

