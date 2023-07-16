// Code generated by "enumer -type=DeviceType"; DO NOT EDIT.

package shelly

import (
	"fmt"
	"strings"
)

const _DeviceTypeName = "PlugSBulbRGBW"

var _DeviceTypeIndex = [...]uint8{0, 5, 13}

const _DeviceTypeLowerName = "plugsbulbrgbw"

func (i DeviceType) String() string {
	if i < 0 || i >= DeviceType(len(_DeviceTypeIndex)-1) {
		return fmt.Sprintf("DeviceType(%d)", i)
	}
	return _DeviceTypeName[_DeviceTypeIndex[i]:_DeviceTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _DeviceTypeNoOp() {
	var x [1]struct{}
	_ = x[PlugS-(0)]
	_ = x[BulbRGBW-(1)]
}

var _DeviceTypeValues = []DeviceType{PlugS, BulbRGBW}

var _DeviceTypeNameToValueMap = map[string]DeviceType{
	_DeviceTypeName[0:5]:       PlugS,
	_DeviceTypeLowerName[0:5]:  PlugS,
	_DeviceTypeName[5:13]:      BulbRGBW,
	_DeviceTypeLowerName[5:13]: BulbRGBW,
}

var _DeviceTypeNames = []string{
	_DeviceTypeName[0:5],
	_DeviceTypeName[5:13],
}

// DeviceTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DeviceTypeString(s string) (DeviceType, error) {
	if val, ok := _DeviceTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _DeviceTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DeviceType values", s)
}

// DeviceTypeValues returns all values of the enum
func DeviceTypeValues() []DeviceType {
	return _DeviceTypeValues
}

// DeviceTypeStrings returns a slice of all String values of the enum
func DeviceTypeStrings() []string {
	strs := make([]string, len(_DeviceTypeNames))
	copy(strs, _DeviceTypeNames)
	return strs
}

// IsADeviceType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DeviceType) IsADeviceType() bool {
	for _, v := range _DeviceTypeValues {
		if i == v {
			return true
		}
	}
	return false
}