package util

const maxPaginationLimit, defaultPaginationLimit int = 1000, 10

func ValidateLimit(l int) int {
	if l < 1 {
		return defaultPaginationLimit
	} else if l > maxPaginationLimit {
		return maxPaginationLimit
	}
	return l
}

func ValidatePage(p int) int {
	if p < 1 {
		return 1
	}
	return p
}

func ValidatePagination(page, limit int) (int, int) {
	return ValidatePage(page), ValidateLimit(limit)
}
