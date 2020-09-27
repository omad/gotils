package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func authPVOutputRequest(req *http.Request) {
	apikey, ok1 := os.LookupEnv("PVOUTPUT_APIKEY")
	systemid, ok2 := os.LookupEnv("PVOUTPUT_SYSTEMID")
	if !ok1 || !ok2 {
		fmt.Printf("Error: Must provide PVOUTPUT_APIKEY and PVOUTPUT_SYSTEMID")
		panic()
	}
	req.Header.Add("X-Pvoutput-Apikey", apikey)
	req.Header.Add("X-Pvoutput-SystemId", systemid)
}

func findMissingPVOutputDates(start string, end string) []string {
	url := "https://pvoutput.org/service/r2/getmissing.jsp"
	req, err := http.NewRequest("GET", url, nil)
	check(err)

	authPVOutputRequest(req)

	params := req.URL.Query()
	params.Add("df", start)
	params.Add("dt", end)
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	check(err)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	fmt.Println("get:\n", string(body))

	missing := strings.Split(string(body), ",")
	fmt.Println(missing)
	return missing
}

func makeUploadString(outputs map[string]int, missing []string) string {
	var toUpload []string
	for _, missingDate := range missing {
		toUpload = append(toUpload, fmt.Sprintf("%s,%d", missingDate, outputs[missingDate]))
	}
	if len(toUpload) > 30 {
		panic("Too many days of data.")
	}

	return strings.Join(toUpload, ";")
}

// https://pvoutput.org/help.html#api-addbatchoutput
func batchSubmitOutputs(uploadData string) {
	url := "https://pvoutput.org/service/r2/addbatchoutput.jsp"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(uploadData)))
	check(err)

	authPVOutputRequest(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// 	exampleData := "20190101,1111;20190202,22222"

func getAuroraOutputs() string {
	auroraCmd := exec.Command("aurora", "-v", "-a", "2", "-Y3", "-l10", "-R1500", "--daily-kwh /dev/ttyUSB0")

	cmdOutput, err := auroraCmd.Output()
	check(err)

	fmt.Println("> ls")
	fmt.Println(string(cmdOutput))
	return string(cmdOutput)
}

func parseAuroraOutputs(s string) map[string]int {
	output := map[string]int{}

	days := strings.Split(s, "\n")
	for _, day := range days {
		words := strings.Fields(day)
		date := words[0]
		if val, err := strconv.ParseFloat(words[1], 32); err == nil {
			fmt.Println(date, val)
			// fmt.Printf("%T, %v\n", s, s)
			output[date] = int(val * 1000)
		}
	}
	fmt.Println(output)
	return output
}
func setOperations() {
	// See https://emersion.fr/blog/2017/sets-in-go/
	// And https://stackoverflow.com/a/34020023
	s1 := map[int]bool{5: true, 2: true}
	_, ok := s1[6] // check for existence
	fmt.Println(ok)
	s1[8] = true  // add element
	delete(s1, 2) // remove element

	s2 := map[int]bool{}
	// Union
	sUnion := map[int]bool{}
	for k := range s1 {
		sUnion[k] = true
	}
	for k := range s2 {
		sUnion[k] = true
	}

	// Intersection

	sIntersection := map[int]bool{}
	for k := range s1 {
		if s2[k] {
			sIntersection[k] = true
		}
	}
}

func main() {
	getAuroraOutputs()
	findMissingPVOutputDates("20190101", "20191231")
}
