package main

import (
	"errors"
	"strings"
)

type ServiceList []Service

func (s ServiceList) Find(file string) int {
	for i, service := range s {
		if service.File == file {
			return i
		}
	}
	return -1
}

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	File string `json:"-"`
}

func getServicesFromString(str string, urlPrefix string) (ServiceList, error) {
	list := make(ServiceList, 0)
	for _, serviceString := range strings.Split(str, "|") {
		partials := strings.Split(serviceString, "::")
		if len(partials) != 2 {
			return list, errors.New("invalid Format for UI_SERVICES")
		}
		list = append(list, Service{
			Name: partials[0],
			URL:  urlPrefix + "/file/" + partials[1],
			File: partials[1],
		})
	}
	return list, nil
}
