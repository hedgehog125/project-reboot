package servercommon

import "github.com/hedgehog125/project-reboot/common"

var ErrUnauthorized = common.NewErrorWithCategory("unauthorized", common.ErrTypeClient)
var ErrNotFound = common.NewErrorWithCategory("not found", common.ErrTypeClient)
