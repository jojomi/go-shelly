package shelly

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPower(t *testing.T) {
	a := assert.New(t)
	c := NewClient(Device{
		Hostname: "http://dark-salmon",
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
