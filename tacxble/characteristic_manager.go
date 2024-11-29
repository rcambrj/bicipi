package tacxble

import (
	"errors"
	"fmt"

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

func (self *ServiceManager) AddService(uuid bluetooth.UUID, characteristics ...bluetooth.CharacteristicConfig) error {
	_, alreadyExists := self.services[uuid]
	if alreadyExists {
		return errors.New(fmt.Sprintf("Service with UUID %s already exists", uuid.String()))
	}

	svc := bluetooth.Service{
		UUID: uuid,
	}

	for _, c := range characteristics {
		var handle bluetooth.Characteristic
		c.Handle = &handle
		svc.Characteristics = append(svc.Characteristics, c)
		key := fmt.Sprintf("%s-%s", svc.UUID.String(), c.UUID.String())
		self.characteristics[key] = &handle
	}

	self.services[uuid] = svc

	return nil
}

func (self *ServiceManager) GetServiceIds() []bluetooth.UUID {
	ids := []bluetooth.UUID{}

	for uuid := range self.services {
		ids = append(ids, uuid)
	}

	return ids
}

func (self *ServiceManager) RegisterServices(adapter *bluetooth.Adapter) error {
	for _, svc := range self.services {
		err := adapter.AddService(&svc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *ServiceManager) GetCharacteristic(serviceUUID bluetooth.UUID, characteristicUUID bluetooth.UUID) (*bluetooth.Characteristic, error) {
	key := fmt.Sprintf("%s-%s", serviceUUID.String(), characteristicUUID.String())

	characteristic, ok := self.characteristics[key]

	if !ok {
		err := errors.New(fmt.Sprintf("No Characteristic found with UUID '%s' for service UUID '%s'", characteristicUUID.String(), serviceUUID.String()))
		return nil, err
	}

	return characteristic, nil
}
