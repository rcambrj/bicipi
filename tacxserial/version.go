package tacxserial

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacxcommon"
	log "github.com/sirupsen/logrus"
)

func getVersion(t commander) (tacxcommon.Version, error) {
	log.Info("requesting tacx version...")
	response, err := t.sendCommand(tacxcommon.GetVersionCommand())
	if err != nil {
		return tacxcommon.Version{}, fmt.Errorf("unable to get version: %w", err)
	}

	version, err := tacxcommon.GetVersionFromResponseBytes(response[4:])
	if err != nil {
		return tacxcommon.Version{}, fmt.Errorf("unable to parse serial version: %w", err)
	}

	log.WithFields(log.Fields{"version": fmt.Sprintf("%+v", version)}).Info("received tacx version")
	return version, nil
}
