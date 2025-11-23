package optional_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/timhugh/optional.go"
)

func TestCreatingOptionals(t *testing.T) {
	t.Run("zero value", func(t *testing.T) {
		var opt optional.Optional[int]
		expectEmpty(t, opt)
	})

	t.Run("Empty", func(t *testing.T) {
		opt := optional.Empty[int]()
		expectEmpty(t, opt)
	})

	t.Run("Of", func(t *testing.T) {
		t.Run("primitives", func(t *testing.T) {
			t.Run("int", func(t *testing.T) {
				t.Run("non-zero", func(t *testing.T) {
					opt := optional.Of(42)
					expectHasValue(t, opt, 42)
				})

				t.Run("zero", func(t *testing.T) {
					opt := optional.Of(0)
					expectHasValue(t, opt, 0)
				})
			})

			t.Run("string", func(t *testing.T) {
				t.Run("non-zero", func(t *testing.T) {
					opt := optional.Of("hello")
					expectHasValue(t, opt, "hello")
				})

				t.Run("zero", func(t *testing.T) {
					opt := optional.Of("")
					expectHasValue(t, opt, "")
				})
			})
		})

		t.Run("pointers", func(t *testing.T) {
			t.Run("non-nil", func(t *testing.T) {
				val := 100
				valPtr := &val
				opt := optional.Of(valPtr)
				expectHasValue(t, opt, valPtr)
			})
			
			t.Run("nil", func(t *testing.T) {
				var val *int = nil
				opt := optional.Of(val)
				expectHasValue(t, opt, val)
			})
		})
	})

	t.Run("OfZeroable", func(t *testing.T) {
		t.Run("primitives", func(t *testing.T) {
			t.Run("int", func(t *testing.T) {
				t.Run("non-zero", func(t *testing.T) {
					opt := optional.OfMaybe(42)
					expectHasValue(t, opt, 42)
				})

				t.Run("zero", func(t *testing.T) {
					opt := optional.OfMaybe(0)
					expectEmpty(t, opt)
				})
			})

			t.Run("string", func(t *testing.T) {
				t.Run("non-zero", func(t *testing.T) {
					opt := optional.OfMaybe("hello")
					expectHasValue(t, opt, "hello")
				})

				t.Run("zero", func(t *testing.T) {
					opt := optional.OfMaybe("")
					expectEmpty(t, opt)
				})
			})
		})

		t.Run("simple structs", func(t *testing.T) {
			type Point struct {
				X, Y int
			}

			t.Run("non-zero", func(t *testing.T) {
				opt := optional.OfMaybe(Point{X: 1, Y: 2})
				expectHasValue(t, opt, Point{X: 1, Y: 2})
			})

			t.Run("zero", func(t *testing.T) {
				opt := optional.OfMaybe(Point{X: 0, Y: 0})
				expectEmpty(t, opt)
			})
		})

		t.Run("slices", func(t *testing.T) {
			t.Run("non-zero", func(t *testing.T) {
				opt := optional.OfMaybeIncomparable([]int{1, 2, 3})
				expectHasValue(t, opt, []int{1, 2, 3})
			})
			
			t.Run("zero", func(t *testing.T) {
				var value []int
				opt := optional.OfMaybeIncomparable(value)
				expectEmpty(t, opt)
			})
		})

		t.Run("pointers", func(t *testing.T) {
			t.Run("non-nil", func(t *testing.T) {
				val := 100
				valPtr := &val
				opt := optional.OfMaybe(valPtr)
				expectHasValue(t, opt, valPtr)
			})
			
			t.Run("nil", func(t *testing.T) {
				var val *int = nil
				opt := optional.OfMaybe(val)
				expectEmpty(t, opt)
			})
		})
	})
}

func TestGet(t *testing.T) {
	optionalWithValue := optional.Of(10)
	emptyOptional := optional.Empty[int]()

	value, err := optionalWithValue.Get()
	expectNoError(t, err)
	expectEqual(t, 10, value)

	_, err = emptyOptional.Get()
	expectError(t, optional.ErrNoValue, err, "getting value from empty optional should return error")
}

func TestConditionals(t *testing.T) {
	t.Run("OrElse", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			value := optional.Of(10).OrElse(20)
			expectEqual(t, 10, value)
		})

		t.Run("empty", func(t *testing.T) {
			value := optional.Empty[int]().OrElse(20)
			expectEqual(t, 20, value)
		})
	})

	t.Run("OrElseGet", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			value := optional.Of(10).OrElseGet(func() int {
				return 20
			})
			expectEqual(t, 10, value)
		})

		t.Run("empty", func(t *testing.T) {
			value := optional.Empty[int]().OrElseGet(func() int {
				return 20
			})
			expectEqual(t, 20, value)
		})
	})

	t.Run("IfPresent", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			called := false
			optional.Of(10).IfPresent(func(v int) {
				called = true
				expectEqual(t, 10, v)
			})
			expectTrue(t, called, "IfPresent should have called the function for optional with value")
		})

		t.Run("empty", func(t *testing.T) {
			optional.Empty[int]().IfPresent(func(v int) {
				t.Error("Function should not have been called for empty optional")
			})
		})
	})

	t.Run("IfPresentOrElse", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			calledIf := false
			optional.Of(10).IfPresentOrElse(
				func(v int) {
					calledIf = true
					expectEqual(t, 10, v)
				},
				func() {
					t.Error("Else function should not have been called for optional with value")
				},
			)
			expectTrue(t, calledIf, "IfPresentOrElse should have called the if function for optional with value")
		})

		t.Run("empty", func(t *testing.T) {
			calledElse := false
			optional.Empty[int]().IfPresentOrElse(
				func(v int) {
					t.Error("If function should not have been called for empty optional")
				},
				func() {
					calledElse = true
				},
			)
			expectTrue(t, calledElse, "IfPresentOrElse should have called the else function for empty optional")
		})
	})

	t.Run("IfEmpty", func(t *testing.T) {
		t.Run("with value", func(t *testing.T) {
			optional.Of(10).IfEmpty(func() {
				t.Error("Function should not have been called for optional with value")
			})
		})

		t.Run("empty", func(t *testing.T) {
			called := false
			optional.Empty[int]().IfEmpty(func() {
				called = true
			})
			expectTrue(t, called, "IfEmpty should have called the function for empty optional")
		})
	})

	t.Run("Map", func(t *testing.T) {
		mapFunc := func(v int) string {
			return "Value is " + fmt.Sprint(v)
		}

		t.Run("with value", func(t *testing.T) {
			mappedOpt := optional.Map(optional.Of(10), mapFunc)
			expectHasValue(t, mappedOpt, "Value is 10")
		})
		
		t.Run("empty", func(t *testing.T) {
			mappedOpt := optional.Map(optional.Empty[int](), mapFunc)
			expectEmpty(t, mappedOpt)
		})
	})
}

func TestMapExampleInReadme(t *testing.T) {
	var ErrReallyBadError = fmt.Errorf("something really bad happened")

	doSomethingSketchy := func() error {
		return fmt.Errorf("disk is full")
	}
	err := doSomethingSketchy()

	var receivedError error
	optional.Map(optional.Of(err), func(err error) error {
		return fmt.Errorf("%w: %s", ErrReallyBadError, err)
	}).IfPresent(func(err error) {
		receivedError = err
	})

	expectTrue(t, errors.Is(receivedError, ErrReallyBadError), "received error should wrap ErrReallyBadError")
	expectEqual(t, "something really bad happened: disk is full", receivedError.Error())
}

func BenchmarkOfMaybe(b *testing.B) {
	for b.Loop() {
		optional.OfMaybe(42)
	}
}

func BenchmarkOfMaybeIncomparable(b *testing.B) {
	for b.Loop() {
		optional.OfMaybeIncomparable(42)
	}
}
