package supercachelx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	applicationJSON = "application/json"
	colonDelimiter  = ":"
)

var (
	errInstructionsAreNil     = errors.New("instructions are nil")
	errRequestToCacheFailed   = errors.New("request to local cache failed")
	errInstructionsFailed     = errors.New("instructions failed to exec")
	errRequestFailedToResolve = errors.New("request failed to resolve instructions")
)

func getCacheSetID(categories ...string) string {
	return strings.Join(categories, colonDelimiter)
}

func fetchPostRequest(
	cacheAddress string,
	instructions *[]interface{},
) (
	*http.Response,
	error,
) {
	bodyBytes := new(bytes.Buffer)
	errJson := json.NewEncoder(bodyBytes).Encode(*instructions)
	if errJson != nil {
		return nil, errJson
	}

	return http.Post(cacheAddress, applicationJSON, bodyBytes)
}

func execInstructionsAndParseInt64(
	cacheAddress string,
	instructions *[]interface{},
) (
	*int64,
	error,
) {
	if instructions == nil {
		return nil, errInstructionsAreNil
	}

	resp, errResp := fetchPostRequest(
		cacheAddress,
		instructions,
	)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errRequestToCacheFailed
	}

	var count int64
	errCount := json.NewDecoder(resp.Body).Decode(&count)

	return &count, errCount
}

func execInstructionsAndParseString(
	cacheAddress string,
	instructions *[]interface{},
) (
	*string,
	error,
) {
	if instructions == nil {
		return nil, errInstructionsAreNil
	}

	resp, errResp := fetchPostRequest(
		cacheAddress,
		instructions,
	)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errRequestFailedToResolve
	}

	var respBodyAsStr string
	errJSONResponse := json.NewDecoder(resp.Body).Decode(&respBodyAsStr)

	return &respBodyAsStr, errJSONResponse
}

func execInstructionsAndParseMultipleInt64(
	cacheAddress string,
	instructions *[]interface{},
) (
	*[]int64,
	error,
) {
	if instructions == nil {
		return nil, errInstructionsAreNil
	}

	resp, errResp := fetchPostRequest(
		cacheAddress,
		instructions,
	)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errRequestFailedToResolve
	}

	var respBodyAsStrArr []string
	errJSONResponse := json.NewDecoder(resp.Body).Decode(&respBodyAsStrArr)
	if errJSONResponse != nil {
		return nil, errJSONResponse
	}

	respBodyAsInts := make([]int64, len(respBodyAsStrArr))
	for index, base64IntStr := range respBodyAsStrArr {
		intBytes, errIntBytes := base64.StdEncoding.DecodeString(base64IntStr)
		if errIntBytes != nil {
			return nil, errIntBytes
		}

		intStr := string(intBytes)
		if intStr == "" {
			respBodyAsInts[index] = 0
			continue
		}

		count, errCount := strconv.ParseInt(intStr, 10, 64)
		if errCount != nil {
			return nil, errCount
		}

		respBodyAsInts[index] = count
	}

	return &respBodyAsInts, nil
}

func execInstructionsAndParseBase64(
	cacheAddress string,
	instructions *[]interface{},
) (
	*string,
	error,
) {
	if instructions == nil {
		return nil, errInstructionsAreNil
	}

	resp, errResp := fetchPostRequest(
		cacheAddress,
		instructions,
	)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errRequestFailedToResolve
	}

	var respBodyAsStr string
	errJSONResponse := json.NewDecoder(resp.Body).Decode(&respBodyAsStr)
	if errJSONResponse != nil {
		return nil, errJSONResponse
	}

	respBodyAsBytes, errRespBodyAsBytes := base64.URLEncoding.DecodeString(
		respBodyAsStr,
	)
	if errRespBodyAsBytes != nil {
		return nil, errRespBodyAsBytes
	}

	respBodyBase64 := string(respBodyAsBytes)

	return &respBodyBase64, nil
}
