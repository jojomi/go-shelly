package shelly

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
)

type ModeColorOptions struct {
	Brightness *int
	Transition *int
	Gain       *int
	Color      *colorful.Color
	Error      error
}

func NewModeColorOptions() *ModeColorOptions {
	return &ModeColorOptions{}
}

func (x *ModeColorOptions) SetBrightness(brightness int) *ModeColorOptions {
	if brightness < 0 {
		x.Error = fmt.Errorf("invalid brightness < 0: %d", brightness)
		return x
	}
	if brightness > 100 {
		x.Error = fmt.Errorf("invalid brightness > 100: %d", brightness)
		return x
	}

	x.Brightness = &brightness
	return x
}

func (x *ModeColorOptions) SetGain(gain int) *ModeColorOptions {
	if gain < 0 {
		x.Error = fmt.Errorf("invalid brightness < 0: %d", gain)
		return x
	}
	if gain > 100 {
		x.Error = fmt.Errorf("invalid brightness > 100: %d", gain)
		return x
	}

	x.Gain = &gain
	return x
}

func (x *ModeColorOptions) SetTransitionTime(milliseconds int) *ModeColorOptions {
	if milliseconds < 0 {
		x.Error = fmt.Errorf("invalid transition time (ms) < 0: %d", milliseconds)
		return x
	}
	if milliseconds > 5_000 {
		x.Error = fmt.Errorf("invalid transition time (ms) > 5000: %d", milliseconds)
		return x
	}

	x.Transition = &milliseconds
	return x
}

func (x *ModeColorOptions) SetColor(color colorful.Color) *ModeColorOptions {
	x.Color = &color
	return x
}

func (x *ModeColorOptions) MustValidate() *ModeColorOptions {
	if x.Error == nil {
		return x
	}
	panic(x.Error)
}
