# Golang Proxy

FastSMSing only operates using three statuses:

* `CONFIRMED` - used when FastSMSing receives a given request and schedules the SMS delivery.
* `FAILED` - used when FastSMSing fails to deliver a message after it has already been `Confirmed`.
* `DELIVERED` - used when the delivery is successful.

Your service provides five different statuses that might be returned via REST API:

* `ACCEPTED` - used when your service accepts a message to be sent in the nearest batch of messages.
* `NOT_FOUND` - used for messages that have not been planned to be sent.
* `CONFIRMED` - used for messages with the `CONFIRMED` status, received via FastSMSing updates mechanism.
* `FAILED` - used for messages with the `FAILED` status, received via FastSMSing updates mechanism.
* `DELIVERED` - used for messages with the `DELIVERED` status, received via FastSMSing updates mechanism.


# Application and tests 
 
## Running app locally

1. `go build -o fastSmsingProxy`
2. `./fastSmsingProxy -port=8080`
 
## Running all tests
 
 ```
     go test ./... -v -count=1 -race
 ```

## Linter used:

`golangci-lint run`

