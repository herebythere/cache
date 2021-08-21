package supercachelx

import (
	"encoding/json"
	"os"
	"testing"
)

const (
	dayInSeconds    = 86400
	emptyStr        = ""
	expCache        = "EX"
	getCache        = "GET"
	hello           = "hello"
	helloWorld      = "hello_world"
	mgetCache       = "MGET"
	setCache        = "SET"
	testExecSetJson = "test_supercache_exec_set_json"
	testExecSetStr  = "test_exec_set_string"
	testidA         = "supercachelx_test_id_a"
	testidB         = "supercachelx_test_id_b"
	testIncrement   = "INCR"
	testSCLX        = "test_supercache_local_interface"
	world           = "world"
)

var (
	localCacheAddress = os.Getenv("LOCAL_CACHE_ADDRESS")
	// localCacheAddress = "http://192.168.1.32:1234"
)

func TestExecInstructionsAndParseInt64(t *testing.T) {
	instructions := []interface{}{testIncrement, testSCLX}
	count, errCount := ExecInstructionsAndParseInt64(
		localCacheAddress,
		&instructions,
	)
	if errCount != nil {
		t.Fail()
		t.Logf(errCount.Error())
	}
	if count == nil {
		t.Fail()
		t.Logf("increment was not successfuul")
	}
	if count != nil && *count < 1 {
		t.Fail()
		t.Logf("increment was less than 1, which means key might be occupied by non integer")
	}
}

func TestExecInstructionsAndParseMultipleInt64(t *testing.T) {
	headID := getCacheSetID(testSCLX, testidA)
	tailID := getCacheSetID(testSCLX, testidB)

	instructions := []interface{}{mgetCache, headID, tailID}
	counts, errCount := ExecInstructionsAndParseMultipleInt64(
		localCacheAddress,
		&instructions,
	)
	if errCount != nil {
		t.Fail()
		t.Logf(errCount.Error())
	}
	if counts == nil {
		t.Fail()
		t.Logf("multiple gets was not successful")
	}
	if counts != nil && len(*counts) < 2 {
		t.Fail()
		t.Logf("less than two counts were returned from instructions")
	}
}

func TestExecGetInstructionsAndParseString(t *testing.T) {
	setInstructions := []interface{}{setCache, testExecSetStr, helloWorld, expCache, dayInSeconds}
	ExecInstructionsAndParseString(localCacheAddress, &setInstructions)

	getInstructions := []interface{}{getCache, testExecSetStr}
	parsedStr, errParsedStr := ExecInstructionsAndParseString(
		localCacheAddress,
		&getInstructions,
	)
	if errParsedStr != nil {
		t.Fail()
		t.Logf(errParsedStr.Error())
	}
	if parsedStr == nil {
		t.Fail()
		t.Logf("fetching string was not successful")
	}
	if parsedStr != nil && *parsedStr == emptyStr {
		t.Fail()
		t.Logf("empty string was returned")
	}
}

func TestExecGetInstructionsAndParseBase64(t *testing.T) {
	helloWorldArray := []string{hello, world}
	helloWorldJSONBytes, errHelloWorldJSONBytes := json.Marshal(helloWorldArray)
	if errHelloWorldJSONBytes != nil {
		t.Fail()
		t.Logf(errHelloWorldJSONBytes.Error())
	}
	helloWorldJSONStr := string(helloWorldJSONBytes)

	setInstructions := []interface{}{setCache, testExecSetJson, helloWorldJSONStr, expCache, dayInSeconds}
	ExecInstructionsAndParseString(localCacheAddress, &setInstructions)

	getInstructions := []interface{}{getCache, testExecSetJson}
	parsedStr, errParsedStr := ExecInstructionsAndParseBase64(
		localCacheAddress,
		&getInstructions,
	)
	if errParsedStr != nil {
		t.Fail()
		t.Logf(errParsedStr.Error())
	}
	if parsedStr == nil {
		t.Fail()
		t.Logf("fetch base64 failed")
	}
	if parsedStr != nil && *parsedStr == emptyStr {
		t.Fail()
		t.Logf("empty string was returned")
	}

	var jsonResult []string
	errJSON := json.Unmarshal([]byte(*parsedStr), &jsonResult)
	if errJSON != nil {
		t.Fail()
		t.Logf(errJSON.Error())
	}

	if len(jsonResult) != 2 {
		t.Fail()
		t.Logf("expected a string array with a length of 2")
	}

	if len(jsonResult) == 2 && jsonResult[0] != hello {
		t.Fail()
		t.Logf("expected first index to contain 'hello'")
	}
	if len(jsonResult) == 2 && jsonResult[1] != world {
		t.Fail()
		t.Logf("expected second index to contain 'world'")
	}
}
