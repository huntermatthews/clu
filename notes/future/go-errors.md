# Error Handling Suggestions

## In errors.go

```go
func WrapWithCaller(err error) error {
  pc, _, line, ok := runtime.Caller(1)
  if !ok {
    return fmt.Errorf("%v: %w", "unknown", err)
  }
  fn := runtime.FuncForPC(pc).Name()
  pkgFunc := path.Base(fn)
  return fmt.Errorf("%v:%d: %w", pkgFunc, line, err)
}
```

## Usage

```go
  func someFunction() error {
    err := someStruct{}.someMethod()
    return WrapWithCaller(err)
  }

  type someStruct struct {
  }

  func (ss someStruct) someMethod() error {
    return WrapWithCaller(errors.New("boom"))
  }
```

## or, do it by hand

```go
  func someFunction() error {
    err := someStruct{}.someMethod()
    return fmt.Errorf("someFunction: %w", err)
  }

  type someStruct struct {
  }

  func (ss someStruct) someMethod() error {
    return fmt.Errorf("someMethod: %w", errors.New("boom"))
  }
```
