package main

import (
	"encoding/hex"
	"errors"
	"strconv"
)

const BUTTON_ROWS int64 = 2500
const BUTTON_COLS int64 = 4096
const BUTTONS_PER_PAGE int64 = 100
const BUTTONS_PER_ROW int64 = BUTTON_COLS * BUTTONS_PER_PAGE

type ButtonState struct {
	id int64
	r  byte
	g  byte
	b  byte
}

type ButtonStateDto struct {
	ID  int64  `json:"id"`
	Hex string `json:"hex,omitempty"`
}

type GridPageDto struct {
	X       int64            `json:"x"`
	Y       int64            `json:"y"`
	Buttons []ButtonStateDto `json:"buttons"`
	Next    string           `json:"next"`
}

func HexToBytes(hex string) ([]byte, error) {
	if len(hex) < 6 {
		return nil, errors.New("hex string is insufficient length")
	}

	if hex[0] == '#' {
		hex = hex[1:7]
	}

	r, errR := strconv.ParseInt(hex[0:2], 16, 64)
	g, errG := strconv.ParseInt(hex[2:4], 16, 64)
	b, errB := strconv.ParseInt(hex[4:6], 16, 64)

	if errR != nil || errG != nil || errB != nil {
		return nil, errors.New("could not parse hex code")
	}

	bits := []byte{byte(r), byte(g), byte(b)}
	return bits, nil
}

func (s *ButtonState) fromHex(hex string) error {
	if len(hex) < 6 {
		return errors.New("hex string is insufficient length")
	}

	if hex[0] == '#' {
		hex = hex[1:7]
	}

	r, errR := strconv.ParseInt(hex[0:2], 16, 64)
	g, errG := strconv.ParseInt(hex[2:4], 16, 64)
	b, errB := strconv.ParseInt(hex[4:6], 16, 64)

	if errR != nil || errG != nil || errB != nil {
		return errors.New("could not parse hex code")
	}

	s.r = byte(r)
	s.g = byte(g)
	s.b = byte(b)
	return nil
}

func (s *ButtonState) ToHex() string {
	return ToHex([]byte{s.r, s.g, s.b})
}

func ToHex(rgb []byte) string {
	return hex.EncodeToString(rgb)
}

func (s *ButtonState) IsEmpty() bool {
	return s.r == 0 && s.g == 0 && s.b == 0
}

type GridPage struct {
	X       int64
	Y       int64
	Buttons []ButtonState
}

func HexCodesAreEquivalent(hex1 string, hex2 string) bool {
	bytes1, err1 := HexToBytes(hex1)
	bytes2, err2 := HexToBytes(hex2)

	if err1 != nil || err2 != nil {
		return false
	}

	return bytes1[0] == bytes2[0] &&
		bytes1[1] == bytes2[1] &&
		bytes1[2] == bytes2[2]
}

func ButtonLocationToIndex(x int64, y int64, id int64) (int64, error) {
	offset := ((y - 1) * BUTTON_COLS * BUTTONS_PER_PAGE) + ((x - 1) * BUTTONS_PER_PAGE)
	ix := id - offset

	if ix < 0 || ix > BUTTONS_PER_PAGE-1 {
		return -1, errors.New("asked for button on incorrect page")
	}

	return ix, nil
}

func (s *GridPage) GetButtonById(id int64) *ButtonState {
	ix, err := ButtonLocationToIndex(s.X, s.Y, id)

	if err != nil {
		panic(err)
	}

	return &s.Buttons[ix]
}

func (s *GridPage) EncodeStates() []byte {
	data := make([]byte, 3*BUTTONS_PER_PAGE)

	for i, dx := 0, 0; i < len(s.Buttons); i, dx = i+1, dx+3 {
		data[dx] = s.Buttons[i].r
		data[dx+1] = s.Buttons[i].g
		data[dx+2] = s.Buttons[i].b
	}

	return data
}

func CreateGridPage(x int64, y int64, data []byte) *GridPage {
	buttonState := make([]ButtonState, BUTTONS_PER_PAGE)

	rowIx := x - 1
	colIx := y - 1

	for i := range buttonState {
		id := colIx*BUTTONS_PER_ROW + (rowIx * BUTTONS_PER_PAGE) + int64(i)
		buttonState[i] = ButtonState{
			id: id,
			r:  data[i*3],
			g:  data[(i*3)+1],
			b:  data[(i*3)+2],
		}
	}

	return &GridPage{
		X:       x,
		Y:       y,
		Buttons: buttonState,
	}
}
