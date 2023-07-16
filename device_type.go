package shelly

//go:generate enumer -type=DeviceType
type DeviceType int

const (
	PlugS DeviceType = iota
	BulbRGBW
)
