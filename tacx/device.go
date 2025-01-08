package tacx

import "github.com/rcambrj/bicipi/tacxcommon"

type TacxDevice interface {
	GetVersion() (tacxcommon.Version, error)
	SendControl(command tacxcommon.ControlCommand) (tacxcommon.ControlResponse, error)
	Close() error
}
