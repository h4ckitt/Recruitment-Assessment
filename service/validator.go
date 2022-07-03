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

//NewValidator : Service To Be Used For Validation Of Country And Code.
func NewValidator() *Validator {

	// return a validator with preset information and regular expressions
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

/*Validate : checks whether the input phone number is valid
	Returns:
		- country <string>
		- code    <string>
		- number  <string>
        - valid   <bool>
*/
func (v *Validator) Validate(phone string) (string, string, string, bool) {
	// trim leading and trailing spaces from the input to avoid true negatives
	phone = strings.TrimSpace(phone)

	// loop through all countries and their regular expressions
	for country, regex := range v.CountryAndRegex {
		// check if the phone number starts with the country code of the current country
		if strings.HasPrefix(phone, fmt.Sprintf("(%s)", country.code)) {
			// check whether the phone number conforms to the registered validation regular expression
			if regex.MatchString(phone) {

				// get all the submatches of the regular expression, which are the country code and phone number (without country code)
				subMatches := regex.FindAllStringSubmatch(phone, -1)

				// it always will be on the first row
				// submatch returns only the number if it doesn't have a submatch.
				subMatch := subMatches[0]

				// return the country, extracted submatches and a value of true for validation.
				return country.name, fmt.Sprintf("+%s", subMatch[1]), subMatch[2], true
			}

			// if the number doesn't match the regex then it is not valid.
			// the country code and ordinary number need to be extracted
			sanitized := strings.TrimSpace(phone[strings.Index(phone, ")")+1:])

			return country.name, fmt.Sprintf("+%s", country.code), sanitized, false
		}
	}

	// there's no match at all, return zero values.
	return "", "", "", false
}

// GetCodeFromCountry : Get's the country code from the input country.
func (v *Validator) GetCodeFromCountry(name string) (string, error) {
	// loop through the default registered countries
	for country := range v.CountryAndRegex {
		// if the input country matches  the requested  country, return it's country code.
		if strings.ToLower(country.name) == strings.ToLower(name) {
			return country.code, nil
		}
	}

	// return a not found error
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
