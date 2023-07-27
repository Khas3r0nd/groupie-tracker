package internal

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"groupie-tracker-filters/models"
)

type response struct {
	Groups                models.Artists
	MinAndMaxCreationDate []int
	MinAndMaxAlbumDate    []string
	MaxArtists            int
	CountriesMap          map[string][]string
	AlbumRanges           []string
	CreationRanges        []int
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		ErrorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	funcMap := template.FuncMap{
		"loop": func(from, to int) <-chan int {
			ch := make(chan int)
			go func() {
				for i := from; i <= to; i++ {
					ch <- i
				}
				close(ch)
			}()
			return ch
		},
	}
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Println("Error parsing template")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	groups, err := ParseJson()
	if err != nil {
		log.Println("Error during parse process")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	countryMap := make(map[string][]string)
	for _, group := range groups {
		for key := range group.Relation.Relations {
			parts := strings.Split(key, "-")
			city := strings.Title(strings.ReplaceAll(parts[0], "_", " "))

			country := strings.Title(strings.ReplaceAll(parts[1], "_", " "))

			if !containsCountry(countryMap[country], city) {
				countryMap[country] = append(countryMap[country], city)
			}
		}
	}

	rangeCreationDate, rangeAlbumDate, maxArtists := findMinAndMax(groups)

	res := response{
		Groups:                groups,
		MinAndMaxCreationDate: rangeCreationDate[:],
		MinAndMaxAlbumDate:    rangeAlbumDate[:],
		MaxArtists:            maxArtists,
		CountriesMap:          countryMap,
		AlbumRanges:           []string{rangeAlbumDate[0], rangeAlbumDate[1]},
		CreationRanges:        []int{rangeCreationDate[0], rangeCreationDate[1]},
	}

	err = t.Execute(w, res)
	if err != nil {
		log.Println("Error executing template:", err)
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	err := r.ParseForm()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if !r.Form.Has("text") {
		ErrorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	funcMap := template.FuncMap{
		"loop": func(from, to int) <-chan int {
			ch := make(chan int)
			go func() {
				for i := from; i <= to; i++ {
					ch <- i
				}
				close(ch)
			}()
			return ch
		},
	}
	lowRangeAlbum, highRangeAlbum := r.FormValue("low-rangeAlbum"), r.FormValue("high-rangeAlbum")
	lowRangeCreation, _ := strconv.Atoi(r.FormValue("low-rangeCreation"))
	highRangeCreation, _ := strconv.Atoi(r.FormValue("high-rangeCreation"))
	isAlbumFilterNeeded, isCreationFilterNeeded := r.FormValue("AlbumDateFilter") == "on", r.FormValue("CreationDateFilter") == "on"
	text := r.FormValue("text")
	groups, err := ParseJson()
	if err != nil {
		log.Println("Error during parse process line 89")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	cityCheckBoxes := r.Form["cities"]
	checkBoxes := r.Form["MemberCountFilter"]
	checkedValues := []int{}

	for _, value := range checkBoxes {
		checkedValue, err := strconv.Atoi(value)
		if err == nil {
			checkedValues = append(checkedValues, checkedValue)
		}
	}

	result := FindMatches(strings.TrimSpace(text), groups, lowRangeAlbum, highRangeAlbum, isAlbumFilterNeeded, lowRangeCreation, highRangeCreation, isCreationFilterNeeded, checkedValues, cityCheckBoxes)
	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		log.Println(err)
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	countryMap := make(map[string][]string)
	for _, group := range groups {
		for key := range group.Relation.Relations {
			parts := strings.Split(key, "-")
			city := strings.Title(strings.ReplaceAll(parts[0], "_", " "))
			country := strings.Title(strings.ReplaceAll(parts[1], "_", " "))

			if !containsCountry(countryMap[country], city) {
				countryMap[country] = append(countryMap[country], city)
			}
		}
	}
	rangeCreationDate, rangeAlbumDate, maxArtists := findMinAndMax(groups)

	res := response{
		result,
		rangeCreationDate[:],
		rangeAlbumDate[:],
		maxArtists,
		countryMap,
		[]string{lowRangeAlbum, highRangeAlbum},
		[]int{lowRangeCreation, highRangeCreation},
	}

	err = t.Execute(w, res)
	if err != nil {
		log.Println("Error executing template:", err)
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func FindMatches(text string, allData models.Artists, lowRangeAlbum, highRangeAlbum string, isAlbumFilterNeeded bool, lowRangeCreation, highRangeCreation int, isCreationFilterNeeded bool, checkedValues []int, cityCheckBoxes []string) models.Artists {
	var sortedArtists models.Artists

	text = strings.ToLower(text)
	for _, artist := range allData {
		if containsText(artist.Name, text) || containsText(artist.FirstAlbum, text) ||
			containsText(strconv.Itoa(artist.CreationDate), text) || containsMember(artist.Members, text) ||
			containsRelationDate(artist.Relation.Relations, text) || containsRelationPlace(artist.Relation.Relations, text) {
			if isAlbumFilterNeeded && (artist.FirstAlbum[6:] < lowRangeAlbum || artist.FirstAlbum[6:] > highRangeAlbum) {
				continue
			}
			if isCreationFilterNeeded && (artist.CreationDate < lowRangeCreation || artist.CreationDate > highRangeCreation) {
				continue
			}
			if !containsNeededMember(len(artist.Members), checkedValues) && len(checkedValues) != 0 {
				continue
			}
			if !containsNeededCity(artist.Relation.Relations, cityCheckBoxes) && len(cityCheckBoxes) != 0 {
				continue
			}
			if !containsID(artist.Id, sortedArtists) {
				sortedArtists = append(sortedArtists, artist)
			}
		}
	}

	return sortedArtists
}

func containsNeededCity(artistsCountries map[string][]string, neededCity []string) bool {
	for key := range artistsCountries {
		parts := strings.Split(key, "-")
		city := strings.Title(strings.ReplaceAll(parts[0], "_", " "))
		for _, value := range neededCity {
			if city == value {
				return true
			}
		}

	}

	return false
}

func containsNeededMember(memberCount int, neededMembers []int) bool {
	for _, value := range neededMembers {
		if memberCount == value {
			return true
		}
	}
	return false
}

func containsText(value, text string) bool {
	return strings.Contains(strings.ToLower(value), text)
}

func containsMember(members []string, text string) bool {
	for _, member := range members {
		if containsText(member, text) {
			return true
		}
	}
	return false
}

func containsRelationDate(relations map[string][]string, text string) bool {
	if len(text) < 2 {
		return false
	}
	text = text[1 : len(text)-1]
	for _, value := range relations {
		for _, date := range value {
			if containsText(date, text) {
				return true
			}
		}
	}
	return false
}

func containsRelationPlace(relations map[string][]string, text string) bool {
	for key := range relations {
		if containsText(key, text) {
			return true
		}
	}
	return false
}

func containsID(id int, data models.Artists) bool {
	for _, artist := range data {
		if artist.Id == id {
			return true
		}
	}
	return false
}

func containsCountry(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func Concert(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	if r.URL.Path == "/artists/" && len(r.URL.Query()) == 0 {
		http.Redirect(w, r, "http://localhost:8080", http.StatusSeeOther)
		return
	}
	if !r.URL.Query().Has("id") {
		ErrorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		ErrorHandler(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	pJ, err := ParseJson()
	if err != nil {
		log.Println("Error during parse proccess")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	if err != nil || id < 1 || id > len(pJ) {
		ErrorHandler(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}

	switch r.Method {
	case http.MethodGet:
		ConcertGet(w, r, id-1)
	default:
		ErrorHandler(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func ConcertGet(w http.ResponseWriter, r *http.Request, artistId int) {
	t, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		log.Println("Error parsing template at")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	artists, err := ParseJson()
	if err != nil {
		log.Println("Error during parse proccess")
		ErrorHandler(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	t.Execute(w, artists[artistId])
}

func findMinAndMax(groups models.Artists) ([]int, []string, int) {
	minAndMaxCreationDate := []int{groups[0].CreationDate, groups[0].CreationDate}
	minAndMaxAlbumDate := []string{groups[0].FirstAlbum[6:], groups[0].FirstAlbum[6:]}
	maxArtists := len(groups[0].Members)
	for _, artist := range groups {
		if len(artist.Members) > maxArtists {
			maxArtists = len(artist.Members)
		}
		if minAndMaxCreationDate[0] > artist.CreationDate {
			minAndMaxCreationDate[0] = artist.CreationDate
		}
		if minAndMaxCreationDate[1] < artist.CreationDate {
			minAndMaxCreationDate[1] = artist.CreationDate
		}
		if minAndMaxAlbumDate[0] > artist.FirstAlbum[6:] {
			minAndMaxAlbumDate[0] = artist.FirstAlbum[6:]
		}
		if minAndMaxAlbumDate[1] < artist.FirstAlbum[6:] {
			minAndMaxAlbumDate[1] = artist.FirstAlbum[6:]
		}
	}
	return minAndMaxCreationDate, minAndMaxAlbumDate, maxArtists
}
