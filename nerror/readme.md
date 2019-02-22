Errors
---------

Errors provides a package for wrapping go based nerror with location of return site.


## Install

```bash
go get -v github.com/gokit/npkg/nerror
```

## Usage

1. Create new nerror


```go
newBadErr = nerror.New("failed connection: %s", "10.9.1.0")
```

2. Create new nerror with stacktrace


```go
newBadErr = nerror.Stacked("failed connection: %s", "10.9.1.0")
```

3. Wrap existing error


```go
newBadErr = nerror.Wrap(BadErr, "something bad happened here")
```

4. Wrap existing error with stacktrace


```go
newBadErr = nerror.WrapStack(BadErr, "something bad happened here")
```

5. Add Stack to package's error type without stack trace.


```go
newBadErr = nerror.StackIt(BadErr)
```
