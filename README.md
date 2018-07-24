# redislike

## Deployment

`$> go install`

`$> $GOPATH/bin/redislike`

After it you will have running server on 8080 port.

## API methods

* /keys - to get all available keys

  `curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/keys`

* /remove/:key - to remove specific key from cache

  `curl -i -w "\n" -X DELETE -H 'Content-Type: application/json' localhost:8080/remove/abc`

* /ttl/:key - to set TTL for existing key

  `curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"value":10}' localhost:8080/ttl/abc`

* /set - to set String value of a key

  `curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"abc", "value":"test value"}' localhost:8080/set`

* /get/:key - to get String value of a key

  `curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/get/abc`

* /push - to append values to a list

  `curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"list", "value":["v1","v2"]}' localhost:8080/push`

* /pop/:key - to remove and get first item in a list

  `curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/pop/list`

* /hset - to set the string value of a dict field

  `curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"dict", "field":"f1", "value": "v1"}' localhost:8080/hset`

* /hget/:key/:field - to get value of a dict field

  `curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/hget/dict/f1`
