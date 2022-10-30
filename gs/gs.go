// gs.go is a sub-package pttk. A package for testing a static content
// Gopher site.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package gs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	// 3rd Party packages
	"git.mills.io/prologic/go-gopher"
)

// Service holds the description needed to startup a service
// e.g. https, http.
type Service struct {
	// Scheme holds the protocol to use, defaults to gopher if not set.
	Scheme string `json:"scheme,omitempty" `
	// Host is the hostname to use, if empty "localhost" is assumed"
	Host string `json:"host,omitempty" `
	// Port is a string holding the port number to listen on
	// An empty strings defaults port to 7000
	Port string `json:"port,omitempty" `
}

func DefaultService() *Service {
	s := new(Service)
	s.Scheme = "gopher"
	s.Host = "localhost"
	s.Port = "7000"
	return s
}

// String renders an URL version of *Service.
func (s *Service) String() string {
	r := []string{}
	if s.Scheme != "" {
		r = append(r, s.Scheme, "://")
	}
	r = append(r, s.Hostname())
	return strings.Join(r, "")
}

// Hostname returns a host+port from a *Service
func (s *Service) Hostname() string {
	r := []string{}
	if s.Host != "" {
		r = append(r, s.Host)
	}
	if s.Port != "" {
		r = append(r, ":", s.Port)
	}
	return strings.Join(r, "")
}

// GopherService describes all the configuration and
// capabilities of running a gopher service.
type GopherService struct {
	// This is the document root for static file services
	// If an empty string then assume current working directory.
	DocRoot string `json:"htdocs,omitempty"`
	// Gopher holds the service description points a *Service
	Gopher *Service `json:"gopher,omitempty"`
}

// DefaultGopherService is gopher, port 7000 on localhost.
func DefaultGopherService() *GopherService {
	gs := new(GopherService)
	gs.Gopher = DefaultService()
	gs.DocRoot = "."
	return gs
}

// LoadGopherService loads a configuration file of *Service
func LoadGopherService(setup string) (*GopherService, error) {
	var (
		gs  *GopherService
		err error
	)

	switch {
	case strings.HasSuffix(setup, ".json"):
		gs, err = loadGopherServiceJSON(setup)
	default:
		err = fmt.Errorf("%q, unknown format.", setup)
	}
	if err != nil {
		return nil, err
	}
	return gs, err
}

// loadGopherServiceJSON loads a *GopherService from a JSON file.
func loadGopherServiceJSON(setup string) (*GopherService, error) {
	src, err := os.ReadFile(setup)
	if err != nil {
		return nil, err
	}
	gs := new(GopherService)
	if err := json.Unmarshal(src, &gs); err != nil {
		return nil, err
	}
	if gs.Gopher == nil {
		gs.Gopher = DefaultService()
	}
	if gs.DocRoot == "" {
		gs.DocRoot = "./"
	}
	return gs, nil
}

// Run() starts a web service(s) described in the *GopherService struct.
func (gs *GopherService) Run() error {
	var err error
	if gs.DocRoot == "" {
		gs.DocRoot, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	log.Printf("Document root %s", gs.DocRoot)
	if gs.Gopher != nil {
		log.Printf("Listening for %s", gs.Gopher.String())
	}
	gopher.Handle("/", gopher.FileServer(gopher.Dir(gs.DocRoot)))
	return gopher.ListenAndServe(gs.Gopher.Hostname(), nil)
}
