package luxafor

import (
	"time"

	"github.com/karalabe/hid"
	"github.com/pkg/errors"
)

// Luxafor is used to access the devices.
type Luxafor struct {
	deviceInfo hid.DeviceInfo
}

type Animation byte

// LED represents the location of an individual LED to control on the Luxafor.
type LED byte

const (
	vendorID uint16 = 0x04d8
	deviceID uint16 = 0xf372

	// Commands recognized by the Luxafor.
	static Animation = 1
	fade   Animation = 2
	strobe Animation = 3
	wave   Animation = 4
	pattrn Animation = 6

	// Exported locations that represent individual--and sets of--LEDs on the Luxafor.
	FrontTop    LED = 1
	FrontMiddle LED = 2
	FrontBottom LED = 3
	BackTop     LED = 4
	BackMiddle  LED = 5
	BackBottom  LED = 6
	FrontAll    LED = 65
	BackAll     LED = 66
	All         LED = 255

	// Wave types available
	SingleSmall = 1
	SingleLarge = 2
	DoubleSmall = 3
	DoubleLarge = 4
)

// Enumerate returns a slice of attached Luxafors
func Enumerate() []Luxafor {
	infos := hid.Enumerate(vendorID, deviceID)
	luxs := make([]Luxafor, len(infos))
	for _, info := range infos {
		lux := Luxafor{
			deviceInfo: info,
		}
		luxs = append(luxs, lux)
	}

	return luxs
}

func (lux Luxafor) sendCommand(command Animation, led LED, r, g, b, speed uint8) (err error) {
	info := lux.deviceInfo
	device, err := info.Open()
	if err != nil {
		return errors.Wrap(err, "open lux")
	}

	defer func() { _ = device.Close() }() // Best effort.

	// Sets specified LED to RGB.
	if _, err := device.Write([]byte{byte(command), byte(led), r, g, b}); err != nil {
		return errors.Wrap(err, "device write")
	}
	return nil
}

// Solid turns the specified luxafor into a solid RGB color.
func (lux Luxafor) Solid(r, g, b uint8) (err error) {
	return lux.Set(All, r, g, b)
}

// Set sets a luxafor.LED to the specific RGB value.
func (lux Luxafor) Set(led LED, r, g, b uint8) (err error) {
	return lux.sendCommand(static, led, r, g, b, 0) // speed isn't used
}

// Sets sets multiple luxafor.LED to the specific RGB value.
func (lux Luxafor) Sets(leds []LED, r, g, b uint8) (err error) {
	for _, led := range leds {
		if err := lux.Set(led, r, g, b); err != nil {
			return errors.Wrap(err, "set led")
		}
	}
	return nil
}

// Fade sets the led to rgb at speed.
func (lux Luxafor) Fade(led LED, r, g, b, speed uint8) (err error) {
	return lux.sendCommand(fade, led, r, g, b, speed)
}

// Police look like da popo
func (lux Luxafor) Police(loops int) (err error) {
	for i := 0; i < loops; i++ {
		lux.Fade(FrontAll, 255, 0, 0, 255)
		lux.Fade(BackAll, 0, 0, 255, 255)
		time.Sleep(500 * time.Millisecond)
		lux.Fade(FrontAll, 0, 0, 255, 255)
		lux.Fade(BackAll, 255, 0, 0, 255)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// Off turns off the luxafor.
func (lux Luxafor) Off() (err error) {
	info := lux.deviceInfo
	device, err := info.Open()
	if err != nil {
		return errors.Wrap(err, "open lux")
	}

	defer func() { _ = device.Close() }() // Best effort.

	// Turns off the leds.
	if _, err := device.Write([]byte{byte(static), byte(All), 0, 0, 0}); err != nil {
		return errors.Wrap(err, "device write")
	}
	return nil
}
