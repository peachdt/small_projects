package utils

import (
	"BoomPayments/cs/core_v0/common"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
)

// TODO get this from config
var verbose = false


func HttpJsonPost(url string, json_body interface{}, headers map[string]string, json_resp interface{} ) (*http.Response, error) {
	return _http_json_upload("POST", url, json_body, headers, json_resp)
}

func HttpJsonPut(url string, json_body interface{}, headers map[string]string, json_resp interface{} ) (*http.Response, error) {
	return _http_json_upload("PUT", url, json_body, headers, json_resp)
}

func HttpGet(url string, headers map[string]string, json_resp interface{}) (int, error){
	return _http_request("GET", url, headers, json_resp)
}

func HttpDelete(url string, headers map[string]string, json_resp interface{}) (int, error){
	return _http_request("DELETE", url, headers, json_resp)
}

func _http_request(method, url string, headers map[string]string, json_resp interface{}) (int, error) {

	// TODO Only print below garbage if in global verbose mode
	var req *http.Request

	req, _ = http.NewRequest(method, url, nil)
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	var client *http.Client
	// return client, req
	client = &http.Client{}
	//    var json_resp []map[string]interface{}

	resp, err := client.Do(req)


	if err != nil {
		panic(fmt.Sprintf("failed on http req to (%v) with error (%v)", url, err))
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		common.LogWarning("Received no content (204) from big commerce")
		return resp.StatusCode , err
	}
	t := reflect.TypeOf(json_resp)
	if t == nil { // They passed in nil for json_resp meaning they dont care about req resp
		return resp.StatusCode, nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if verbose {
		common.LogInfo("**********************")
		for k, v := range resp.Header {
			common.LogInfo(fmt.Sprintf("%v: %v", k, v))
		}
		common.LogInfo(strconv.Itoa(resp.StatusCode))

		common.LogInfo("**********************")
		common.LogInfo(string(body))
	}
	errorr := json.Unmarshal(body, json_resp)
	if errorr != nil {
		common.LogWarningf(fmt.Sprintf("Cannot unmarshal! err: [%v]", errorr))
	}
	return resp.StatusCode, errorr
}


func _http_json_upload(method, url string, json_body interface{}, headers map[string]string, json_resp interface{} ) (*http.Response, error) {
	var req *http.Request
	b, err := json.Marshal(json_body)
	if err != nil {
		panic("cannot marshal json_body")
	}
	req, _ = http.NewRequest(method, url, bytes.NewBuffer([]byte(b)))
	req.Header.Set("Content-Type", "application/json")

	for header, value := range headers {
		req.Header.Set(header, value)
	}
	var client *http.Client
	client = &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(fmt.Sprintf("failed on http req to (%v) with error (%v)", url, err))
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		common.LogWarning("Received no content (204) from big commerce")
		return resp, err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	if verbose {
		common.LogInfo("**********************")
		for k, v := range resp.Header {
			common.LogInfo(fmt.Sprintf("%v: %v", k, v))
		}
		common.LogInfo(strconv.Itoa(resp.StatusCode))
		common.LogInfo("**********************")
		common.LogInfo(string(body))
	}
	t := reflect.TypeOf(json_resp)
	if t == nil { // They passed in nil for json_resp meaning they dont care about req resp
		return resp, nil
	}
	errorr := json.Unmarshal(body, json_resp)
	if errorr != nil {
		common.LogWarning(fmt.Sprintf("Cannot unmarshal! err: [%v]", errorr))
	}
	return resp, errorr

}
