package tacx

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Device string
}

func Start(config Config) {
	port, err := connect(config.Device)
	if err != nil {
		log.Fatalf("unable to connect to tacx: %v", err)
	}

	commander := makeCommander(port)

	version, err := getVersion(commander)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("done %+v", version)
}
