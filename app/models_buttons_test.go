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

func TestButtonLocationToIndex(t *testing.T) {
	tests := []struct {
		name        string
		x           int64
		y           int64
		id          int64
		expectedIx  int64
		expectError bool
	}{
		{
			name:        "valid button at (1,1) first button",
			x:           1,
			y:           1,
			id:          0,
			expectedIx:  0,
			expectError: false,
		},
		{
			name:        "valid button at (1,1) last button",
			x:           1,
			y:           1,
			id:          99,
			expectedIx:  99,
			expectError: false,
		},
		{
			name:        "valid button at (2,1)",
			x:           2,
			y:           1,
			id:          100,
			expectedIx:  0,
			expectError: false,
		},
		{
			name:        "valid button at (1,2)",
			x:           1,
			y:           2,
			id:          BUTTON_COLS * BUTTONS_PER_PAGE,
			expectedIx:  0,
			expectError: false,
		},
		{
			name:        "button id too low for page",
			x:           2,
			y:           1,
			id:          50,
			expectedIx:  -1,
			expectError: true,
		},
		{
			name:        "button id too high for page",
			x:           1,
			y:           1,
			id:          150,
			expectedIx:  -1,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ix, err := ButtonLocationToIndex(test.x, test.y, test.id)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if ix != -1 {
					t.Errorf("expected index -1 but got %d", ix)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if ix != test.expectedIx {
					t.Errorf("expected index %d but got %d", test.expectedIx, ix)
				}
			}
		})
	}
}

func TestHexCodesAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		hex1     string
		hex2     string
		expected bool
	}{
		{
			name:     "identical hex codes",
			hex1:     "ffffff",
			hex2:     "ffffff",
			expected: true,
		},
		{
			name:     "identical hex codes with hash",
			hex1:     "#ffffff",
			hex2:     "#ffffff",
			expected: true,
		},
		{
			name:     "mixed hash and no hash",
			hex1:     "#ffffff",
			hex2:     "ffffff",
			expected: true,
		},
		{
			name:     "case insensitive",
			hex1:     "FFFFFF",
			hex2:     "ffffff",
			expected: true,
		},
		{
			name:     "different colors",
			hex1:     "000000",
			hex2:     "ffffff",
			expected: false,
		},
		{
			name:     "invalid first hex",
			hex1:     "gggggg",
			hex2:     "ffffff",
			expected: false,
		},
		{
			name:     "invalid second hex",
			hex1:     "ffffff",
			hex2:     "gggggg",
			expected: false,
		},
		{
			name:     "both invalid hex",
			hex1:     "gggggg",
			hex2:     "hhhhhh",
			expected: false,
		},
		{
			name:     "too short hex",
			hex1:     "fff",
			hex2:     "ffffff",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := HexCodesAreEquivalent(test.hex1, test.hex2)
			if result != test.expected {
				t.Errorf("HexCodesAreEquivalent(%q, %q) = %v; expected %v",
					test.hex1, test.hex2, result, test.expected)
			}
		})
	}
}

func TestHexToBytes(t *testing.T) {
	tests := []struct {
		name      string
		hex       string
		expected  []byte
		expectErr bool
	}{
		{
			name:      "valid hex without hash",
			hex:       "ffffff",
			expected:  []byte{255, 255, 255},
			expectErr: false,
		},
		{
			name:      "valid hex with hash",
			hex:       "#ffffff",
			expected:  []byte{255, 255, 255},
			expectErr: false,
		},
		{
			name:      "valid hex black",
			hex:       "000000",
			expected:  []byte{0, 0, 0},
			expectErr: false,
		},
		{
			name:      "valid hex mixed case",
			hex:       "5F2300",
			expected:  []byte{95, 35, 0},
			expectErr: false,
		},
		{
			name:      "too short hex",
			hex:       "fff",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "invalid hex characters",
			hex:       "gggggg",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "empty string",
			hex:       "",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "long hex string should work",
			hex:       "#ffaabb00",
			expected:  []byte{255, 170, 187},
			expectErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := HexToBytes(test.hex)

			if test.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != len(test.expected) {
					t.Errorf("expected length %d, got %d", len(test.expected), len(result))
				}
				for i := 0; i < len(test.expected); i++ {
					if result[i] != test.expected[i] {
						t.Errorf("expected byte[%d] = %d, got %d", i, test.expected[i], result[i])
					}
				}
			}
		})
	}
}

func TestCreateGridPage(t *testing.T) {
	tests := []struct {
		name     string
		x        int64
		y        int64
		data     []byte
		expected GridPage
	}{
		{
			name: "create page at (1,1)",
			x:    1,
			y:    1,
			data: make([]byte, 300), // 100 buttons * 3 bytes each = 300 bytes
			expected: GridPage{
				X: 1,
				Y: 1,
				Buttons: func() []ButtonState {
					buttons := make([]ButtonState, 100)
					for i := 0; i < 100; i++ {
						buttons[i] = ButtonState{id: int64(i), r: 0, g: 0, b: 0}
					}
					return buttons
				}(),
			},
		},
		{
			name: "create page at (2,3) with some color data",
			x:    2,
			y:    3,
			data: func() []byte {
				data := make([]byte, 300)
				// Set first button to red (255, 0, 0)
				data[0] = 255
				data[1] = 0
				data[2] = 0
				// Set second button to green (0, 255, 0)
				data[3] = 0
				data[4] = 255
				data[5] = 0
				return data
			}(),
			expected: GridPage{
				X: 2,
				Y: 3,
				Buttons: func() []ButtonState {
					buttons := make([]ButtonState, 100)
					// Calculate base ID for page (2,3)
					baseId := (3-1)*BUTTON_COLS*BUTTONS_PER_PAGE + (2-1)*BUTTONS_PER_PAGE
					buttons[0] = ButtonState{id: baseId, r: 255, g: 0, b: 0}
					buttons[1] = ButtonState{id: baseId + 1, r: 0, g: 255, b: 0}
					for i := 2; i < 100; i++ {
						buttons[i] = ButtonState{id: baseId + int64(i), r: 0, g: 0, b: 0}
					}
					return buttons
				}(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := CreateGridPage(test.x, test.y, test.data)

			if result.X != test.expected.X {
				t.Errorf("expected X=%d, got X=%d", test.expected.X, result.X)
			}
			if result.Y != test.expected.Y {
				t.Errorf("expected Y=%d, got Y=%d", test.expected.Y, result.Y)
			}
			if len(result.Buttons) != len(test.expected.Buttons) {
				t.Errorf("expected %d buttons, got %d", len(test.expected.Buttons), len(result.Buttons))
			}

			// Test first few buttons for correctness
			for i := 0; i < min(len(result.Buttons), 5); i++ {
				if result.Buttons[i].id != test.expected.Buttons[i].id {
					t.Errorf("button[%d]: expected id=%d, got id=%d", i, test.expected.Buttons[i].id, result.Buttons[i].id)
				}
				if result.Buttons[i].r != test.expected.Buttons[i].r {
					t.Errorf("button[%d]: expected r=%d, got r=%d", i, test.expected.Buttons[i].r, result.Buttons[i].r)
				}
				if result.Buttons[i].g != test.expected.Buttons[i].g {
					t.Errorf("button[%d]: expected g=%d, got g=%d", i, test.expected.Buttons[i].g, result.Buttons[i].g)
				}
				if result.Buttons[i].b != test.expected.Buttons[i].b {
					t.Errorf("button[%d]: expected b=%d, got b=%d", i, test.expected.Buttons[i].b, result.Buttons[i].b)
				}
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestGridPageGetButtonById(t *testing.T) {
	// Create a test grid page at (1,1)
	data := make([]byte, 300)
	data[0] = 255  // First button red
	data[1] = 0
	data[2] = 0
	data[9] = 0    // Fourth button green
	data[10] = 255
	data[11] = 0
	
	page := CreateGridPage(1, 1, data)

	tests := []struct {
		name        string
		id          int64
		expectedBtn ButtonState
		shouldPanic bool
	}{
		{
			name: "get first button",
			id:   0,
			expectedBtn: ButtonState{
				id: 0,
				r:  255,
				g:  0,
				b:  0,
			},
			shouldPanic: false,
		},
		{
			name: "get fourth button",
			id:   3,
			expectedBtn: ButtonState{
				id: 3,
				r:  0,
				g:  255,
				b:  0,
			},
			shouldPanic: false,
		},
		{
			name:        "invalid button id should panic",
			id:          150,
			shouldPanic: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("expected panic but none occurred")
					}
				}()
			}

			result := page.GetButtonById(test.id)

			if !test.shouldPanic {
				if result.id != test.expectedBtn.id {
					t.Errorf("expected id=%d, got id=%d", test.expectedBtn.id, result.id)
				}
				if result.r != test.expectedBtn.r {
					t.Errorf("expected r=%d, got r=%d", test.expectedBtn.r, result.r)
				}
				if result.g != test.expectedBtn.g {
					t.Errorf("expected g=%d, got g=%d", test.expectedBtn.g, result.g)
				}
				if result.b != test.expectedBtn.b {
					t.Errorf("expected b=%d, got b=%d", test.expectedBtn.b, result.b)
				}
			}
		})
	}
}

func TestGridPageEncodeStates(t *testing.T) {
	tests := []struct {
		name     string
		page     *GridPage
		expected []byte
	}{
		{
			name: "encode empty buttons",
			page: &GridPage{
				X: 1,
				Y: 1,
				Buttons: func() []ButtonState {
					buttons := make([]ButtonState, 100)
					for i := 0; i < 100; i++ {
						buttons[i] = ButtonState{id: int64(i), r: 0, g: 0, b: 0}
					}
					return buttons
				}(),
			},
			expected: make([]byte, 300), // All zeros
		},
		{
			name: "encode with colored buttons",
			page: &GridPage{
				X: 1,
				Y: 1,
				Buttons: func() []ButtonState {
					buttons := make([]ButtonState, 100)
					for i := 0; i < 100; i++ {
						buttons[i] = ButtonState{id: int64(i), r: 0, g: 0, b: 0}
					}
					// Set first button to red
					buttons[0] = ButtonState{id: 0, r: 255, g: 0, b: 0}
					// Set second button to green
					buttons[1] = ButtonState{id: 1, r: 0, g: 255, b: 0}
					return buttons
				}(),
			},
			expected: func() []byte {
				data := make([]byte, 300)
				data[0] = 255  // First button red
				data[1] = 0
				data[2] = 0
				data[3] = 0    // Second button green
				data[4] = 255
				data[5] = 0
				return data
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.page.EncodeStates()

			if len(result) != len(test.expected) {
				t.Errorf("expected length %d, got %d", len(test.expected), len(result))
			}

			for i := 0; i < len(test.expected); i++ {
				if result[i] != test.expected[i] {
					t.Errorf("byte[%d]: expected %d, got %d", i, test.expected[i], result[i])
				}
			}
		})
	}
}

func TestButtonStateIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		state    ButtonState
		expected bool
	}{
		{
			name:     "empty black button",
			state:    ButtonState{r: 0, g: 0, b: 0},
			expected: true,
		},
		{
			name:     "red button",
			state:    ButtonState{r: 255, g: 0, b: 0},
			expected: false,
		},
		{
			name:     "green button",
			state:    ButtonState{r: 0, g: 255, b: 0},
			expected: false,
		},
		{
			name:     "blue button",
			state:    ButtonState{r: 0, g: 0, b: 255},
			expected: false,
		},
		{
			name:     "white button",
			state:    ButtonState{r: 255, g: 255, b: 255},
			expected: false,
		},
		{
			name:     "very dark but not empty",
			state:    ButtonState{r: 1, g: 0, b: 0},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.state.IsEmpty()
			if result != test.expected {
				t.Errorf("IsEmpty() = %v; expected %v", result, test.expected)
			}
		})
	}
}
