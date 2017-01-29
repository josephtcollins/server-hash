# server-hash
Returns a base64 encoded password using SHA512.

### Install
```
git clone the project
go run server-hash.go
```

### To Use
```
curl -d password=anExamplePassword localhost:8080
```

### For Graceful Shutdown: '/shutdown' 
```
curl localhost:8080/shutdown
```
Server processes any active requests, denies any further requests, and shuts down after completing

