package usb

import (
	"fmt"
	"time"

	"github.com/rcambrj/bicipi/tacx/common"
	log "github.com/sirupsen/logrus"
)

func getVersion(t commander) (common.Version, error) {
	// while the serial connection responds reliably, the USB headunit doesn't.
	// the first few (3-5) requests appear to fail before, and a delay greater
	// than 1s between requests seems to make them all fail.
	tries := 10
	for {
		log.Info("requesting tacx version...")
		response, err := t.sendCommand(common.GetVersionCommand())
		if err != nil {
			return common.Version{}, fmt.Errorf("unable to get version: %w", err)
		}
		if !isValidFrame(response, frameTypeVersion) {
			log.Warn("received invalid frame")
			if tries > 0 {
				tries--
				time.Sleep(100 * time.Millisecond)
				continue
			}
			return common.Version{}, ErrReceivedInvalidFrame
		}

		version, err := common.GetVersionFromResponseBytes(response[24:48])
		if err != nil {
			return common.Version{}, fmt.Errorf("unable to parse usb version: %w", err)
		}

		log.WithFields(log.Fields{"version": fmt.Sprintf("%+v", version)}).Info("received tacx version")
		return version, nil
	}
}
