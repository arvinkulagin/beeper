# beeper
Pub/Sub messaging server using websockets
## Installation
`go get github.com/arvinkulagin/beeper/...`
## Usage
### Start server
`user@host$ beeper -ws localhost:8888 -rest localhost:8889`
### Start command line interface
`user@host$ beeper-cli -s localhost:8889`
`>> help`
### Add topic
`>> add test0`
### Subscribe to topic
Connect to ws://localhost:8888/test0
### Publish message
`>> pub test0 Hello World!`
### Delete topic
`>> del test0`
## API
You can get golang REST API wrapper [here](https://github.com/arvinkulagin/beeperapi).
### Add topic
`POST http://localhost:8889/topic` with topic ID in request body
### Publish message
`POST http://localhost:8889/topic/{id}` with message in request body
### Delete topic
`DELETE http://localhost:8889/topic/{id}`