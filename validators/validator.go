package validators

import multierror "github.com/hashicorp/go-multierror"

type Validator interface {
	Validate() (bool, error)
}

func Validate(items ...Validator) (bool, error) {
	var result error
	var check bool = true
	for _, item := range items {
		itemCheck, itemErr := item.Validate()

		if itemCheck == false {
			check = false
		}
		if itemErr != nil {
			result = multierror.Append(result, itemErr)
		}
	}
	if check == true {
		return check, nil
	}
	return check, result
}
