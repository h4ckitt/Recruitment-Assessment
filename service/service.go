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

	if len(result) == 0 {
		return model.Result{}, apperror.NotFound
	}

	var meta model.Meta

	if len(result) == lim+1 {
		result = result[:len(result)-1]
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
		data    []model.Data
		meta    model.Meta
		numbers []string
	)
	lim, err := strconv.Atoi(limit)
	off, _ := strconv.Atoi(page)
	off = lim*off - lim

	if err != nil {
		return model.Result{}, apperror.BadRequest
	}
	for len(data) <= lim {
		numbers, err = s.repository.FetchPaginatedPhoneNumbersByCode(code, off, lim+1)

		if err != nil {
			if err == apperror.NotFound {
				break
			}
			log.Println(err)
			return model.Result{}, apperror.ServerError
		}

		if len(numbers) == 0 {
			return model.Result{
				Data: data,
				Meta: model.Meta{
					Next:        false,
					Prev:        false,
					CurrentPage: page,
				},
			}, nil
		}

		if len(numbers) == lim+1 {
			numbers = numbers[:len(numbers)-1]
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
		off += lim
	}

	if len(data) == 0 {
		return model.Result{}, apperror.NotFound
	}

	if len(numbers) == lim+1 {
		meta.Next = true
	}

	if off >= lim {
		meta.Prev = true
	}

	meta.CurrentPage = page

	return model.Result{
		Data: data,
		Meta: meta,
	}, nil
}
