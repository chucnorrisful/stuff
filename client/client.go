package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var httpClient = &http.Client{}

type SimpleClient struct {
	BaseUrl string
	Headers map[string]string
}

func noRedirect(req *http.Request, via []*http.Request) error {
	_, _ = spew.Println("redirect catched.")
	return http.ErrUseLastResponse
}

//jsonEncoded body
//func (simCl SimpleClient) Call(method, path string, queryParams, bodyParams, v interface{}) error {
//
//	// Check if the BaseUrl has been set.
//	if simCl.BaseUrl == "" {
//		return errors.New("simpleClient: BaseUrl has not been set")
//	}
//
//	// Build the API url using the given endpoint.
//	u := simCl.BaseUrl + path
//
//	// Handle the query parameters.
//	if queryParams != nil {
//		q, err := query.Encode(queryParams)
//		if err != nil {
//			return err
//		}
//
//		u += "?" + q
//	}
//
//	// Handle the body parameters.
//	b := new(bytes.Buffer)
//	if bodyParams != nil {
//		if err := json.NewEncoder(b).Encode(bodyParams); err != nil {
//			return err
//		}
//	}
//
//	// Build the Request.
//	req, err := http.NewRequest(method, u, b)
//	if err != nil {
//		return err
//	}
//
//	// Set headers.
//
//	if len(simCl.Headers) > 0 {
//		for k, v := range simCl.Headers {
//			req.Header.Add(k, v)
//		}
//	}
//	//req.SetBasicAuth("", key)
//
//	// Send request.
//	httpClient.Timeout = 16 * time.Minute
//	resp, err := httpClient.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//
//	// Handle API error. may be broke
//	if resp.StatusCode >= 400 {
//		var apiErr = errors.New(strconv.Itoa(resp.StatusCode))
//		return apiErr
//	}
//
//	if v != nil {
//		return json.NewDecoder(resp.Body).Decode(v)
//	}
//
//	return nil
//}

//uses map instead of struct for paras, urlEncodes values, jsonMarshals bodyParams
func (simCl SimpleClient) Call2(method, path string, queryParams map[string]string, bodyParams, v interface{}) error {

	// Check if the BaseUrl has been set.
	if simCl.BaseUrl == "" {
		return errors.New("simpleClient: BaseUrl has not been set")
	}

	// Build the API url using the given endpoint.
	u := simCl.BaseUrl + path

	// Handle the query parameters.
	if queryParams != nil {

		u += "?"
		for k, v := range queryParams {
			u += k + "=" + url.PathEscape(v) + "&"
		}

		strings.TrimRight(u, "&")
	}

	// Handle the body parameters.
	b := new(bytes.Buffer)
	if bodyParams != nil {
		if err := json.NewEncoder(b).Encode(bodyParams); err != nil {
			return err
		}
	}

	// Build the Request.
	req, err := http.NewRequest(method, u, b)
	if err != nil {
		return err
	}

	// Set headers.

	if len(simCl.Headers) > 0 {
		for k, v := range simCl.Headers {
			req.Header.Add(k, v)
		}
	}
	//req.SetBasicAuth("", key)

	// Send request.
	httpClient.Timeout = 16 * time.Minute
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Handle API error. may be broke
	if resp.StatusCode >= 400 {
		var apiErr = errors.New(strconv.Itoa(resp.StatusCode))
		//json.NewDecoder(resp.Body).Decode(v)
		return apiErr
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	_ = resp.Body.Close()
	return nil
}

// queryParas and bodyParas both are taken as maps and are urlEncoded. header x-www-form-urlencoded should be set!
func (simCl SimpleClient) Call3(method, path string, queryParams, bodyParams map[string]string, v interface{}) error {

	// Check if the BaseUrl has been set.
	if simCl.BaseUrl == "" {
		return errors.New("simpleClient: BaseUrl has not been set")
	}

	// Build the API url using the given endpoint.
	u := simCl.BaseUrl + path

	// Handle the query parameters.
	if queryParams != nil {

		u += "?"
		for k, v := range queryParams {
			u += k + "=" + url.PathEscape(v) + "&"
		}

		strings.TrimRight(u, "&")
	}

	// Handle the body parameters.
	bString := ""
	if bodyParams != nil {
		for k, v := range bodyParams {
			bString += k + "=" + url.PathEscape(v) + "&"
		}

		strings.TrimRight(bString, "&")
	}
	b := strings.NewReader(bString)

	// Build the Request.
	req, err := http.NewRequest(method, u, b)
	if err != nil {
		return err
	}

	// Set headers.

	if len(simCl.Headers) > 0 {
		for k, v := range simCl.Headers {
			req.Header.Add(k, v)
		}
	}
	//req.SetBasicAuth("", key)

	// Send request.
	httpClient.Timeout = 2 * time.Minute
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Handle API error. may be broke
	if resp.StatusCode >= 400 {
		spew.Dump("Walhalla")
		var apiErr = errors.New(strconv.Itoa(resp.StatusCode))
		_ = json.NewDecoder(resp.Body).Decode(v)
		return apiErr
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}

// queryParas are taken as maps and are urlEncoded. body is string of type text/xml. dont forget to set text/xml header.
// also response is []byte, no longer interface.
func (simCl SimpleClient) Call4(method, path string, queryParams map[string]string, bodyParams string, v *[]byte) error {

	// Check if the BaseUrl has been set.
	if simCl.BaseUrl == "" {
		return errors.New("simpleClient: BaseUrl has not been set")
	}

	// Build the API url using the given endpoint.
	u := simCl.BaseUrl + path

	// Handle the query parameters.
	if queryParams != nil {

		u += "?"
		for k, v := range queryParams {
			u += k + "=" + url.PathEscape(v) + "&"
		}

		strings.TrimRight(u, "&")
	}

	// Handle the body parameters.
	b := strings.NewReader(bodyParams)

	// Build the Request.
	req, err := http.NewRequest(method, u, b)
	if err != nil {
		return err
	}

	// Set headers.

	if len(simCl.Headers) > 0 {
		for k, v := range simCl.Headers {
			req.Header.Add(k, v)
		}
	}
	//req.SetBasicAuth("", key)

	// Send request.
	httpClient.Timeout = 16 * time.Minute
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	// Handle API error. may be broke
	if resp.StatusCode >= 400 {
		var apiErr = errors.New(strconv.Itoa(resp.StatusCode))
		*v, _ = ioutil.ReadAll(resp.Body)
		return apiErr
	}

	if v != nil {
		*v, _ = ioutil.ReadAll(resp.Body)
		return nil
	}

	_ = resp.Body.Close()
	return nil
}

// with Cookies
// queryParas arenc taken as maps and are urlEoded. body is string of type text/xml. dont forget to set text/xml header.
// also response is []byte, no longer interface.
func (simCl SimpleClient) Call5(method, path string, queryParams map[string]string, bodyParams map[string]string, cs []*http.Cookie, v *[]byte) error {

	// Check if the BaseUrl has been set.
	if simCl.BaseUrl == "" {
		return errors.New("simpleClient: BaseUrl has not been set")
	}

	// Build the API url using the given endpoint.
	u := simCl.BaseUrl + path

	// Handle the query parameters.
	if queryParams != nil {

		u += "?"
		for k, v := range queryParams {
			u += k + "=" + url.PathEscape(v) + "&"
		}

		strings.TrimRight(u, "&")
	}

	// Handle the body parameters.
	bodyHelp := ""
	for k, v := range bodyParams {
		bodyHelp += k + "=" + url.PathEscape(v) + "&"
	}
	strings.TrimRight(bodyHelp, "&")
	b := strings.NewReader(bodyHelp)

	// Build the Request.
	req, err := http.NewRequest(method, u, b)
	if err != nil {
		return err
	}

	// Set headers.

	if len(simCl.Headers) > 0 {
		for k, v := range simCl.Headers {
			req.Header.Add(k, v)
		}
	}
	//req.SetBasicAuth("", key)

	//add cookies ;)
	for _,c := range cs {
		req.AddCookie(c)
	}

	//rd, _ := req.GetBody()
	//spew.Dump(ioutil.ReadAll(rd))

	// Send request.
	httpClient.Timeout = 16 * time.Minute
	resp, err := httpClient.Do(req)
	if err != nil {
		spew.Dump(resp.Cookies())

		return err
	}

	// Handle API error. may be broke
	if resp.StatusCode >= 400 {
		var apiErr = errors.New(strconv.Itoa(resp.StatusCode))
		*v, _ = ioutil.ReadAll(resp.Body)
		return apiErr
	}

	if v != nil {
		*v, _ = ioutil.ReadAll(resp.Body)
		return nil
	}

	_ = resp.Body.Close()
	return nil
}
