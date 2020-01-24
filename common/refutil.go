package common

func StringRef(str string) *string {
	return &str
}

func IntRef(number int) *int {
	return &number
}

func Int64Ref(number int64) *int64 {
	return &number
}

func Int64Value(number *int64) int64 {
	if number == nil {
		return 0
	}
	return *number
}

func IntRefOrNil(number int) *int {
	if number == 0 {
		return nil
	}
	return IntRef(number)
}

func Int64RefOrNil(number int64) *int64 {
	if number == 0 {
		return nil
	}
	return Int64Ref(number)
}

func StringRefOrNil(str string) *string {
	if str == "" {
		return nil
	}
	return StringRef(str)
}
