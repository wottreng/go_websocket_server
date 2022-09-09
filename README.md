
# websocket pub/sub microservice
websocket pub/sub style microservice server written in golang

## Usage
* run server: `go run server.go`
* connect to the server with arguments `?topic=<some_topic>`
* send messages to the server for that topic
* receive messages from the server for that topic

## Features
* supports multiple topics
* supports multiple clients per topic
* optional api key define in `/utils/system_utils/system_utils.go`
* multithreaded
* cmd line arguments for development(bind to localhost) and production(bind to public address)
* logging output to file

## testing
basic example html served from `/test` route to test websocket connection

## layout
* server.go: main file
  * starts websocket server
  * contains the handler function
* file_utils.go: file utilities
  * contains functions to read and write files
* time_utils.go: time utilities
  * contains functions to get the current time and date
* http_utils.go: http utilities
  * contains functions to handle http requests
  * contains functions to handle websocket requests
* system_utils.go: system utilities
  * handles cmd line arguments

## websocket server library
* [gorilla websocket](https://github.com/gorilla/websocket)
    * Gorilla WebSocket is a Go implementation of the WebSocket protocol.

Cheers,
Mark