package helpers

// GetLimitOffset returns limit and offset based on page number
func GetLimitOffset(limit int, page int) (int, int) {
	offset := 0
	if limit == 0 {
		limit = 20
	}
	offset = limit * (page - 1)
	if offset < 0 {
		offset = 0
	}
	return offset, limit
}
