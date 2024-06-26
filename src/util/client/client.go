package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type NetClientRequest struct {
	NetClient  *http.Client
	RequestUrl string
	QueryParam []QueryParams
}

type QueryParams struct {
	Param string
	Value string
}

type Response struct {
	Res    []byte
	Err    error
	Status int
}

func (ncr *NetClientRequest) AddQueryParam(param, value string) {
	ncr.QueryParam = append(ncr.QueryParam, QueryParams{Param: param, Value: value})
}

func (ncr *NetClientRequest) Get(channel chan Response) {
	url := ncr.RequestUrl
	if len(ncr.QueryParam) > 0 {
		url += "?"
		for _, param := range ncr.QueryParam {
			url += param.Param + "=" + param.Value + "&"
		}
		url = url[:len(url)-1]
	}

	bResp, err := ncr.NetClient.Get(url)
	if err != nil {
		channel <- Response{Err: err}
		return
	}
	defer bResp.Body.Close()

	resBody, err := io.ReadAll(bResp.Body)
	if err != nil {
		channel <- Response{Err: err}
		return
	}

	channel <- Response{Res: resBody, Status: bResp.StatusCode}
}

func (ncr *NetClientRequest) Post(load interface{}, channel chan Response) {
	go func() {
		marshalled, err := json.Marshal(load)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		// Construct URL with query parameters
		urlObj, err := url.Parse(ncr.RequestUrl)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		if len(ncr.QueryParam) > 0 {
			query := urlObj.Query()
			for _, param := range ncr.QueryParam {
				query.Add(param.Param, param.Value)
			}
			urlObj.RawQuery = query.Encode()
		}

		// Create a new POST request
		req, err := http.NewRequest("POST", urlObj.String(), bytes.NewBuffer(marshalled))
		if err != nil {
			channel <- Response{Err: err}
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Perform the request
		bResp, err := ncr.NetClient.Do(req)
		if err != nil {
			channel <- Response{Err: err}
			return
		}
		defer bResp.Body.Close()

		// Read the response body
		respBody, err := io.ReadAll(bResp.Body)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		channel <- Response{Res: respBody, Status: bResp.StatusCode}
	}()
}

func Put(netClient *http.Client, uri string, load interface{}, channel chan Response) {
	go func() {
		marshalledLoad, err := json.Marshal(load)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		req, err := http.NewRequest(http.MethodPut, uri, bytes.NewBuffer(marshalledLoad))
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := netClient.Do(req)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		channel <- Response{
			Res:    respBody,
			Status: resp.StatusCode,
		}
	}()
}

func Delete(netClient *http.Client, uri string, load interface{}, channel chan Response) {
	go func() {
		marshalledLoad, err := json.Marshal(load)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		req, err := http.NewRequest(http.MethodDelete, uri, bytes.NewBuffer(marshalledLoad))
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := netClient.Do(req)
		if err != nil {
			channel <- Response{Err: err}
			return
		}

		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		channel <- Response{
			Res:    respBody,
			Status: resp.StatusCode,
		}
	}()
}
