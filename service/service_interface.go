package service

type NumberValidator interface {
	Validate(string) (string, string, string, bool)
	GetCodeFromCountry(name string) (string, error)
}
