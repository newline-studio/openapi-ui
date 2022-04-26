package main

import (
	"errors"
	"os"
	"regexp"
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

func getServicesFromString(envPrefix string, envList []string, urlPrefix string) (ServiceList, error) {
	list := make(ServiceList, 0)
	exp := regexp.MustCompile("=.*$")
	for _, env := range envList {
		if strings.Index(env, envPrefix) == 0 {
			suffix := exp.ReplaceAllString(strings.Replace(env, envPrefix, "", 1), "")
			partials := strings.Split(os.Getenv(envPrefix+suffix), "|")
			if len(partials) != 2 {
				return list, errors.New("invalid format for env " + env)
			}
			list = append(list, Service{
				Name:    suffix,
				Title:   partials[0],
				FileUrl: urlPrefix + "/files/" + partials[1],
				DocUrl:  urlPrefix + "/" + suffix,
				File:    partials[1],
			})
		}
	}
	return list, nil
}
