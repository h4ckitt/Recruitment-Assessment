package service

import (
	"assessment/apperror"
	"assessment/model"
	"assessment/repository"
	"log"
	"regexp"
	"strings"
)

var numberRegex = regexp.MustCompile(`^\d+$`)

type NumberService struct {
	validator  NumberValidator
	repository repository.PhoneNumberRepository
}

/*NewNumberService : This starts a new service which handles the business logic of returning
phone numbers with the specified criteria
*/
func NewNumberService(validator NumberValidator, repository repository.PhoneNumberRepository) *NumberService {
	return &NumberService{
		validator:  validator,
		repository: repository,
	}
}

/*FetchPhoneNumbers : Fetches all the phone numbers in the database
Returns a paginated list of all phone numbers in the database.
*/
func (s *NumberService) FetchPhoneNumbers(page, limit string) (model.Result, error) {
	off, lim, err := validateParams("", page, limit)

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	off = lim*off - lim

	result, err := s.repository.FetchPaginatedPhoneNumbers(off, lim+1)

	if err != nil {
		log.Println(err)
		return model.Result{}, apperror.ServerError
	}

	var meta model.Meta

	if len(result) == lim+1 {
		result = result[:lim]
		meta.Next = true
	}

	if off >= lim {
		meta.Prev = true

	}

	meta.CurrentPage = page

	var data []model.Data

	for _, number := range result {
		country, code, number, valid := s.validator.Validate(number)

		state := "OK"

		if !valid {
			state = "NOK"
		}

		data = append(data, model.Data{
			Country:     country,
			CountryCode: code,
			PhoneNumber: number,
			State:       state,
		})
	}

	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	}

	finalResult := model.Result{
		Data: data,
		Meta: meta,
	}

	return finalResult, nil
}

func (s *NumberService) FilterByState(state, page, limit string) (model.Result, error) {
	/*if state != "OK" && state != "NOK" {
		return model.Result{}, apperror.BadRequest
	}*/

	p, lim, err := validateParams(state, page, limit)

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	offset := 0

	var (
		data        = make([]model.Data, lim)
		meta        model.Meta
		numbers     []string
		next        bool
		currentPage = 1
	)

	/*Fetching paginated data from the database based on status is a bit more complex because that data is non-contiguous. e.g:
		1. OK
		2. NOK
		3. NOK
		4. OK
		5. NOK
		6. OK
		7. OK
		8. NOK

	Fetching page 1 of phone numbers which are nok with a limit of 3 puts the final offset for page 1 on number 5 because that's where the third NOK number is.

	In order to fetch page 2, the offset will need to start from number 6 in order to prevent number 5 from entering the result set.

	What this function does is to start from page 1 up to page n and try to calculate the start offset to use for page n since it's not really deterministic.

	This block runs in 0(N^2) Time
	*/
	for currentPage <= p {
		if len(data) < lim {
			data = nil
			break
		}

		data = data[lim:]
		for len(data) < lim {

			numbers, err = s.repository.FetchPaginatedPhoneNumbers(offset, lim+1)

			if err != nil {
				log.Println(err)
				return model.Result{}, apperror.ServerError
			}

			if len(numbers) == 0 {
				next = false
				break
			}

			if len(numbers) == lim+1 {
				next = true
				numbers = numbers[:lim]
			} else {
				next = false
			}

			for _, number := range numbers {
				country, code, number, valid := s.validator.Validate(number)
				if (state == "OK" && valid) || (state == "NOK" && !valid) {
					data = append(data, model.Data{
						Country:     country,
						CountryCode: code,
						PhoneNumber: number,
						State:       state,
					})
				}
			}
			offset += lim
		}
		if len(data) > lim {
			next = true
		}
		currentPage++
	}

	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	} else {
		if len(data) > lim {
			data = data[:lim]
		}
		meta.Next = next
		if (p*lim - lim) >= lim {
			meta.Prev = true
		}
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil

}

func (s *NumberService) FilterByCountry(country, page, limit string) (model.Result, error) {

	p, lim, err := validateParams("", page, limit)

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

	code, err := s.validator.GetCodeFromCountry(country)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	offset := lim*p - lim

	var (
		state   string
		data    []model.Data
		meta    model.Meta
		numbers []string
		next    bool
		//currentPage = 1
	)

	for len(data) < lim {
		numbers, err = s.repository.FetchPaginatedPhoneNumbersByCode(code, offset, lim+1)

		if err != nil {
			log.Println(err)
			return model.Result{}, apperror.ServerError
		}

		if len(numbers) == 0 {
			break
		}

		if len(numbers) == lim+1 {
			next = true
			numbers = numbers[:lim]
		} else {
			next = false
		}

		for _, number := range numbers {
			country, code, number, valid := s.validator.Validate(number)
			if valid {
				state = "OK"
			} else {
				state = "NOK"
			}
			data = append(data, model.Data{
				Country:     country,
				CountryCode: code,
				PhoneNumber: number,
				State:       state,
			})
		}
		offset += lim
	}

	meta.Next = next

	if (p*lim - lim) >= lim {
		meta.Prev = true
	}

	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}

func (s *NumberService) FilterByCountryAndState(country, state, page, limit string) (model.Result, error) {

	p, lim, err := validateParams(state, page, limit)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	code, err := s.validator.GetCodeFromCountry(country)

	if err != nil {
		return model.Result{}, apperror.NotFound
	}

	state = strings.ToUpper(state)

	offset := 0

	var (
		data        = make([]model.Data, lim)
		meta        model.Meta
		numbers     []string
		next        bool
		currentPage = 1
	)

	/*Fetching paginated data from the database based on status is a bit more complex because that data is non-contiguous. e.g:
		1. OK
		2. NOK
		3. NOK
		4. OK
		5. NOK
		6. OK
		7. OK
		8. NOK

	Fetching page 1 of phone numbers which are nok with a limit of 3 puts the final offset for page 1 on number 5 because that's where the third NOK number is.

	In order to fetch page 2, the offset will need to start from number 6 in order to prevent number 5 from entering the result set.

	What this function does is to start from page 1 up to page n and try to calculate the start offset to use for page n since it's not really deterministic.

	This block runs in 0(N^2) Time
	*/
	for currentPage <= p {
		// Results From The Previous Page Are Lesser Than The Limit, Meaning There's No More Data On The Next Pl
		if len(data) < lim {
			data = nil
			break
		}

		data = data[lim:]
		for len(data) < lim {
			numbers, err = s.repository.FetchPaginatedPhoneNumbersByCode(code, offset, lim+1)

			if err != nil {
				log.Println(err)
				return model.Result{}, apperror.ServerError
			}

			if len(numbers) == 0 {
				next = false
				break
			}

			if len(numbers) == lim+1 {
				next = true
				numbers = numbers[:lim]
			} else {
				next = false
			}

			for _, number := range numbers {
				country, code, number, valid := s.validator.Validate(number)
				if (state == "OK" && valid) || (state == "NOK" && !valid) {
					data = append(data, model.Data{
						Country:     country,
						CountryCode: code,
						PhoneNumber: number,
						State:       state,
					})
				}
			}
			offset += lim
		}
		if len(data) > lim {
			next = true
		}
		currentPage++
	}

	if len(data) == 0 {
		meta.Next = false
		meta.Prev = false
	} else {
		if len(data) > lim {
			data = data[:lim]
		}
		meta.Next = next

		if (p*lim - lim) >= lim {
			meta.Prev = true
		}
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}
