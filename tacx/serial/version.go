package serial

import (
	"fmt"

	"github.com/rcambrj/bicipi/tacx/common"
	log "github.com/sirupsen/logrus"
)

func getVersion(t commander) (common.Version, error) {
	log.Info("requesting tacx version...")
	response, err := t.sendCommand(common.GetVersionCommand())
	if err != nil {
		return common.Version{}, fmt.Errorf("unable to get version: %w", err)
	}

	version, err := common.GetVersionFromResponseBytes(response[4:])
	if err != nil {
		return common.Version{}, fmt.Errorf("unable to parse serial version: %w", err)
	}

	log.WithFields(log.Fields{"version": fmt.Sprintf("%+v", version)}).Info("received tacx version")
	return version, nil
}
