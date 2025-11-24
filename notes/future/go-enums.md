# Enums in Go

## Use stringer tool

`go install golang.org/x/tools/cmd/stringer@latest`

```go
//go:generate stringer -type=Color
type Color int
```

## implement marshalling method (interface)

```go
func (c Color) MarshalText() ([]byte, error) {
  return []byte(c.String()), nil
}
// ...
bts, _ := json.Marshal(MyResponse{Color: Blue})
fmt.Println(string(bts)) // {"Color":"Blue"}

```

## Unmarshalling

```go
var ErrInvalidColor = errors.New("invalid color")

func ParseColor(in string) (Color, error) {
  switch in {
    case Red.String():
      return Red, nil
    case Green.String():
      return Green, nil
    case Blue.String():
      return Blue, nil
  }

  return Red, fmt.Errorf("%q is not a valid color: %w", in, ErrInvalidColor)
}

func (c *Color) UnmarshalText(text []byte) error {
  parsed, err := ParseColor(string(text))
  if err != nil {
    return err
  }

  *c = parsed
  return nil
}
```
