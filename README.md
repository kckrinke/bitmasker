# Bitmasker [![GoDoc](https://godoc.org/github.com/kckrinke/bitmasker?status.png)](https://godoc.org/github.com/kckrinke/bitmasker)

Bitmasker is a tool used to automate the creation of helper methods when
dealing with bitmask-type constant flags. Given the name of an unsigned
integer type T that has constants defined, bitmasker will create a new
self-contained Go source file implementing the BitMask and fmt.Stringer
interfaces.

## Getting Started

To get up and running, follow the normal go install procedure:

```
go install github.com/kckrinke/bitmasker
```

## Example Usage

Bitmasker is intended to be used with go:generate but can operated standalone
as well. For example:

Given this snippet:

```
package mental

//go:generate bitmasker -type=State
type State uint

const (
	Unconscious State = 0
  Conscious State = (1 << iota)
  Meditative
  Distracted
  Entertained = Distracted
)
```

Standalone usage:
```
bitmasker -type=State
```

Using go-generate
```
go generate
```

In both cases a new file named "state_bitmask.go" will be created with the
following contents:

```
// Code generated by "bitmasker -type=State"; DO NOT EDIT.

package mental

import "strconv"

type IState interface {
	Has(m State) bool
	Set(m State) State
	Clear(m State) State
	Toggle(m State) State
	String() string
}

func (i State) Has(m State) bool {
	return i&m != 0
}

func (i State) Set(m State) State {
	return i | m
}

func (i State) Clear(m State) State {
	return i &^ m
}

func (i State) Toggle(m State) State {
	return i ^ m
}

const (
	_State_name_0 = "Unconscious"
	_State_name_1 = "Conscious"
	_State_name_2 = "Meditative"
	_State_name_3 = "Distracted"
)

func (i State) String() string {
	switch {
	case i == 0:
		return _State_name_0
	case i == 2:
		return _State_name_1
	case i == 4:
		return _State_name_2
	case i == 8:
		return _State_name_3
	default:
		return "State(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
```

## Running the tests

Unit tests are provided and can be invoked using the normal Go pattern:

```
go test
```

## Authors

* **Kevin C. Krinke** - *Bitmasker author* - [kckrinke](https://github.com/kckrinke)
* **The Go Authors** - *Stringer derived sources* - [stringer](https://golang.org/x/tools/cmd/stringer)

## License

This project is licensed under the LGPL (specifically v3) - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

* Thanks to Golang.org for the [stringer](https://golang.org/x/tools/cmd/stringer) program that bitmasker sources are derived from.
