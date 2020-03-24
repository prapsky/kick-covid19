package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"kick-covid19/models"
)

func logerr(n int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}

func GetCCandDate(w http.ResponseWriter, r *http.Request) (string, string) {
	param := r.URL.Path[16:]
	index := strings.Index(param, "/")
	static := param[index+1:]
	if len(static) == 0 {
		http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusMovedPermanently)
		return "", ""
	}

	param = path.Clean(param)
	cc := filepath.Dir(param)

	inputDate := filepath.Base(param)
	addLastTime := "23:00:00"
	inputDateTime := inputDate + " " + addLastTime

	return cc, inputDateTime
}

func DataByCountryandDate(w http.ResponseWriter, cc string, inputDateTime string) models.Response {
	url := "http://covid19.soficoop.com/country/" + cc
	req, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.Response{}
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.Response{}
	}

	request := models.Request{}
	logerr(1, json.Unmarshal(body, &request))

	selectedIndex := 0
	if inputDateTime < (request.Snapshots[0].Timestamp).Format("2006-01-02 15:04:05") {
		selectedIndex = 0
	}

	lengthSnapshots := len(request.Snapshots)
	for i := 0; i < lengthSnapshots; i++ {
		r, _ := regexp.MatchString(inputDateTime, (request.Snapshots[i].Timestamp).Format("2006-01-02 15:04:05"))
		if r {
			selectedIndex = i
		}
	}

	if inputDateTime > (request.Snapshots[lengthSnapshots-1].Timestamp).Format("2006-01-02 15:04:05") {
		selectedIndex = lengthSnapshots - 1
	}

	response := models.Response{}
	response.Timestamp = request.Snapshots[selectedIndex].Timestamp.Format("2006-01-02")
	response.CountryName = request.CountryName
	response.Cases = request.Snapshots[selectedIndex].Cases
	response.Deaths = request.Snapshots[selectedIndex].Deaths
	response.Recovered = request.Snapshots[selectedIndex].Recovered

	return response
}

func GetRealtimeData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		cc, inputDateTime := GetCCandDate(w, r)

		response := DataByCountryandDate(w, cc, inputDateTime)

		jsonInBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logerr(w.Write(jsonInBytes))
	default:
		w.WriteHeader(http.StatusNotFound)
		logerr(w.Write([]byte(`{"message": "not found"}`)))
	}
}
