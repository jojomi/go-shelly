package shelly

import (
	"fmt"
)

type ModeWhiteOptions struct {
	Brightness *int
	Transition *int
	Temp       *int
	Error      error
}

func NewModeWhiteOptions() *ModeWhiteOptions {
	return &ModeWhiteOptions{}
}

func (x *ModeWhiteOptions) SetBrightness(brightness int) *ModeWhiteOptions {
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

func (x *ModeWhiteOptions) SetTransitionTime(milliseconds int) *ModeWhiteOptions {
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

func (x *ModeWhiteOptions) SetTemp(kelvin int) *ModeWhiteOptions {
	if kelvin < 3_000 {
		x.Error = fmt.Errorf("invalid temp (kelvin) < 0: %d", kelvin)
		return x
	}
	if kelvin > 6_500 {
		x.Error = fmt.Errorf("invalid temp (kelvin) > 5000: %d", kelvin)
		return x
	}

	x.Temp = &kelvin
	return x
}

func (x *ModeWhiteOptions) MustValidate() *ModeWhiteOptions {
	if x.Error == nil {
		return x
	}
	panic(x.Error)
}
