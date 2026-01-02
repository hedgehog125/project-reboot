package servercommon

import (
	"encoding/base64"
	"net/http"

	"github.com/google/uuid"
)

func ParseUUID(str string) (uuid.UUID, *Error) {
	parsedID, stdErr := uuid.Parse(str)
	if stdErr != nil {
		return uuid.UUID{}, NewError(stdErr).
			SetStatus(http.StatusBadRequest).
			AddDetail(ErrorDetail{
				Message: "ID is not a valid UUID",
				Code:    "ID_NOT_VALID_UUID",
			})
	}
	return parsedID, nil
}

func DecodeBase64(str string) ([]byte, *Error) {
	parsedBytes, stdErr := base64.StdEncoding.DecodeString(str)
	if stdErr != nil {
		return nil, NewError(stdErr).
			SetStatus(http.StatusBadRequest).
			AddDetail(ErrorDetail{
				Message: "auth code is not valid base64",
				Code:    "MALFORMED_AUTH_CODE",
			})
	}
	return parsedBytes, nil
}
