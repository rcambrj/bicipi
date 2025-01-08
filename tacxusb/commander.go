package tacxusb

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

var usbVendorId gousb.ID = 0x3561
var usbProductId gousb.ID = 0x1932
var usbInEndpointAddress = 0x82
var usbOutEndpointAddress = 0x02

type c struct {
	ctx         *gousb.Context
	dev         *gousb.Device
	iface       *gousb.Interface
	ifaceDone   func()
	inEndpoint  *gousb.InEndpoint
	outEndpoint *gousb.OutEndpoint
}

func makeCommander() (commander, error) {
	ctx := gousb.NewContext()
	ctx.Debug(1)

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == usbVendorId && desc.Product == usbProductId
	})
	if err != nil {
		return &c{}, fmt.Errorf("unable to open usb devices: %w", err)
	}

	if len(devs) < 1 {
		return &c{}, fmt.Errorf("unable to find usb device. is it connected?")
	}
	if len(devs) > 1 {
		return &c{}, fmt.Errorf("found too many matching usb devices")
	}
	dev := devs[0]

	iface, ifaceDone, err := dev.DefaultInterface()
	if err != nil {
		return &c{}, fmt.Errorf("unable to open usb interface: %w", err)
	}

	inEndpoint, err := iface.InEndpoint(usbInEndpointAddress)
	if err != nil {
		return &c{}, fmt.Errorf("unable to open usb input endpoint: %w", err)
	}
	outEndpoint, err := iface.OutEndpoint(usbOutEndpointAddress)
	if err != nil {
		return &c{}, fmt.Errorf("unable to open usb output endpoint: %w", err)
	}

	log.WithFields(log.Fields{
		"iface":       fmt.Sprintf("%+v", iface),
		"inEndpoint":  fmt.Sprintf("%+v", inEndpoint),
		"outEndpoint": fmt.Sprintf("%+v", outEndpoint),
	}).Info("connected to usb device")

	return &c{
		ctx:         ctx,
		dev:         dev,
		iface:       iface,
		ifaceDone:   ifaceDone,
		inEndpoint:  inEndpoint,
		outEndpoint: outEndpoint,
	}, nil
}

var startOfFrame = []byte{0x12, 0xCC}
var (
	frameTypeVersion = []byte{0x03, 0x0C}
	frameTypeControl = []byte{0x03, 0x13}
)

// checks whether the frame is valid because sometimes the head unit will
// return garbage, but certain bytes seem to be correct when the response
// is valid - check those
// it seems that the incorrect response is sometimes due to a delay: perhaps the
// headunit takes some time to issue a request to the motor brake and update
// its internal cache? it then seems that this cache expires after ~1 second.
// TODO: there might be a checksum in the frame, check it
func isValidFrame(frame []byte, frameType []byte) bool {
	return len(frame) >= 48 && bytes.Equal(frame[0:2], startOfFrame) && bytes.Equal(frame[24:26], frameType)
}

var ErrReceivedInvalidFrame = errors.New("received invalid frame")

func (c *c) sendCommand(command []byte) ([]byte, error) {
	log.WithFields(log.Fields{"command": command}).Trace("sending usb command")

	// TODO: reset input buffer before sending command
	_, err := c.outEndpoint.Write(command)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to write to usb port: %w", err)
	}

	response := make([]byte, 64)
	_, err = c.inEndpoint.Read(response)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read from usb port: %w", err)
	}

	log.WithFields(log.Fields{"response": response}).Trace("received usb response")
	return response, nil
}

func (u *c) close() error {
	if u.iface != nil {
		u.iface.Close()
	}
	if u.ifaceDone != nil {
		u.ifaceDone()
	}
	if u.dev != nil {
		err := u.dev.Close()
		if err != nil {
			return fmt.Errorf("unable to close usb: %w", err)
		}
	}
	if u.ctx != nil {
		err := u.ctx.Close()
		if err != nil {
			return fmt.Errorf("unable to close usb: %w", err)
		}
	}
	return nil
}
