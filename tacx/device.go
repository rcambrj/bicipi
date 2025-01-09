package tacx

import "github.com/rcambrj/bicipi/tacx/common"

type TacxDevice interface {
	GetVersion() (common.Version, error)
	SendControl(command common.ControlCommand) (common.ControlResponse, error)
	Close() error
}
