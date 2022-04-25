package main

import (
	"errors"
	"strings"
)

type ServiceList []Service

func (s ServiceList) Find(value string, reducer func(*Service) string) int {
	for i, service := range s {
		if reducer(&service) == value {
			return i
		}
	}
	return -1
}

type Service struct {
	Name    string
	Title   string
	FileUrl string
	DocUrl  string
	File    string
}

func getServicesFromString(str string, urlPrefix string) (ServiceList, error) {
	list := make(ServiceList, 0)
	for _, serviceString := range strings.Split(str, "|") {
		partials := strings.Split(serviceString, "::")
		if len(partials) != 3 {
			return list, errors.New("invalid Format for UI_SERVICES")
		}
		list = append(list, Service{
			Name:    partials[0],
			Title:   partials[1],
			FileUrl: urlPrefix + "/files/" + partials[2],
			DocUrl:  urlPrefix + "/" + partials[0],
			File:    partials[2],
		})
	}
	return list, nil
}
