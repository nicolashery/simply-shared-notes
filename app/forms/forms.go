package forms

type Errors = map[string][]string

func EmptyErrors() Errors {
	return make(map[string][]string)
}

func HasError(errors Errors, name string) bool {
	_, ok := errors[name]
	return ok
}

func GetErrors(errors Errors, name string) []string {
	errs, ok := errors[name]

	if !ok {
		return []string{}
	}

	return errs
}
