package service

import (
	"assessment/apperror"
	"fmt"
	"regexp"
	"strings"
)

type country struct {
	name string
	code string
}
type Validator struct {
	CountryAndRegex map[country]*regexp.Regexp
}

func NewValidator() *Validator {
	return &Validator{
		CountryAndRegex: map[country]*regexp.Regexp{
			country{
				name: "Cameroon",
				code: "237",
			}: regexp.MustCompile(`\((237)\) ?([2368]\d{7,8})$`),
			country{
				name: "Ethiopia",
				code: "251",
			}: regexp.MustCompile(`\((251)\) ?([1-59]\d{8})$`),
			country{
				name: "Morocco",
				code: "212",
			}: regexp.MustCompile(`\((212)\) ?([5-9]\d{8})$`),
			country{
				name: "Mozambique",
				code: "258",
			}: regexp.MustCompile(`\((258)\) ?([28]\d{7,8})$`),
			country{
				name: "Uganda",
				code: "256",
			}: regexp.MustCompile(`\((256)\) ?(\d{9})$`),
		},
	}
}

func (v *Validator) Validate(phone string) (string, string, string, bool) {
	phone = strings.TrimSpace(phone)
	for country, regex := range v.CountryAndRegex {
		if strings.HasPrefix(phone, fmt.Sprintf("(%s)", country.code)) {
			if regex.MatchString(phone) {
				subMatches := regex.FindAllStringSubmatch(phone, -1)

				subMatch := subMatches[0]

				return country.name, fmt.Sprintf("+%s", subMatch[1]), subMatch[2], true
			}

			sanitized := strings.TrimSpace(phone[strings.Index(phone, ")")+1:])

			return country.name, fmt.Sprintf("+%s", country.code), sanitized, false
		}
	}

	return "", "", "", false
}

func (v *Validator) GetCodeFromCountry(name string) (string, error) {
	for country := range v.CountryAndRegex {
		if strings.ToLower(country.name) == strings.ToLower(name) {
			return country.code, nil
		}
	}

	return "", apperror.NotFound
}

/*func (v *Validator) GetCountryFromCode(code string) (string, error) {
	for country := range v.CountryAndRegex {
		if country.code == code {
			return country.name, nil
		}
	}

	return "", apperror.NotFound
}*/
