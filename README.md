googlemaps
==========

Introduction
------------
Google Maps WebServices HTTP Client API for the Go Programming Language.

The idea is simple. You create a service, you assign parameters to it and you execute it with the Execute function.

This client API supports the two types of business authentication, described here - https://developers.google.com/maps/documentation/business/webservices/auth

For "key" authentication you must set the API Key for the service with the function
* func (s *Service) SetApiKey(aKey string)

For "client" authentication you must set a Client ID and Client Secret with the functions
* func (s *Service) SetClientId(cId string)
* func (s *Service) SetClientSecret(cSecret string)

Official API documentation - https://developers.google.com/maps/documentation/webservices/. Please note that you should familiarize yourself with the terms of Google Maps API before using it.

Constants
---------
* googlemaps.GOOGLEMAPS_AUTH_CLIENT - used if we want to use "client" type of authentication. For business users. Note the limits explained in the Google Maps API.
* googlemaps.GOOGLEMAPS_AUTH_KEY - used if we want to use "key" type of authentication. For business users. Note the limits explained in the Google Maps API.
* googlemaps.GOOGLEMAPS_AUTH_EMPTY - used if we do not want to authenticate. Note the limits explained in the Google Maps API.

Functions
---------

### func NewService(sPath string, rMethod string, aType string) (*Service, error)
NewService creates a new Google Maps API service. The passed parameters are:

* sPath string - this is the name of the service you wish to call. For example: "place/event/add" or "directions".
* rMethod string - this is the request method. You will most often use GET, but with the case of the Places API you may need to use POST.
* aType string - can be one of the following: googlemaps.GOOGLEMAPS_AUTH_KEY, googlemaps.GOOGLEMAPS_AUTH_CLIENT, googlemaps.GOOGLEMAPS_AUTH_EMPTY

### func (s *Service) Execute(requestBody string) (map[string]interface{}, error)
Execute is the main function of this client API. You should pass the request body as a string.

In most cases, this should be an empty string, since the requests to the API are mostly GET. In the case of the Places API the requestBody string should contain a JSON object with the parameters, like in this link: https://developers.google.com/places/documentation/actions#event_add

#### Example:
    package main
    
    import (
    	"fmt"
    	"net/googlemaps"
    )
    
    func main() {
    	s, err := googlemaps.NewService("elevation", "GET", googlemaps.GOOGLEMAPS_AUTH_EMPTY)
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

### func (s *Service) AddParam(key, value string)
AddParam adds a request parameter to the service request. Note that you should always use this function to set the "sensor" parameter to either "true" or "false"

### func (s *Service) RemoveParam(key string)
Removes a parameter.

### func (s *Service) SetApiKey(aKey string)
If you have used the "key" authentication type, you need to set the API Key with this function before calling Execute.

### func (s *Service) SetAuthType(aType string) error
SetAuthType sets the authentication type of the service. Can (and should) be one of these:
* googlemaps.GOOGLEMAPS_AUTH_CLIENT
* googlemaps.GOOGLEMAPS_AUTH_KEY
* googlemaps.GOOGLEMAPS_AUTH_EMPTY

### func (s *Service) SetClientId(cId string)
If you have used the "client" authentication type, you need to set the Client ID with this function before calling Execute.

### func (s *Service) SetClientSecret(cSecret string)
If you have used the "client" authentication type, you need to set the Client Secret with this function before calling Execute.

Note that the client secret is a base-64 encoded string for URL's. Basically, the one you get in the e-mail after you register a business account with Google Maps. An example key is "vNIXE0xscrmjlyV-12Nj_BvUPaw="

More information here - https://developers.google.com/maps/documentation/business/webservices/auth#generating_valid_signatures

### func (s *Service) SetRequestMethod(rMethod string) error
SetRequestMethod sets the request method of the service.

In most cases, you will be using GET, but in the case of the Places API, you might need to use POST.

### func (s *Service) SetServicePath(sPath string)
SetServicePath sets the path of the service.

Example paths are: "place/event/add" or "directions"

### func (s *Service) String() string
String outputs the service method and URL as a string.
