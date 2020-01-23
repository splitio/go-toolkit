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
