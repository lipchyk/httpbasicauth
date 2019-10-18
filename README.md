# httpbasicauth
[![CircleCI](https://circleci.com/gh/yspro/httpbasicauth.svg?style=svg)](https://circleci.com/gh/yspro/httpbasicauth)
[![codecov](https://codecov.io/gh/yspro/httpbasicauth/branch/master/graph/badge.svg)](https://codecov.io/gh/yspro/httpbasicauth)

An HTTP Basic Auth middleware for Go

## Usage

```go
import (
    "net/http"
    "github.com/yspro/httpbasicauth"
)

// credentials
creds := map[string]string{"u$eR": "$ecret"}
middleware := httpbasicauth.Handle(creds, "Restricted Zone")

yourhandler := http.HandlerFunc(
    func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "Hello World")
    },
)

http.Handle("/secret", middleware(yourhandler))
err := http.ListenAndServe(":8080", nil)
if err != nil {
    panic(err)
}
```
