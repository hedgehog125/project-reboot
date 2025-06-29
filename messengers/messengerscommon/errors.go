package messengerscommon

import "github.com/hedgehog125/project-reboot/common"

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).SetChild(common.ErrWrapperDatabase)
