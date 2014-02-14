// Copyright 2014 Simeon Totev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Google Maps WebServices HTTP Client API.

Official API documentation - https://developers.google.com/maps/documentation/webservices/. Please note that you should familiarize yourself with the terms of Google Maps API before using it.
*/
package googlemaps

import (
	"fmt"
	"strings"
	"testing"
)

func TestNewService(t *testing.T) {
	_, err := NewService("timezone", "GET", GOOGLEMAPS_AUTH_EMPTY)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestExecute_forEmptyAuth(t *testing.T) {
	s, err := NewService("distancematrix", "GET", GOOGLEMAPS_AUTH_EMPTY)
	if err != nil {
		t.Error(err)
		return
	}

	s.AddParam("sensor", "false")
	s.AddParam("origins", "Sofia, Bulgaria")
	s.AddParam("destinations", "Berlin, Germany")

	_, err = s.Execute("")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestExecute_forKeyAuth(t *testing.T) {
	s, err := NewService("geocode", "GET", GOOGLEMAPS_AUTH_KEY)
	if err != nil {
		t.Error(err)
		return
	}

	s.AddParam("sensor", "false")
	s.AddParam("address", "5 James Bourchier, Sofia, Bulgaria")

	s.SetApiKey("TEST_API_KEY")

	_, err = s.Execute("")

	if err != nil {
		t.Error(err)
		return
	}
}

func TestExecute_forClientAuth(t *testing.T) {
	s, err := NewService("elevation", "GET", GOOGLEMAPS_AUTH_CLIENT)
	if err != nil {
		t.Error(err)
		return
	}

	s.AddParam("sensor", "false")
	s.AddParam("locations", "40.714728,-73.998672")

	s.SetClientId("TEST_CLIENT_ID")
	s.SetClientSecret("vNIXE0xscrmjlyV-12Nj_BvUPaw=")

	_, err = s.Execute("")

	// We are sure that it will give 403 Forbidden, because we have entered fake credentials.
	if !strings.Contains(err.Error(), "403 Forbidden") {
		t.Error(err)
		return
	}
}

func ExampleService_Execute_forGET() {
	s, err := NewService("elevation", "GET", GOOGLEMAPS_AUTH_EMPTY)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.AddParam("sensor", "false")
	s.AddParam("locations", "40.714728,-73.998672")

	resp, err := s.Execute("")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(resp["status"])
	//Output: OK
}

// As in the example in https://developers.google.com/places/documentation/actions#PlaceReports
func ExampleService_Execute_forPOST() {
	s, err := googlemaps.NewService("place/add", "POST", googlemaps.GOOGLEMAPS_AUTH_KEY)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.AddParam("sensor", "false")
	s.SetApiKey("TEST_API_KEY")

	// We pass a JSON string as the request body
	resp, err := s.Execute("{\"location\": {\"lat\": -33.8669710, \"lng\": 151.1958750}, \"accuracy\": 50, \"name\": \"Google Shoes!\", \"types\": [\"shoe_store\"], \"language\": \"en-AU\"}")

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//Output: REQUEST_DENIED, since the Api Key is not valid
}
