package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//-------------------------------------------------------------------------------------------------------
//----------------------------------------------DAO.GO---------------------------------------------------
//-------------------------------------------------------------------------------------------------------

var user = "root"
var password = "myssql#5682"
var dbname = "coronadb"

type allstats struct {
	ActiveCases         int    `json:"ActiveCases"`
	TotalRecovered      int    `json:"TotalRecovered"`
	TotalDeaths         int    `json:"TotalDeaths"`
	TotalCases          int    `json:"TotalCases"`
	Date                string `json:"Date"`
	CurrentDayCases     int    `json:"CurrentDayCases"`
	CurrentDayRecovered int    `json:"CurrentDayRecovered"`
	CurrentDayDeaths    int    `json:"CurrentDayDeaths"`
}

type data struct {
	Name  string   `json:"Name"`
	Stats allstats `json:"Stats"`
}

type worldwideData struct {
	Eachday []data `json:"Eachday"`
}

func getDB() *sql.DB {
	db, err := sql.Open("mysql", user+":"+password+"@/"+dbname)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func getdate() string {
	db := getDB()
	// defer the close till after the function has finished executing
	defer db.Close()
	var date string
	err := db.QueryRow(`SELECT MAX(Date) FROM CountryWiseData;`).Scan(&date)
	if err != nil {
		panic(err.Error())
	}
	return date
}

func getcountrynames() *[]string {
	db := getDB()
	// defer the close till after the function has finished executing
	defer db.Close()
	var countrynames []string
	countries, err := db.Query("SELECT DISTINCT CountryName FROM CountryWiseData")
	for countries.Next() {
		var cname string
		err = countries.Scan(&cname)
		if err != nil {
			panic(err.Error())
		}
		countrynames = append(countrynames, cname)
	}
	return &countrynames
}

func getalldata(from string, to string) *worldwideData {

	db := getDB()
	// defer the close till after the function has finished executing
	defer db.Close()

	var result worldwideData
	//var result data
	var count int
	err := db.QueryRow(`SELECT COUNT(DISTINCT CountryName) From CountryWiseData`).Scan(&count)
	if err != nil {
		panic(err.Error())
	}

	query := fmt.Sprintf(`SELECT CountryName, SUM(ActiveCases), SUM(TotalRecovered), SUM(TotalDeaths), Date
							FROM CountryWiseData where Date Between '%s' And '%s' GROUP BY Date;`, from, to)
	countries1, err := db.Query(query)

	for countries1.Next() {
		var currentdata1 data
		err = countries1.Scan(
			&currentdata1.Name,
			&currentdata1.Stats.ActiveCases,
			&currentdata1.Stats.TotalRecovered,
			&currentdata1.Stats.TotalDeaths,
			&currentdata1.Stats.Date)
		if err != nil {
			panic(err.Error())
		}
		currentdata1.Name = "WorldWide"
		result.Eachday = append(result.Eachday, currentdata1)
	}

	return &result

}

func getallcountrydata(from string, to string) *worldwideData {

	db := getDB()
	// defer the close till after the function has finished executing
	defer db.Close()

	var result worldwideData
	query := fmt.Sprintf(`SELECT CountryName, SUM(ActiveCases), SUM(TotalRecovered), SUM(TotalDeaths), Date
    					  FROM CountryWiseData where Date Between '%s' And '%s' GROUP BY countryname,date;`, from, to)
	countries, err := db.Query(query)

	for countries.Next() {
		var currentdata data
		err = countries.Scan(
			&currentdata.Name,
			&currentdata.Stats.ActiveCases,
			&currentdata.Stats.TotalRecovered,
			&currentdata.Stats.TotalDeaths,
			&currentdata.Stats.Date)
		if err != nil {
			panic(err.Error())
		}
		result.Eachday = append(result.Eachday, currentdata)
	}

	return &result

}

func getcountrydata(country string, from string, to string) *worldwideData {

	db := getDB()
	// defer the close till after the function has finished executing
	defer db.Close()

	var result worldwideData
	query := fmt.Sprintf(`SELECT CountryName, SUM(ActiveCases), SUM(TotalRecovered), SUM(TotalDeaths), Date
    					  FROM CountryWiseData where countryname='%s' and Date Between '%s' And '%s' GROUP BY countryname,date;`, country, from, to)
	countries, err := db.Query(query)

	for countries.Next() {
		var currentdata data
		err = countries.Scan(
			&currentdata.Name,
			&currentdata.Stats.ActiveCases,
			&currentdata.Stats.TotalRecovered,
			&currentdata.Stats.TotalDeaths,
			&currentdata.Stats.Date)
		if err != nil {
			panic(err.Error())
		}
		result.Eachday = append(result.Eachday, currentdata)
	}

	return &result
}

//-------------------------------------------------------------------------------------------------------
//-----------------------------------------SERVICE.GO----------------------------------------------------
//-------------------------------------------------------------------------------------------------------

func worldwideservice(from string, to string) *worldwideData {
	recievedData := getalldata(from, to)
	length := len(recievedData.Eachday)
	for i := 1; i < length; i++ {
		currentDayData := &recievedData.Eachday[i]
		prevDayData := &recievedData.Eachday[i-1]
		currentDayData.Stats.CurrentDayCases = currentDayData.Stats.ActiveCases - prevDayData.Stats.ActiveCases
		currentDayData.Stats.CurrentDayRecovered = currentDayData.Stats.TotalRecovered - prevDayData.Stats.TotalRecovered
		currentDayData.Stats.CurrentDayDeaths = currentDayData.Stats.TotalDeaths - prevDayData.Stats.TotalDeaths
		currentDayData.Stats.TotalCases = currentDayData.Stats.ActiveCases + currentDayData.Stats.TotalRecovered + currentDayData.Stats.TotalDeaths
		//recievedData.Stats.TotalCases += currentCountryData.Stats.TotalCases
	}
	// recievedData.Eachday = remove(recievedData.Eachday, recievedData.Eachday[0])
	return recievedData
}

func allcountryservice(from string, to string) *worldwideData {
	recievedData := getallcountrydata(from, to)
	length := len(recievedData.Eachday)
	for i := 0; i < length; i++ {
		currentDayData := &recievedData.Eachday[i]
		prevDayData := &recievedData.Eachday[i-1]
		currentDayData.Stats.CurrentDayCases = currentDayData.Stats.ActiveCases - prevDayData.Stats.ActiveCases
		currentDayData.Stats.CurrentDayRecovered = currentDayData.Stats.TotalRecovered - prevDayData.Stats.TotalRecovered
		currentDayData.Stats.CurrentDayDeaths = currentDayData.Stats.TotalDeaths - prevDayData.Stats.TotalDeaths
		currentDayData.Stats.TotalCases = currentDayData.Stats.ActiveCases + currentDayData.Stats.TotalRecovered + currentDayData.Stats.TotalDeaths
		//recievedData.Stats.TotalCases += currentCountryData.Stats.TotalCases
	}
	return recievedData
}

func countryservice(country string, from string, to string) *worldwideData {
	recievedData := getcountrydata(country, from, to)
	length := len(recievedData.Eachday)
	for i := 1; i < length; i++ {
		currentDayData := &recievedData.Eachday[i]
		prevDayData := &recievedData.Eachday[i-1]
		currentDayData.Stats.CurrentDayCases = currentDayData.Stats.ActiveCases - prevDayData.Stats.ActiveCases
		currentDayData.Stats.CurrentDayRecovered = currentDayData.Stats.TotalRecovered - prevDayData.Stats.TotalRecovered
		currentDayData.Stats.CurrentDayDeaths = currentDayData.Stats.TotalDeaths - prevDayData.Stats.TotalDeaths
		currentDayData.Stats.TotalCases = currentDayData.Stats.ActiveCases + currentDayData.Stats.TotalRecovered + currentDayData.Stats.TotalDeaths
		//recievedData.Stats.TotalCases += currentCountryData.Stats.TotalCases
	}
	return recievedData
}

func nameservice() *[]string {
	recievedData := getcountrynames()
	return recievedData
}

//-------------------------------------------------------------------------------------------------------
// ----------------------------------controller.go-------------------------------------------------------
//-------------------------------------------------------------------------------------------------------

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Corona API")

}
func worldwidestats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	if from == "" {
		from = "2020-01-23"
	}
	if to == "" {
		to = getdate()
	}
	myDate, err := time.Parse("2006-01-02 15:04", from+" 15:04")
	if err != nil {
		panic(err)
	}
	date := myDate.AddDate(0, 0, -1)
	from = date.Format("2006-01-02")
	log.Println("from:", from)
	log.Println("to:", to)
	var result []data
	result = worldwideservice(from, to).Eachday
	result = result[1:]
	// // returning json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func all(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	if from == "" {
		from = "2020-01-22"
	}
	if to == "" {
		to = getdate()
	}
	log.Println("to:", to)
	result := allcountryservice(from, to)
	// returning json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func countrywise(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	country := vars["countryname"]
	query := r.URL.Query()
	from := query.Get("from")
	to := query.Get("to")
	if from == "" {
		from = "2020-01-23"
	}
	if to == "" {
		to = getdate()
	}
	myDate, err := time.Parse("2006-01-02 15:04", from+" 15:04")
	if err != nil {
		panic(err)
	}
	date := myDate.AddDate(0, 0, -1)
	from = date.Format("2006-01-02")
	log.Println("u r in /country/india/from&to")
	log.Println("to:", to)
	result := countryservice(country, from, to).Eachday
	result = result[1:]
	// returning json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func countrynames(w http.ResponseWriter, r *http.Request) {
	result := nameservice()
	// returning json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

//-------------------------------------------------------------------------------------------------------
//------------------------------------------------- server.go--------------------------------------------
//-------------------------------------------------------------------------------------------------------

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", root)
	rtr.HandleFunc("/home", worldwidestats).Methods("GET")
	rtr.HandleFunc("/all", all).Methods("GET")
	rtr.HandleFunc("/countrynames", countrynames).Methods("GET")
	rtr.HandleFunc("/country/{countryname}", countrywise).Methods("GET")
	http.Handle("/", rtr)
	log.Println("server starting at 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
