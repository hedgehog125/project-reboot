package servercommon

import (
	"net/http"

	"github.com/google/uuid"
)

func ParseObjectID(uuidStr string) (uuid.UUID, *Error) {
	parsed, stdErr := uuid.Parse(uuidStr)
	if stdErr != nil {
		return uuid.Nil, NewError(ErrWrapperParseObjectID.Wrap(stdErr)).
			SetStatus(http.StatusBadRequest).
			AddDetail(ErrorDetail{
				Message: "ID is not a valid UUID",
				Code:    "ID_NOT_VALID_UUID",
			}).
			DisableLogging()
	}
	return parsed, nil
}
