# libraryfuzzing

Run the follow command to perform modification on a provided testcase:

```go
go run . -s
```

To modify a existing package, e.g. `websocket`, run:

```go
go run . -s -f TestHandshake websocket
```
