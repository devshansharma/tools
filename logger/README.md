# logger 

logger package to use in rest based apis, taking reference from various online sources.

## How to use
```
package main

import (
	"context"
	"log/slog"

	"github.com/devshansharma/tools/logger"
	"github.com/google/uuid"
	"github.com/mdobak/go-xerrors"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("name", u.FirstName+" "+u.LastName),
	)
}

type ContextKey string

var (
	requestIDCtx  ContextKey = "requestID"
	customerIDCtx ContextKey = "customerID"
)

func recorder(ctx context.Context, rec slog.Record) (slog.Record, error) {
	requestID, ok := ctx.Value(requestIDCtx).(string)
	if ok {
		rec.Add("requestID", requestID)
	}

	customerID, ok := ctx.Value(customerIDCtx).(string)
	if ok {
		rec.Add("customerID", customerID)
	}

	return rec, nil
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, requestIDCtx, uuid.NewString())
	ctx = context.WithValue(ctx, customerIDCtx, uuid.NewString())

	log := logger.New(
		logger.WithJSON(true),
		logger.WithHandle(recorder),
		logger.WithSource(true),
		logger.WithReplaceAttr(logger.WithShortFileNameAndErrorTrace),
	)
	slog.SetDefault(log)

	u := &User{
		ID:        "user-12234",
		FirstName: "Jan",
		LastName:  "Doe",
		Email:     "jan@example.com",
		Password:  "pass-12334",
	}

	err := xerrors.New("something happened")

	slog.ErrorContext(ctx, "new error message", "user", u, slog.Any("error", err))

}

```

The result will be:
```
{
  "time": "2024-05-18T12:37:47.892882993+05:30",
  "level": "ERROR",
  "source": {
    "function": "main.main",
    "file": "main.go",
    "line": 71
  },
  "msg": "new error message",
  "user": {
    "id": "user-12234",
    "name": "Jan Doe"
  },
  "error": {
    "msg": "something happened",
    "trace": [
      {
        "func": "main.main",
        "source": "tools/main.go",
        "line": 69
      },
      {
        "func": "runtime.main",
        "source": "runtime/proc.go",
        "line": 271
      },
      {
        "func": "runtime.goexit",
        "source": "runtime/asm_amd64.s",
        "line": 1695
      }
    ]
  },
  "requestID": "e6bd5c5b-896d-4933-995a-27bdc5dc2298",
  "customerID": "1a2ab12e-fa3e-4538-9896-2b070925b029"
}
```