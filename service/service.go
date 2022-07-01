package service

import (
	"assessment/apperror"
	"assessment/model"
	"assessment/repository"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var numberRegex = regexp.MustCompile(`^\d+$`)

type NumberService struct {
	validator  NumberValidator
	repository repository.PhoneNumberRepository
}

func NewNumberService(validator NumberValidator, repository repository.PhoneNumberRepository) *NumberService {
	return &NumberService{
		validator:  validator,
		repository: repository,
	}
}

func (s *NumberService) FetchPhoneNumbers(page, limit string) (model.Result, error) {
	if page == "" {
		page = "1"
	}

	if limit == "" || limit == "0" {
		limit = "10"
	}

	if !numberRegex.MatchString(limit) || !numberRegex.MatchString(page) {
		return model.Result{}, apperror.BadRequest
	}

	off, _ := strconv.Atoi(page)
	lim, _ := strconv.Atoi(limit)

	off = lim*off - lim

	result, err := s.repository.FetchPaginatedPhoneNumbers(off, lim+1)

	if err != nil {
		log.Println(err)
		return model.Result{}, apperror.ServerError
	}

	/*if len(result) == 0 {
		return model.Result{}, apperror.NotFound
	}*/

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

	finalResult := model.Result{
		Data: data,
		Meta: meta,
	}

	return finalResult, nil
}

func (s *NumberService) FilterByCountryAndState(country, state, page, limit string) (model.Result, error) {
	if state == "" {
		state = "OK"
	}

	if page == "" {
		page = "1"
	}

	if limit == "" || limit == "0" {
		limit = "10"
	}

	if !numberRegex.MatchString(limit) || !numberRegex.MatchString(page) {
		return model.Result{}, apperror.BadRequest
	}

	state = strings.ToUpper(state)

	if state != "OK" && state != "NOK" {
		return model.Result{}, apperror.BadRequest
	}

	if country == "" {
		return model.Result{}, apperror.BadRequest
	}

	code, err := s.validator.GetCodeFromCountry(country)

	if err != nil {
		log.Println("Country isn't recognized by the validator at the moment")
		return model.Result{}, apperror.NotFound
	}

	var (
		data        []model.Data
		meta        model.Meta
		numbers     []string
		next        bool
		currentPage int = 1
	)
	lim, err := strconv.Atoi(limit)
	p, _ := strconv.Atoi(page)
	offset := 0

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}

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
		data = nil
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
				numbers = numbers[:len(numbers)-1]
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
		currentPage++
	}

	meta.Next = next

	if (p*lim - lim) >= lim {
		meta.Prev = true
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}
