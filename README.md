# Go Kafka Sample Producer

Just a simple project to test some kafka features. Produce a message for any random topic 


### used frameworks
[1 Gonic] (https://github.com/gin-gonic/gin)

[2 franz go] (https://github.com/twmb/franz-go)

### how to use

1. Start local with

```
go run main.go
```

2. Sample post on `http://localhost:8091/send/[foo-topic]`

```json
{
  "foo": "bar"
}
```

This post will send it to a local kafka on port `29092`.
Feel free to change it in `app/producer.go`. The project also support the kafka on the azure eventhub
