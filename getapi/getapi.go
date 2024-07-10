package getapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Api struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []string            `json:"locations"`
	ConcertDates []string            `json:"concertDates"`
	Relation     map[string][]string `json:"datesLocations"`
}

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location     []string `json:"locations"`
	Date         string   `json:"concertDates"`
	Relations    string   `json:"datesLocations"`
}

type Locations struct {
	Index []LocIndex `json:"index"`
}
type LocIndex struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}
type Dates struct {
	Index []DateIndex `json:"index"`
}
type DateIndex struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type Relations struct {
	Index []RelIndex `json:"index"`
}
type RelIndex struct {
	ID       int                 `json:"id"`
	Relation map[string][]string `json:"datesLocations"`
}

var (
	artists    []Artist
	artistsall []Api
	locations  Locations
	dates      Dates
	relations  Relations
)

const url = "https://groupietrackers.herokuapp.com/api"

// recreate a structure to make all data usable
func GetFullData() (error, int) {
	err1 := GetArtists() // check that "url + /artist" does not exist, same for the 3 other functions.
	err2 := GetDates()
	err3 := GetLocations()
	err4 := GetRelations()
	limit := 0
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("ERROR TO GET DATAS, ARTISTS, LOCATIONS, DATES"), limit
	}

	// Range Loop to assign all values of the api into defined structures.
	for i := range artists {
		var artistTemplate Api
		artistTemplate.ID = i + 1
		artistTemplate.Image = artists[i].Image
		artistTemplate.Name = artists[i].Name
		artistTemplate.Members = artists[i].Members
		artistTemplate.CreationDate = artists[i].CreationDate
		artistTemplate.FirstAlbum = artists[i].FirstAlbum
		artistTemplate.Locations = locations.Index[i].Locations
		artistTemplate.ConcertDates = dates.Index[i].Dates
		artistTemplate.Relation = relations.Index[i].Relation
		artistsall = append(artistsall, artistTemplate)
		limit = i + 1
	}
	return nil, limit // these outputs check if function is perfectly executed.
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
		return
	}
	err2 := GetArtists()
	if err2 != nil {
		return
	}

	// get creation date filter input
	CreationDateFrom := r.FormValue("createdFrom")
	CreationDateTill := r.FormValue("createdTill")

	//get number of members filter input
	var m1, m2, m3, m4, m5, m6, m7, m8 int
	var err1, errZ, err3, err4, err5, err6, err7, err8 error
	m1, err1 = strconv.Atoi(r.FormValue("m1"))
	if err1 != nil {
		m1 = 0
	}
	m2, errZ = strconv.Atoi(r.FormValue("m2"))
	if errZ != nil {
		m2 = 0
	}
	m3, err3 = strconv.Atoi(r.FormValue("m3"))
	if err3 != nil {
		m3 = 0
	}
	m4, err4 = strconv.Atoi(r.FormValue("m4"))
	if err4 != nil {
		m4 = 0
	}
	m5, err5 = strconv.Atoi(r.FormValue("m5"))
	if err5 != nil {
		m5 = 0
	}
	m6, err6 = strconv.Atoi(r.FormValue("m6"))
	if err6 != nil {
		m6 = 0
	}
	m7, err7 = strconv.Atoi(r.FormValue("m7"))
	if err7 != nil {
		m7 = 0
	}
	m8, err8 = strconv.Atoi(r.FormValue("m8"))
	if err8 != nil {
		m8 = 0
	}
	members := []int{m1, m2, m3, m4, m5, m6, m7, m8}
	sum := 0
	for _, n := range members {
		sum += n
	}

	// get first album filter input
	firstAlbumFrom := r.FormValue("dateFrom")
	firstAlbumTill := r.FormValue("dateTill")

	if CreationDateFrom != "" || CreationDateTill != "" {
		if CreationDateFrom == "" {
			CreationDateFrom = "1900"
		}
		if CreationDateTill == "" {
			CreationDateTill = "2020"
		}
		artists = CreationDateFilter(artists, CreationDateFrom, CreationDateTill)
	}

	if firstAlbumFrom != "" || firstAlbumTill != "" {
		if firstAlbumFrom == "" {
			firstAlbumFrom = "1900-01-01"
		}
		if firstAlbumTill == "" {
			firstAlbumTill = "2020-03-03"
		}
		artists = AlbumDateFilter(artists, firstAlbumFrom, firstAlbumTill)
	}

	if sum != 0 {
		artists = MembersFilter(artists, members)
	}

	Maketmpl(w, "index", artists)
}

func IndividualHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	err, limit := GetFullData()
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
		return
	}
	if r.URL.Path != "/individual" && Atoi(id) < 1 || Atoi(id) > limit || id == "" {
		http.Error(w, "ERROR 404: NOT FOUND", http.StatusNotFound)
		return
	}
	if err != nil {
		fmt.Printf("UNABLE TO CONTACT SERVER %s", err.Error())
	}
	artist, error := GetArtistbyID(Atoi(id))
	if error != nil {
		fmt.Println("UNABLE TO RETRIEVE DATAS")
	}
	Maketmpl(w, "individual", artist)
}

func Maketmpl(w http.ResponseWriter, tmplName string, data interface{}) {
	templateCache, err := createTemplateCache()
	if err != nil {
		panic(err)
	}
	tpl, err2d2 := templateCache[tmplName+".html"]
	if !err2d2 {
		http.Error(w, "The template doesn't exist !", http.StatusInternalServerError)
		return
	}
	buff := new(bytes.Buffer)
	tpl.Execute(buff, data)
	buff.WriteTo(w)
}

func createTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.html")
	if err != nil {
		return cache, nil
	}

	for _, page := range pages {
		name := filepath.Base(page)
		tmpl := template.Must(template.ParseFiles(page))
		if err != nil {
			return cache, nil
		}
		cache[name] = tmpl
	}
	return cache, nil
}

func Atoi(s string) int {
	numbs := 0
	c := 0
	esq := []rune(s)
	for _, char := range esq {
		if char >= '0' && char <= '9' {
			for i := '0'; i < char; i++ {
				c++
			}
			numbs = numbs*10 + c
			c = 0
		} else {
			return 0
		}
	}
	return numbs
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
}

type SInput struct {
	Input string
	Id    int
}

func GetArtists() error {
	r, err := http.Get(url + "/artists")
	if err != nil {
		fmt.Println("CANNOT RETRIEVE DATAS", err.Error())
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&artists)
	return nil
}
func GetLocations() error {
	r, err := http.Get(url + "/locations")
	if err != nil {
		fmt.Println("CANNOT RETRIEVE DATAS", err.Error())
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&locations)
	return nil
}
func GetDates() error {
	r, err := http.Get(url + "/dates")
	if err != nil {
		fmt.Println("CANNOT RETRIEVE DATAS", err.Error())
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&dates)
	return nil
}

func GetRelations() error {
	r, err := http.Get(url + "/relation")
	if err != nil {
		fmt.Println("CANNOT RETRIEVE DATAS", err.Error())
	}
	defer r.Body.Close()
	json.NewDecoder(r.Body).Decode(&relations)
	return nil
}

func GetArtistbyID(id int) (Api, error) {
	for _, elem := range artistsall {
		if elem.ID == id {
			return elem, nil
		}
	}
	return Api{}, errors.New("UNABLE TO RETRIEVE DATA FOR THIS ARTIST")
}

func GetArtistbyId(id int) (Artist, error) {
	for _, elem := range artists {
		if elem.ID == id {
			return elem, nil
		}
	}
	return Artist{}, errors.New("UNABLE TO RETRIEVE DATA FOR THIS ARTIST")
}

func Search(w http.ResponseWriter, r *http.Request) {

	searchInput := r.FormValue("q")
	id, sth, _, artone, limit := FindSearch(searchInput, artists)

	fmt.Println(len(sth))
	fmt.Println("ARTIST:", artone)

	Maketmpl(w, "Search", sth)

	fmt.Println("id: ", id)
	fmt.Println("limit:", limit)

}

func FindSearch(searchInput string, artists []Artist) (id int, assc []Artist, ass []Api, artst Artist, limit int) {
	err1 := GetArtists() // check that "url + /artist" does not exist, same for the 3 other functions.
	err2 := GetDates()
	err3 := GetLocations()
	err4 := GetRelations()
	//limit := 0
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		log.Fatal("ERROR TO GET DATAS, ARTISTS, LOCATIONS, DATES")
	}
	//assc := make([]Artist, 52)

	var country []Artist
	var artistTemplate Api
	for i := range artists {
		ok := strconv.Itoa(artists[i].CreationDate)

		limit = 0
		switch {
		case strings.EqualFold(searchInput, ok):
			artst = artists[i]
			artistTemplate.ID = artists[i].ID
			id = artists[i].ID
			country = append(country, artst)
		case GetOneChar(artists[i].Name, strings.Title(searchInput)) || strings.EqualFold(searchInput, artists[i].Name):
			artst = artists[i]
			artistTemplate.Name = searchInput
			artistTemplate.ID = artists[i].ID
			id = artists[i].ID
			country = append(country, artst)
		case strings.EqualFold(searchInput, artists[i].FirstAlbum):
			artistTemplate.FirstAlbum = searchInput
			artst = artists[i]
			country = append(country, artst)
			artistTemplate.ID = artists[i].ID
			id = artists[i].ID
		case GetOneMemChar(strings.Title(searchInput), artists[i].Members) || GetRange(strings.Title(searchInput), artists[i].Members):
			artistTemplate.Members = artists[i].Members
			artst = artists[i]
			country = append(country, artst)
			artistTemplate.ID = artists[i].ID
		default:
			for key, _ := range relations.Index[i].Relation {

				if strings.EqualFold(searchInput, key) {
					artistTemplate.Relation = relations.Index[i].Relation
					artst = artists[i]
					country = append(country, artst)
					artistTemplate.ID = artists[i].ID
					id = relations.Index[i].ID
				}
				indKey := strings.Split(key, "-")
				for _, k := range indKey {
					if strings.Contains(strings.Title(k), strings.Title(searchInput)) {
						artistTemplate.Relation = relations.Index[i].Relation
						artst = artists[i]
						country = append(country, artst)
						artistTemplate.ID = artists[i].ID
						id = relations.Index[i].ID
					}
				}
			}
		}
		limit = i + 1

	}
	ass = append(artistsall, artistTemplate)
	assc = country

	return id, assc, ass, artst, limit
}

/*func DisplayRange(s []string) {

}*/

/*func ShowRelation(ok string, mp map[string][]string) bool {
	for key, _ := range mp {
		if strings.EqualFold(ok, key) {
			return true
		}
	}
	return false
}*/

func GetRange(st string, s []string) bool {
	for _, i := range s {
		if i == st {
			return true
		}
	}
	return false
}

func GetOneChar(s, y string) bool {
	ok := strings.Split(s, " ")
	//var nd string
	for _, i := range ok {
		if strings.Contains(i, y) {
			return true
		}
	}
	return false
}

func GetOneMemChar(st string, s []string) bool {
	for _, i := range s {

		if strings.Contains(i, st) {
			return true
		}
	}
	return false
}

func AutoComplete(cool Artist, s string) bool {
	c := 0
	var x int
	for i := 0; i < len(s); i++ {
		c++
		x = strings.Compare(cool.Name, string(s[i]))
		x = strings.Compare(cool.Members[c], string(s[i]))
		x = strings.Compare(cool.Relations, string(s[i]))
		if x == 0 || x == -1 {
			return true
		}

	}
	return false
}

//----------------------------------------------------------

func MembersFilter(att []Artist, members []int) []Artist {
	var artistSlice []Artist
	for _, mem := range att {
		for _, num := range members {
			if len(mem.Members) == num {
				artistSlice = append(artistSlice, mem)
			}

		}
	}
	return artistSlice
}

func CreationDateFilter(att []Artist, from string, till string) []Artist {
	if from == "" || till == "" {
		return att
	}
	fromDate, err1 := strconv.Atoi(from)
	tillDate, err2 := strconv.Atoi(till)
	if err1 != nil || err2 != nil {
		errors.New("Error by filter by creation date data")
	}

	var artistSlice []Artist

	for _, artist := range att {
		if artist.CreationDate >= fromDate && artist.CreationDate <= tillDate {
			artistSlice = append(artistSlice, artist)
		}
	}
	return artistSlice
}

func AlbumDateFilter(att []Artist, fromDate string, tillDate string) []Artist {
	layOut := "2006-01-02"

	From, err1 := time.Parse(layOut, fromDate)
	Till, err2 := time.Parse(layOut, tillDate)
	if err1 != nil || err2 != nil {
		errors.New("Error by First Album date convert for filter")
	}
	var artistSlice []Artist

	layOutData := "02-01-2006"
	for _, band := range att {
		date, err := time.Parse(layOutData, band.FirstAlbum)
		if err != nil {
			errors.New("Error by First Album date convert for filter")
		}
		if From.Before(date) && Till.After(date) {
			artistSlice = append(artistSlice, band)
		}
	}
	return artistSlice
}

//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------
//-----------------------------------------------------------------------------

/*func FilterByLocation(data []Artist, locations string) []Artist {
	var tmp []Artist

	for _, band := range data {
		for _, loc := range band.Location {
			if locations == string(loc) {
				tmp = append(tmp, band)
			}
		}
	}
	return tmp
}*/

//const googleApiUri = "https://maps.googleapis.com/maps/api/js?key=AIzaSyBFC7OK5o8blljL25668lSZmGbdt2x78Qo&callback=initMap"

/*func getCityCoordinates() error {
	err := GetLocations()
	if err != nil {
		log.Fatal(err)
	}
	var city string
	var geoApiUrl string
	for _, v := range locations.Index {
		for i := range v.Locations {
			sp := strings.Split(v.Locations[i], "-")
			city = sp[len(sp)-1]
			geoApiUrl = fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=AIzaSyBFC7OK5o8blljL25668lSZmGbdt2x78Qo", city)
		}
	}
	fmt.Println("Fetching latitude and longitude of the city ...")

	resp, err := http.Get(geoApiUrl)

	if err != nil {
		log.Fatal("Fetching google api uri data error: ", err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Reading google api data error: ", err)
	}

	var data googleApiResponse
	//var cood Coordinates
	json.Unmarshal(bytes, &data)
	//cood.Latitude = data.Results[0].Geometry.Location.Latitude
	//cood.Longitude = data.Results[0].Geometry.Location.Longitude

	return nil
}*/
