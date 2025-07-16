package jobscommon

import "github.com/hedgehog125/project-reboot/common"

const (
	ErrTypeParseVersionedType = "parse versioned type"
)

var ErrMalformedVersionedType = common.NewErrorWithCategories(
	"malformed versioned type", common.ErrTypeJobs,
)
