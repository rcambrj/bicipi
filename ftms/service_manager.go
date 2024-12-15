package ftms

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"tinygo.org/x/bluetooth"
)

type ServiceManager struct {
	characteristics map[string]*bluetooth.Characteristic
	services        map[bluetooth.UUID]bluetooth.Service
}

func NewServiceManager() ServiceManager {
	return ServiceManager{
		characteristics: map[string]*bluetooth.Characteristic{},
		services:        map[bluetooth.UUID]bluetooth.Service{},
	}
}

func (sm *ServiceManager) AddService(uuid bluetooth.UUID, characteristics ...bluetooth.CharacteristicConfig) error {
	_, alreadyExists := sm.services[uuid]
	if alreadyExists {
		return fmt.Errorf("service with UUID %s already exists", uuid.String())
	}

	svc := bluetooth.Service{
		UUID: uuid,
	}

	log.WithField("service", svc.UUID.String()).Debug("registered service")

	for _, c := range characteristics {
		var handle bluetooth.Characteristic
		c.Handle = &handle
		svc.Characteristics = append(svc.Characteristics, c)
		key := getCharacteristicKey(svc.UUID, c.UUID)
		log.WithFields(log.Fields{
			"service":        svc.UUID.String(),
			"characteristic": c.UUID.String(),
		}).Debug("registered characteristic")
		sm.characteristics[key] = &handle
	}

	sm.services[uuid] = svc

	return nil
}

func getCharacteristicKey(svcUUID bluetooth.UUID, characteristicUUID bluetooth.UUID) string {
	return fmt.Sprintf("%s-%s", svcUUID.String(), characteristicUUID.String())
}

func (sm *ServiceManager) GetServiceIds() []bluetooth.UUID {
	ids := []bluetooth.UUID{}

	for uuid := range sm.services {
		ids = append(ids, uuid)
	}

	return ids
}

func (sm *ServiceManager) PublishServices(adapter *bluetooth.Adapter) error {
	for _, svc := range sm.services {
		err := adapter.AddService(&svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sm *ServiceManager) GetCharacteristic(serviceUUID bluetooth.UUID, characteristicUUID bluetooth.UUID) (*bluetooth.Characteristic, error) {
	key := getCharacteristicKey(serviceUUID, characteristicUUID)

	characteristic, ok := sm.characteristics[key]

	if !ok {
		err := fmt.Errorf("no Characteristic found with UUID '%s' for service UUID '%s'", characteristicUUID.String(), serviceUUID.String())
		return nil, err
	}

	return characteristic, nil
}

func (sm *ServiceManager) WriteToCharacteristic(serviceUUID bluetooth.UUID, characteristicUUID bluetooth.UUID, message []byte) (int, error) {
	c, err := sm.GetCharacteristic(serviceUUID, characteristicUUID)
	if err != nil {
		return 0, fmt.Errorf("unable to get characteristic: %w", err)
	}

	n, err := c.Write(message)
	if err != nil {
		return n, fmt.Errorf("unable to write to ble: %w", err)
	}

	return n, nil
}
