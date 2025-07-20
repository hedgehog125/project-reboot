package messengerscommon

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeReadUserContacts = "read user contacts"
)

var ErrWrapperDatabase = common.NewErrorWrapper(common.ErrTypeMessengers).SetChild(common.ErrWrapperDatabase)

var ErrNoTxInContext = common.ErrNoTxInContext.AddCategory(common.ErrTypeMessengers)
