# walk
--
    import "vimagination.zapto.org/javascript/walk"

Package walk provides a javascript type walker

## Usage

#### func  Walk

```go
func Walk(t javascript.Type, fn Handler) error
```
Walk calls the Handle function on the given interface for each non-nil,
non-Token field of the given javascript type.

#### type Handler

```go
type Handler interface {
	Handle(javascript.Type) error
}
```

Handler is used to process javascript types.

#### type HandlerFunc

```go
type HandlerFunc func(javascript.Type) error
```

HandlerFunc wraps a func to implement Handler interface.

#### func (HandlerFunc) Handle

```go
func (h HandlerFunc) Handle(t javascript.Type) error
```
Handle implements the Handler interface.
