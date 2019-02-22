History
-----------
History embodies a simple, easy structured logger with the focus around a simple API with hierarchical logs.

## Install

```bash
go get -v github.com/gokit/history
```


## Example


```go

// create a history source.
src := history.WithHandlers(std.Std)

// create a ctx for a single log session.
ctx := src.FromTags("core", "wombat")

// log a status for the context to provide
// progress status.
ctx.Info("Waiting to count")
ctx.Yellow("Waiting to count %d", 10)
ctx.Error(ErrNoClient,"Waiting to count")

// Add fields to the context.
ctx.With("k", "bob")

// Add a collection with key 'l.v' to the context.
ctx.Collect("l.v", 1, 32, 3)

// Add more items to 'l.v' collection
ctx.Collect("l.v", 4, 5, 10)

// Have context finalzied and sent to handlers.
ctx.Done()

```

