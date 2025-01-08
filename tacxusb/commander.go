package tacxusb

import (
	"fmt"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

var usbVendorId gousb.ID = 0x3561
var usbProductId gousb.ID = 0x1932
var usbInEndpointAddress = 0x82
var usbOutEndpointAddress = 0x02

type commander struct {
	ctx         *gousb.Context
	dev         *gousb.Device
	iface       *gousb.Interface
	ifaceDone   func()
	inEndpoint  *gousb.InEndpoint
	outEndpoint *gousb.OutEndpoint
}

func makeCommander() (*commander, error) {
	ctx := gousb.NewContext()
	ctx.Debug(1)

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == usbVendorId && desc.Product == usbProductId
	})
	if err != nil {
		return &commander{}, fmt.Errorf("unable to open usb devices: %w", err)
	}

	if len(devs) < 1 {
		return &commander{}, fmt.Errorf("unable to find usb device. is it connected?")
	}
	if len(devs) > 1 {
		return &commander{}, fmt.Errorf("found too many matching usb devices")
	}
	dev := devs[0]
	log.WithFields(log.Fields{"dev": fmt.Sprintf("%+v", dev)}).Info("found usb device")

	iface, ifaceDone, err := dev.DefaultInterface()
	if err != nil {
		return &commander{}, fmt.Errorf("unable to open usb interface: %w", err)
	}

	inEndpoint, err := iface.InEndpoint(usbInEndpointAddress)
	if err != nil {
		return &commander{}, fmt.Errorf("unable to open usb input endpoint: %w", err)
	}
	outEndpoint, err := iface.OutEndpoint(usbOutEndpointAddress)
	if err != nil {
		return &commander{}, fmt.Errorf("unable to open usb output endpoint: %w", err)
	}

	return &commander{
		ctx:         ctx,
		dev:         dev,
		iface:       iface,
		ifaceDone:   ifaceDone,
		inEndpoint:  inEndpoint,
		outEndpoint: outEndpoint,
	}, nil
}

func (c *commander) sendCommand(command []byte) ([]byte, error) {
	log.WithFields(log.Fields{"command": command}).Trace("sending usb command")

	// TODO: reset input buffer before sending command
	_, err := c.outEndpoint.Write(command)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to write to usb port: %w", err)
	}

	response := make([]byte, 0, 64)
	_, err = c.inEndpoint.Read(response)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to read from usb port: %w", err)
	}

	return response, nil
}

func (u *commander) close() error {
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
