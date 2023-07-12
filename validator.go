package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var ErrorFiledNotValid error = errors.New("filed not valid")
var ErrorTagNotValid error = errors.New("tag not valid")

func Validate(v interface{}) (error, ValidationErrors) {
	vValue := reflect.ValueOf(v)
	if vValue.Kind() != reflect.Struct {
		return errors.New("not struct type"), nil
	}
	validationErrors := make(ValidationErrors, 0, 10)
	vType := vValue.Type()
	compile := regexp.MustCompile("^validate:*")
	for i := 0; i < vType.NumField(); i++ {
		tag := vType.Field(i).Tag
		stringTag := string(tag)
		if math := compile.MatchString(stringTag); !math {
			validationErrors = append(validationErrors, ValidationError{
				Field: vType.Field(i).Name,
				Err:   errors.New("cant find validate tag or tag matching"),
			})
		}
		typeFiled := vType.Field(i).Type
		trim := strings.Trim(stringTag[10:], "\"")
		split := strings.Split(trim, "|")
		var validateErr ValidationError
		switch typeFiled.Kind() {
		case reflect.Int:
			validateErr = intValidate(split, vType.Field(i).Name, vValue.Field(i).Int())
		case reflect.String:
			validateErr = stringValidate(split, vType.Field(i).Name, vValue.Field(i).String())
		case reflect.Slice:
			if typeFiled.Elem().Kind() == reflect.Int {
				for j := 0; j < vValue.Field(i).Len(); j++ {
					sliceElem := vValue.Field(i).Index(j)
					validateErr = intValidate(split, vType.Field(i).Name, sliceElem.Int())
					validationErrors = append(validationErrors, validateErr)
				}
			} else if typeFiled.Elem().Kind() == reflect.String {
				for j := 0; j < vValue.Field(i).Len(); j++ {
					sliceElem := vValue.Field(i).Index(j)
					validateErr = stringValidate(split, vType.Field(i).Name, sliceElem.String())
					validationErrors = append(validationErrors, validateErr)
				}
			}
		}
		if validateErr.Err != nil {
			validationErrors = append(validationErrors, validateErr)
		}
	}
	return nil, validationErrors
}

func intValidate(tagArray []string, filedName string, filedValue int64) ValidationError {
	minTagRegExp := regexp.MustCompile("^min:*")
	maxTagRegExp := regexp.MustCompile("^max:*")
	inTagRegExp := regexp.MustCompile("^in:*")
	for _, tag := range tagArray {
		if match := minTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "min:")
			min, err := strconv.Atoi(trimStr)
			if err != nil {
				fmt.Println("tag not valid min")
				return ValidationError{
					Field: filedName,
					Err:   ErrorTagNotValid,
				}
			}
			if int64(min) < filedValue {
				fmt.Println("valid min")
			} else {
				fmt.Println("not valid min")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
		if match := maxTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "max:")
			max, err := strconv.Atoi(trimStr)
			if err != nil {
				fmt.Println("tag not valid max")
				return ValidationError{
					Field: filedName,
					Err:   ErrorTagNotValid,
				}
			}
			if int64(max) > filedValue {
				fmt.Println("valid max")
			} else {
				fmt.Println("not valid max")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
		if match := inTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "in:")
			splitTrimStr := strings.Split(trimStr, ",")
			var isFind bool
			for _, s := range splitTrimStr {
				convS, err := strconv.Atoi(s)
				if err != nil {
					fmt.Println("tag not valid max")
					return ValidationError{
						Field: filedName,
						Err:   ErrorTagNotValid,
					}
				}
				if filedValue == int64(convS) {
					isFind = true
					break
				}
			}
			if isFind {
				fmt.Println("valid in")
				continue
			} else {
				fmt.Println("not valid in")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
	}
	return ValidationError{}
}

func stringValidate(tagArray []string, filedName string, filedValue string) ValidationError {
	lenTagRegExp := regexp.MustCompile("^len:*")
	regexpTagRegExp := regexp.MustCompile("^regexp:*")
	inTagRegExp := regexp.MustCompile("^in:*")
	for _, tag := range tagArray {
		if match := lenTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "len:")
			length, err := strconv.Atoi(trimStr)
			if err != nil {
				fmt.Println("tag not valid length")
				return ValidationError{
					Field: filedName,
					Err:   ErrorTagNotValid,
				}
			}
			if length == len(filedValue) {
				fmt.Println("valid len")
			} else {
				fmt.Println("not valid len")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
		if match := regexpTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "regexp:")
			reg, err := regexp.Compile(trimStr)
			if err != nil {
				fmt.Println("tag not valid length")
				return ValidationError{
					Field: filedName,
					Err:   ErrorTagNotValid,
				}
			}
			if reg.MatchString(filedValue) {
				fmt.Println("valid regexp")
			} else {
				fmt.Println("not valid regexp")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
		if match := inTagRegExp.MatchString(tag); match {
			_, trimStr, _ := strings.Cut(tag, "in:")
			splitTrimStr := strings.Split(trimStr, ",")
			var isFind bool
			for _, s := range splitTrimStr {
				if filedValue == s {
					isFind = true
					break
				}
			}
			if isFind {
				fmt.Println("valid in")
				continue
			} else {
				fmt.Println("not valid in")
				return ValidationError{
					Field: filedName,
					Err:   ErrorFiledNotValid,
				}
			}
		}
	}
	return ValidationError{}
}
