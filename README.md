# optional.go

Because this should probably already exist 👀

No dependencies, no magic, just a value wrapper with some convenience methods.

## Usage

```go
import "github.com/timhugh/optional.go"
```

Go's zero values haven't always sat right with me. It's fine some of the time, but sometimes there's an important difference between a zero value and an empty value.

If you're familiar with Java or C++ optionals, this will probably look familiar.

```go
// can be explicitly constructed from known values:
notMessage := optional.Empty[string]
notMessage.Empty() // => true
notMessage.HasValue() // => false
message, err := notMessage.Get() // => err: "Optional has no value"

definitelyMessage := optional.Of("Hello!")
definitelyMessage.Empty() // => false
definitelyMessage.HasValue() // => true
message, err = definitelyMessage.Get() // => message: "Hello!"

// can also infer presence based on the zero value of the type:
message = "Hola!"
maybeMessage = optional.OfMaybe(message)
maybeMessage.HasValue() // => true

message = ""
maybeMessage = optional.OfMaybe(message)
maybeMessage.HasValue() // => false

// supports some java-inspired chaining methods:
maybeMessage.OrElse("Goodbye!") // => "Goodbye!"
maybeMessage.OrElseGet(func() string {
  return "Goodbye!"
}) // => "Goodbye!"
maybeMessage.IfPresent(func(v string) {
  fmt.Sprintf("Here's a message for you! %s", v)
})

// and some simple monad operations:
var ErrReallyBadError = fmt.Errorf("something really bad happened")
err := doSomethingSketchy()
optional.Map(optional.Of(err), func(err error) error {
  return fmt.Errorf("%w: %s", ErrReallyBadError, err)
}).IfPresent(func(err error) {
  panic(err)
})
```

For more examples, check out [the tests](optional_test.go).
