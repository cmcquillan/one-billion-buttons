package main

import (
	"strings"
	"testing"
)

var convertTests = map[string]struct {
	r byte
	g byte
	b byte
}{
	"000000": {
		r: 0, g: 0, b: 0,
	},

	"ffffff": {
		r: 255, g: 255, b: 255,
	},
	"5F2300": {
		r: 95, g: 35, b: 0,
	},
	"#379217": {
		r: 55, g: 146, b: 23,
	},
}

func TestButtonStateFromHex(t *testing.T) {

	for name, test := range convertTests {
		t.Run(name, func(t *testing.T) {
			state := &ButtonState{}
			state.fromHex(name)

			if state.r != test.r || state.g != test.g || state.b != test.b {
				t.Errorf("state (%d, %d, %d) does not match test (%d, %d, %d)",
					state.r, state.g, state.b, test.r, test.g, test.b)
			}
		})
	}
}

func TestButtonStateToHex(t *testing.T) {
	for name, test := range convertTests {
		t.Run(name, func(t *testing.T) {
			state := &ButtonState{
				r: test.r,
				g: test.g,
				b: test.b,
			}

			hex := state.ToHex()

			if hex != strings.ToLower(strings.TrimLeft(name, "#")) {
				t.Errorf("expected %s but received %s for (%d, %d, %d)",
					name, hex, state.r, state.g, state.b)
			}
		})
	}
}
