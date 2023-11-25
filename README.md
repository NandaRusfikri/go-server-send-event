
# Server Sent Event Golang

How to sent message from server to client ?
you can use alternative SSE (Server Sent Event)



![Logo](https://res.cloudinary.com/practicaldev/image/fetch/s--k4GZeQBW--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://cdn-images-1.medium.com/max/2000/1%2AUg-BosrJefTOOEBmtA2H4Q.jpeg)


## Usage/Examples

Run Program

```shell
go run main.go
```

Example Send Event with trigger http call

POST http://localhost:4444/api/sse

```json
{
    "message" :"NandaRusfikri"
}
```

