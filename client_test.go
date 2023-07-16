package shelly

import (
	"github.com/lucasb-eyer/go-colorful"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPower(t *testing.T) {
	a := assert.New(t)
	c := NewClient(Device{
		Hostname: "http://dark-salmon",
		Type:     PlugS,
	})

	err := c.SetPowerOn()
	a.Nil(err)
	on, err := c.IsPowerOn()
	a.Nil(err)
	a.True(on)

	time.Sleep(5 * time.Second)

	// need a consumer attached and running, so that we have a value > 0 here
	pow, err := c.GetCurrentPower()
	a.Nil(err)
	a.True(pow > 0)

	err = c.SetPowerOff()
	a.Nil(err)
	on, err = c.IsPowerOn()
	a.Nil(err)
	a.False(on)

	time.Sleep(2 * time.Second)

	pow, err = c.GetCurrentPower()
	a.Nil(err)
	a.Equal(0.0, pow)
}

func TestBulbColor(t *testing.T) {
	a := assert.New(t)
	c := NewClient(Device{
		Hostname: "http://hot-pink",
		Type:     BulbRGBW,
	})

	err := c.SetPowerOn()
	a.Nil(err)
	on, err := c.IsPowerOn()
	a.Nil(err)
	a.True(on)

	whiteOpts := NewModeWhiteOptions().
		SetBrightness(34).
		SetTransitionTime(2000).
		SetTemp(6500).
		MustValidate()
	err = c.SetModeWhite(*whiteOpts)
	a.Nil(err)

	isWhite, err := c.IsModeWhite()
	a.Nil(err)
	a.True(isWhite)
	//brightness, err := c.GetBrightness()

	time.Sleep(3 * time.Second)

	color, err := colorful.Hex("#5599aa")
	a.Nil(err)
	colorOpts := NewModeColorOptions().
		SetBrightness(12).
		SetTransitionTime(2000).
		SetGain(100).
		SetColor(color).
		MustValidate()
	err = c.SetModeColor(*colorOpts)
	a.Nil(err)

	isColor, err := c.IsModeColor()
	a.Nil(err)
	a.True(isColor)
	isWhite, err = c.IsModeWhite()
	a.Nil(err)
	a.False(isWhite)

	time.Sleep(3 * time.Second)

	err = c.SetPowerOff()
	a.Nil(err)
	on, err = c.IsPowerOn()
	a.Nil(err)
	a.False(on)
}
