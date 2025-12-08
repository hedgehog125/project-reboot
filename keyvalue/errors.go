package keyvalue

import "github.com/NicoClack/cryptic-stash/common"

const (
	ErrTypeGetValue = "get value"
	ErrTypeSetValue = "set value"
	ErrTypeInitAll  = "init all values"

	// Lower level
	ErrTypeEncode          = "encode value"
	ErrTypeDecode          = "decode value"
	ErrTypeInitInvalidType = "init invalid type"
)

var ErrUnknownName = common.NewErrorWithCategories("unknown name")
var ErrWrongPointerType = common.NewErrorWithCategories("wrong pointer type, must be a pointer to the same type as the definition")

var ErrWrapperGetValue = common.NewErrorWrapper(ErrTypeGetValue)
var ErrWrapperSetValue = common.NewErrorWrapper(ErrTypeSetValue)
var ErrWrapperInitAll = common.NewErrorWrapper(ErrTypeInitAll)

var ErrWrapperEncode = common.NewErrorWrapper(ErrTypeEncode)
var ErrWrapperDecode = common.NewErrorWrapper(ErrTypeDecode)
var ErrWrapperInitInvalidType = common.NewErrorWrapper(ErrTypeInitInvalidType)
