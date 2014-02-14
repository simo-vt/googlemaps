// Copyright 2014 Simeon Totev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Google Maps WebServices HTTP Client API.

The idea is simple. You create a service, you assign parameters to it and you execute it with the Execute function. More details in the Examples section.

This client API supports the two types of business authentication, described here - https://developers.google.com/maps/documentation/business/webservices/auth
	// For "key" authentication you must set the API Key for the service with the function
	func (s *Service) SetApiKey(aKey string)
	// For "client" authentication you must set a Client ID and Client Secret with the functions
	func (s *Service) SetClientId(cId string)
	func (s *Service) SetClientSecret(cSecret string)

Official API documentation - https://developers.google.com/maps/documentation/webservices/. Please note that you should familiarize yourself with the terms of Google Maps API before using it.
*/
package googlemaps

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
The constants you will be using are:
	googlemaps.GOOGLEMAPS_AUTH_CLIENT - used if we want to use "client" type of authentication. For business users. Note the limits explained in the Google Maps API.
	googlemaps.GOOGLEMAPS_AUTH_KEY - used if we want to use "key" type of authentication. For business users. Note the limits explained in the Google Maps API.
	googlemaps.GOOGLEMAPS_AUTH_EMPTY - used if we do not want to authenticate. Note the limits explained in the Google Maps API.
*/
const (
	GOOGLEMAPS_AUTH_CLIENT   string = "client"
	GOOGLEMAPS_AUTH_KEY      string = "key"
	GOOGLEMAPS_AUTH_EMPTY    string = ""
	GOOGLEMAPS_URL           string = "https://maps.googleapis.com/maps/api/"
	GOOGLEMAPS_RESPONSE_JSON string = "json"
	GOOGLEMAPS_RESPONSE_XML  string = "xml"
)

type Service struct {
	servicePath, requestMethod, requestType, authType, apiKey, clientId, clientSecret string
	business                                                                          bool
	params                                                                            url.Values
}

type Stringer interface {
	String() string
}

/*
NewService creates a new Google Maps API service. The passed parameters are:

sPath string - this is the name of the service you wish to call. For example: "place/event/add" or "directions".

rMethod string - this is the request method. You will most often use GET, but with the case of the Places API you may need to use POST.

aType string - can be one of the following:
	googlemaps.GOOGLEMAPS_AUTH_KEY
	googlemaps.GOOGLEMAPS_AUTH_CLIENT
	googlemaps.GOOGLEMAPS_AUTH_EMPTY

All three are explained in the constants section.
*/
func NewService(sPath string, rMethod string, aType string) (*Service, error) {
	s := &Service{servicePath: sPath}

	err := s.setRequestType(GOOGLEMAPS_RESPONSE_JSON)
	if err != nil {
		return nil, err
	}

	err = s.SetRequestMethod(rMethod)
	if err != nil {
		return nil, err
	}

	err = s.SetAuthType(aType)
	if err != nil {
		return nil, err
	}

	s.params = make(url.Values)

	return s, nil
}

/*
String outputs the service method and URL as a string.
*/
func (s *Service) String() string {
	query, err := s.constructQuery()

	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	return fmt.Sprintf("%v %v\n", s.requestMethod, query)
}

/*
SetServicePath sets the path of the service.

Example paths are: "place/event/add" or "directions"
*/
func (s *Service) SetServicePath(sPath string) {
	s.servicePath = sPath
}

/*
If you have used the "key" authentication type, you need to set the API Key with this function before calling Execute.
*/
func (s *Service) SetApiKey(aKey string) {
	s.apiKey = aKey
}

/*
If you have used the "client" authentication type, you need to set the Client ID with this function before calling Execute.
*/
func (s *Service) SetClientId(cId string) {
	s.clientId = cId
}

/*
If you have used the "client" authentication type, you need to set the Client Secret with this function before calling Execute.

Note that the client secret is a base-64 encoded string for URL's. Basically, the one you get in the e-mail after you register a business account with Google Maps. An example key is "vNIXE0xscrmjlyV-12Nj_BvUPaw="

More information here - https://developers.google.com/maps/documentation/business/webservices/auth#generating_valid_signatures
*/
func (s *Service) SetClientSecret(cSecret string) {
	s.clientSecret = cSecret
}

/*
SetRequestMethod sets the request method of the service.

In most cases, you will be using GET, but in the case of the Places API, you might need to use POST.
*/
func (s *Service) SetRequestMethod(rMethod string) error {
	allowedTypes := []string{"POST", "GET"}
	if !inSlice(rMethod, allowedTypes) {
		err := errors.New("Invalid request method passed. Should be either POST or GET")
		return err
	}
	s.requestMethod = rMethod
	return nil
}

/*
SetAuthType sets the authentication type of the service. Can (and should) be one of these:
	googlemaps.GOOGLEMAPS_AUTH_CLIENT
	googlemaps.GOOGLEMAPS_AUTH_KEY
	googlemaps.GOOGLEMAPS_AUTH_EMPTY
*/
func (s *Service) SetAuthType(aType string) error {
	allowedTypes := []string{GOOGLEMAPS_AUTH_EMPTY, GOOGLEMAPS_AUTH_KEY, GOOGLEMAPS_AUTH_CLIENT}
	if !inSlice(aType, allowedTypes) {
		err := errors.New("Invalid authentication type passed. Should be either \"\", \"client\" or \"key\"")
		return err
	}
	s.authType = aType
	s.business = s.authType != GOOGLEMAPS_AUTH_EMPTY

	return nil
}

/*
AddParam adds a request parameter to the service request. Note that you should always use this function to set the "sensor" parameter to either "true" or "false"
*/
func (s *Service) AddParam(key, value string) {
	s.params.Add(key, value)
}

/*
Removes a parameter
*/
func (s *Service) RemoveParam(key string) {
	s.params.Del(key)
}

/*
Execute is the main function of this client API. You should pass the request body as a string.
In most cases, this should be an empty string, since the requests to the API are mostly GET. In the case of the Places API the requestBody string should contain a JSON object with the parameters, like in this link: https://developers.google.com/places/documentation/actions#event_add
*/
func (s *Service) Execute(requestBody string) (map[string]interface{}, error) {
	var m interface{}
	err := s.executeReal(requestBody, &m)
	if err != nil {
		return nil, err
	}

	n := m.(map[string]interface{})

	return n, nil
}

func (s *Service) setRequestType(rType string) error {
	allowedTypes := []string{GOOGLEMAPS_RESPONSE_XML, GOOGLEMAPS_RESPONSE_JSON}
	if !inSlice(rType, allowedTypes) {
		err := errors.New("Invalid request type passed. Should be either " + GOOGLEMAPS_RESPONSE_JSON + " or " + GOOGLEMAPS_RESPONSE_XML)
		return err
	}
	s.requestType = rType
	return nil
}

func (s *Service) executeReal(requestBody string, v interface{}) error {
	url, err := s.constructQuery()

	if err != nil {
		return err
	}

	var resp *http.Response

	if s.requestMethod == "GET" {
		resp, err = http.Get(url)
	} else if s.requestMethod == "POST" {
		reader := strings.NewReader(requestBody)

		var rType string

		if s.requestType == GOOGLEMAPS_RESPONSE_JSON {
			rType = "application/json"
		} else if s.requestType == GOOGLEMAPS_RESPONSE_XML {
			rType = "application/xml"
		}

		resp, err = http.Post(url, rType, reader)
	}
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.Status != "200 OK" {
		httperr := errors.New(resp.Status + ": " + string(responseBody))
		return httperr
	}

	if s.requestType == GOOGLEMAPS_RESPONSE_JSON {
		json.Unmarshal(responseBody, &v)
	} else if s.requestType == GOOGLEMAPS_RESPONSE_XML {

	}

	return nil
}

func (s *Service) constructQuery() (string, error) {
	s.params.Del("client")
	s.params.Del("signature")
	s.params.Del("key")

	tempUrl := GOOGLEMAPS_URL + s.servicePath + "/" + s.requestType + "?"

	if s.business {
		if s.authType == GOOGLEMAPS_AUTH_CLIENT {
			if s.clientId == "" || s.clientSecret == "" {
				err := errors.New("Client ID and/or Client Secret are not set. You should set them if you wish to use the \"client\" authentication method.")
				return "", err
			}
			s.params.Add("client", s.clientId)

			parsedUrl, err := url.Parse(tempUrl)
			if err != nil {
				return "", err
			}

			pathQuery := parsedUrl.Path + "?" + s.params.Encode()

			newSecret := strings.Replace(s.clientSecret, "-", "+", -1)
			newSecret = strings.Replace(newSecret, "_", "/", -1)

			decodedSecret, err := base64.StdEncoding.DecodeString(newSecret)
			if err != nil {
				return "", err
			}

			mac := hmac.New(sha1.New, decodedSecret)
			mac.Write([]byte(pathQuery))

			secretResult := mac.Sum(nil)

			newSignature := base64.StdEncoding.EncodeToString(secretResult)
			newSignature = strings.Replace(newSignature, "+", "-", -1)
			newSignature = strings.Replace(newSignature, "/", "_", -1)

			s.params.Add("signature", newSignature)
		} else if s.authType == GOOGLEMAPS_AUTH_KEY {
			if s.apiKey == "" {
				err := errors.New("API Key is not set. You should set them if you wish to use the \"key\" authentication method.")
				return "", err
			}
			s.params.Add("key", s.apiKey)
		}
	}

	return tempUrl + s.params.Encode(), nil
}

func inSlice(key string, slice []string) bool {
	for _, val := range slice {
		if val == key {
			return true
		}
	}
	return false
}
