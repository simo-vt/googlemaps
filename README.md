googlemaps
==========

Installation
------------

Simply tell Go to do the following thing:

	go get github.com/simo-vt/googlemaps

How to use
------------

#### Example:
    package main
    
    import (
    	"fmt"
    	"googlemaps"
    )
    
    func main() {
        // We instantiate a service from the Google API and select the authentication method. It uses one of the available constants, explained below.
    	s, err := googlemaps.NewService("elevation", "GET", googlemaps.GOOGLEMAPS_AUTH_EMPTY)
    	if err != nil {
    		fmt.Println(err.Error())
    		return
    	}
        
        // We add parameters to the service.
    	s.AddParam("sensor", "false")
    	s.AddParam("locations", "40.714728,-73.998672")
    
        // And we execute with a request body. There are cases when we will need to pass JSON code.
    	resp, err := s.Execute("")
    
    	if err != nil {
    		fmt.Println(err.Error())
    		return
    	}
    
    	fmt.Println(resp["status"])
    	//Output: OK
    }
    
### Managing parameters

AddParam(key, value string) adds a request parameter to the service request. Note that you should always use this function to set the "sensor" parameter to either "true" or "false"

RemoveParam(key string) removes a parameter.

### Authenticating

Almost all of the Google services offer a corporate upgrade for which you will need to authenticate. There are two ways to do this.

#### Setting the authentication type

SetAuthType(aType string) sets the authentication type of the service. Can (and should) be one of these:
* googlemaps.GOOGLEMAPS_AUTH_CLIENT
* googlemaps.GOOGLEMAPS_AUTH_KEY
* googlemaps.GOOGLEMAPS_AUTH_EMPTY

#### API Key
You should use SetApiKey(aKey string) if you have used the "key" authentication type, you need to set the API Key with this function before calling Execute.

#### Client
Use SetClientId(cId string) and SetClientSecret(cSecret string) if you have used the "client" authentication type. You need to set the Client ID and Client Secret with this function before calling Execute.

Note that the client secret is a base-64 encoded string for URL's. Basically, the one you get in the e-mail after you register a business account with Google Maps. An example key is "vNIXE0xscrmjlyV-12Nj_BvUPaw="

More information here - https://developers.google.com/maps/documentation/business/webservices/auth#generating_valid_signatures

License
----------
Copyright (c) 2014 Simeon Totev. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Simeon Totev nor the names of his
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

