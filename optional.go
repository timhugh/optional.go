// Package optional provides a generic Optional type to represent values that may or may not be present
// and functions to interact with them.
package optional

import (
	"fmt"
	"reflect"
)

// ErrNoValue is returned when attempting to get a value from an empty Optional.
var ErrNoValue = fmt.Errorf("Optional has no value")

// Optional is a generic immutable type that may or may not contain a value of type T.
type Optional[T any] struct {
	value    T
	hasValue bool
}

// Of creates an Optional with a guaranteed value.
func Of[T any](value T) Optional[T] {
	return Optional[T]{
		value:    value,
		hasValue: true,
	}
}

// OfMaybe creates an Optional from a value of type T. If the value is equal to the zero value of T,
// the Optional will be empty; otherwise, it will contain the value. This requires T to be comparable.
func OfMaybe[T comparable](value T) Optional[T] {
	var hasValue bool
	if value != *new(T) {
		hasValue = true
	}
	return Optional[T]{
		value:    value,
		hasValue: hasValue,
	}
}

func OfMaybeIncomparable[T any](value T) Optional[T] {
	var hasValue bool
	if !reflect.DeepEqual(value, *new(T)) {
		hasValue = true
	}
	return Optional[T]{
		value:    value,
		hasValue: hasValue,
	}
}

// Empty creates an Optional guaranteed to have no value.
func Empty[T any]() Optional[T] {
	return Optional[T]{
		hasValue: false,
	}
}

// Get returns the contained value of type T if present, otherwise returns an ErrNoValue error.
func (o Optional[T]) Get() (T, error) {
	if o.Empty() {
		return o.value, ErrNoValue
	}
	return o.value, nil
}

// HasValue returns true if the Optional contains a value.
func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

// Empty returns true if the Optional does not contain a value.
func (o Optional[T]) Empty() bool {
	return !o.HasValue()
}

// OrElse returns the contained value of type T if present, otherwise returns the provided
// elseValue, which must be of type T.
func (o Optional[T]) OrElse(elseValue T) T {
	if o.hasValue {
		return o.value
	}
	return elseValue
}

// OrElseGet returns the contained value of type T if present, otherwise invokes the provided
// elseFunc to obtain a value of type T.
func (o Optional[T]) OrElseGet(elseFunc func() T) T {
	if o.hasValue {
		return o.value
	}
	return elseFunc()
}

// IfPresent invokes the provided action function with the contained value of type T if present.
func (o Optional[T]) IfPresent(action func(T)) {
	if o.hasValue {
		action(o.value)
	}
}

// IfPresentOrElse invokes the provided action function with the contained value of type T if present,
// otherwise invokes the elseAction function.
func (o Optional[T]) IfPresentOrElse(action func(T), elseAction func()) {
	if o.hasValue {
		action(o.value)
	} else {
		elseAction()
	}
}

// IfEmpty invokes the provided action function if the Optional is empty.
func (o Optional[T]) IfEmpty(action func()) {
	if !o.hasValue {
		action()
	}
}

// Map applies the provided mapper function to the contained value of type T if present,
// returning a new Optional containing the result of type U. If the Optional is empty,
// an empty Optional of type U is returned.
func Map[T any, U any](opt Optional[T], mapper func(T) U) Optional[U] {
	if opt.HasValue() {
		return Of(mapper(opt.value))
	}
	return Empty[U]()
}
