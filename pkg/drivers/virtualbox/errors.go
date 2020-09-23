package virtualbox

import "errors"

// virtualbox errors
var (
	ErrVBoxManageCommandNotFound = errors.New("the VBoxManage executable cannot be found")
)
