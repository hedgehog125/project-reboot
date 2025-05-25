package twofactoractions

import (
	"encoding/json"

	"github.com/hedgehog125/project-reboot/common"
)

const (
	ErrTypeEncoding = "encoding"
	// Lower level
	ErrTypeInvalidData = "invalid data"
)

// TODO: this doesn't work because ErrTypeTwoFactorAction needs to remain the second highest level category as more categories are added
// Maybe need to specify an insertion point?
// Make HighestCategory a slice. Update the constructors to either take something like this: []string{<highest>, "+", <normal>}
// Or split into 2 slice arguments

// I think something like this makes sense to do for implementing the package categories. But should the "database" category be the highest level? Or the lowest? Maybe test some scenarios with the different categories

// Replace HighestCategory with the concept of packages and subpackages, would be a slice.
// AddCategory would append after the packages and a new InsertCategory would insert before the package categories
// Example: Encode would use InsertCategory but if another package wanted to add categories it would use AddCategory

// Actually, AddCategory could always insert and AddPackage has to be used instead to commit the package path and then append after it

// Create ErrorWrapper to reduce repetition. Has a Wrap method that returns an *common.Error given an error

// Log categorisation should be done with GeneralCategory field. The common error values could have a prefix so the general category can be anywhere in the hierarchy
// e.g "constraint", common.ErrTypeDatabase, "create user"

var ErrUnknownActionType = common.NewErrorWithCategories(
	"unknown action type", common.ErrTypeTwoFactorAction,
)

func (registry *Registry) Encode(fullType string, data any) (string, *common.Error) {
	actionDef, ok := registry.actions[fullType]
	if !ok {
		return "", ErrUnknownActionType.AddCategory(ErrTypeEncoding)
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return "", common.WrapErrorWithCategories(
			err, ErrTypeInvalidData, ErrTypeEncoding, common.ErrTypeTwoFactorAction,
		)
	}

	// TODO: is there a better way to do this? With reflection maybe?
	temp := actionDef.BodyType
	err = json.Unmarshal(encoded, &temp)
	if err != nil {
		return "", common.WrapErrorWithCategories(
			err, ErrTypeInvalidData, ErrTypeEncoding, common.ErrTypeTwoFactorAction,
		)
	}

	return string(encoded), nil
}
