package color

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtilFuncs(t *testing.T) {
	is := assert.New(t)

	// IsConsole
	is.True(IsConsole(os.Stdin))
	is.True(IsConsole(os.Stdout))
	is.True(IsConsole(os.Stderr))
	is.False(IsConsole(&bytes.Buffer{}))
	ff, err := os.OpenFile(".travis.yml", os.O_WRONLY, 0)
	is.NoError(err)
	is.False(IsConsole(ff))

	// IsMSys
	oldVal := os.Getenv("MSYSTEM")
	is.NoError(os.Setenv("MSYSTEM", "MINGW64"))
	is.True(IsMSys())
	is.NoError(os.Unsetenv("MSYSTEM"))
	is.False(IsMSys())
	_ = os.Setenv("MSYSTEM", oldVal)

	// IsSupport256Color
	oldVal = os.Getenv("TERM")
	_ = os.Unsetenv("TERM")
	is.False(IsSupportColor())
	is.False(IsSupport256Color())

	// ConEmuANSI
	mockEnvValue("ConEmuANSI", "ON", func(_ string) {
		is.True(IsSupportColor())
	})

	// ANSICON
	mockEnvValue("ANSICON", "189x2000 (189x43)", func(_ string) {
		is.True(IsSupportColor())
	})

	// "COLORTERM=truecolor"
	mockEnvValue("COLORTERM", "truecolor", func(_ string) {
		is.True(IsSupportTrueColor())
	})

	// TERM
	mockEnvValue("TERM", "screen-256color", func(_ string) {
		is.True(IsSupportColor())
	})

	// TERM
	mockEnvValue("TERM", "tmux-256color", func(_ string) {
		is.True(IsSupportColor())
	})

	// TERM
	mockEnvValue("TERM", "rxvt-unicode-256color", func(_ string) {
		is.True(IsSupportColor())
	})

	is.NoError(os.Setenv("TERM", "xterm-vt220"))
	is.True(IsSupportColor())
	// revert
	if oldVal != "" {
		is.NoError(os.Setenv("TERM", oldVal))
	} else {
		is.NoError(os.Unsetenv("TERM"))
	}
}

func TestRgbTo256Table(t *testing.T) {
	index := 0
	for hex, c256 := range RgbTo256Table() {
		Hex(hex).Print("RGB:", hex)
		fmt.Print(" = ")
		C256(c256).Print("C256:", c256)
		fmt.Print(" | ")
		index++
		if index%5 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestC256ToRgbV1(t *testing.T) {
	for i :=0; i < 256; i++ {
		c256 := uint8(i)
		C256(c256).Printf("C256:%d", c256)
		fmt.Print(" => ")
		rgb := C256ToRgbV1(c256)
		RGBFromSlice(rgb).Printf("RGB:%v | ", rgb)
		// assert.Equal(t, item.want, rgb, fmt.Sprint("256 code:", c256))
		if i%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestC256ToRgb(t *testing.T) {
	for i :=0; i < 256; i++ {
		c256 := uint8(i)
		C256(c256).Printf("C256:%d", c256)
		fmt.Print(" => ")
		rgb := C256ToRgb(c256)
		RGBFromSlice(rgb).Printf("RGB:%v | ", rgb)
		// assert.Equal(t, item.want, rgb, fmt.Sprint("256 code:", c256))
		if i%4 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestHexToRgb(t *testing.T) {
	tests := []struct {
		given string
		want  []int
	}{
		{"666", []int{102, 102, 102}},
		{"ccc", []int{204, 204, 204}},
		{"#abc", []int{170, 187, 204}},
		{"#aa99cd", []int{170, 153, 205}},
	}

	for _, item := range tests {
		assert.Equal(t, HexToRgb(item.given), item.want)
		assert.Equal(t, HexToRGB(item.given), item.want)
		assert.Equal(t, Hex2rgb(item.given), item.want)
	}

	assert.Len(t, HexToRgb(""), 0)
	assert.Len(t, HexToRgb("13"), 0)
}

func TestRgbToHex(t *testing.T) {
	tests := []struct {
		want  string
		given []int
	}{
		{"666666", []int{102, 102, 102}},
		{"cccccc", []int{204, 204, 204}},
		{"aabbcc", []int{170, 187, 204}},
		{"aa99cd", []int{170, 153, 205}},
	}

	for _, item := range tests {
		assert.Equal(t, RgbToHex(item.given), item.want)
		assert.Equal(t, Rgb2hex(item.given), item.want)
	}
}

func TestRgbToAnsi(t *testing.T) {
	tests := []struct {
		want uint8
		rgb  []uint8
		isBg bool
	}{
		{40, []uint8{102, 102, 102}, true},
		{37, []uint8{204, 204, 204}, false},
		{47, []uint8{170, 78, 204}, true},
		{37, []uint8{170, 153, 245}, false},
		{30, []uint8{127, 127, 127}, false},
		{40, []uint8{127, 127, 127}, true},
		{90, []uint8{128, 128, 128}, false},
		{97, []uint8{34, 56, 255}, false},
		{31, []uint8{134, 56, 56}, false},
		{30, []uint8{0, 0, 0}, false},
		{40, []uint8{0, 0, 0}, true},
		{97, []uint8{255, 255, 255}, false},
		{107, []uint8{255, 255, 255}, true},
	}

	for _, item := range tests {
		r, g, b := item.rgb[0], item.rgb[1], item.rgb[2]

		assert.Equal(
			t,
			item.want,
			RgbToAnsi(r, g, b, item.isBg),
			fmt.Sprint("rgb=", item.rgb, ", is bg? ", item.isBg),
		)
		assert.Equal(t, item.want, Rgb2ansi(r, g, b, item.isBg))
	}
}
