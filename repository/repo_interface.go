package repository

type PhoneNumberRepository interface {
	FetchPaginatedPhoneNumbers(offset, limit int) ([]string, error)
	FetchPaginatedPhoneNumbersByCode(code string, offset, limit int) ([]string, error)
}
