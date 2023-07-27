package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"groupie-tracker-filters/models"
)

func ParseJson() (models.Artists, error) {
	musicLinks := []string{
		"https://open.spotify.com/embed/track/1lCRw5FEZ1gPDNPzy1K4zW?utm_source=generator",
		"https://open.spotify.com/embed/track/1HnriuDThLq7bEl1QKiaJL?utm_source=generator", "https://open.spotify.com/embed/track/7LlOGWlocEpRMccB6Gaehf?utm_source=generator",
		"https://open.spotify.com/embed/track/3ovjw5HZZv43SxTwApooCM?utm_source=generator", "https://open.spotify.com/embed/track/3ee8Jmje8o58CHK66QrVC2?utm_source=generator",
		"https://open.spotify.com/embed/track/4gT3mNJA8lnlkYFqGZ8IA2?utm_source=generator", "https://open.spotify.com/embed/track/69fnnbeoZO9h16CUZGIbTG?utm_source=generator",
		"https://open.spotify.com/embed/track/7KXjTSCq5nL1LoYtL7XAwS?utm_source=generator", "https://open.spotify.com/embed/track/57bgtoPSgt236HzfBOd8kj?utm_source=generator",
		"https://open.spotify.com/embed/track/1L94M3KIu7QluZe63g64rv?utm_source=generator", "https://open.spotify.com/embed/track/4jbmgIyjGoXjY01XxatOx6?utm_source=generator",
		"https://open.spotify.com/embed/track/6O20JhBJPePEkBdrB5sqRx?utm_source=generator", "https://open.spotify.com/embed/track/7safX55XidhznxK5eDdDm5?utm_source=generator",
		"https://open.spotify.com/embed/track/4YwbSZaYeYja8Umyt222Qf?utm_source=generator", "https://open.spotify.com/embed/track/5CQ30WqJwcep0pYcV4AMNc?utm_source=generator",
		"https://open.spotify.com/embed/track/1Eolhana7nKHYpcYpdVcT5?utm_source=generator", "https://open.spotify.com/embed/track/2JoZzpdeP2G6Csfdq5aLXP?utm_source=generator",
		"https://open.spotify.com/embed/track/2JhJOPGvtqMpj5RQC8cIYf?utm_source=generator", "https://open.spotify.com/embed/track/4JfuiOWlWCkjP6OKurHjSn?utm_source=generator",
		"https://open.spotify.com/embed/track/4nFNJmjfgBF7jwv2oBC45b?utm_source=generator", "https://open.spotify.com/embed/track/6cuLsMmHlR08BdoWAIwYVh?utm_source=generator",
		"https://open.spotify.com/embed/track/1PYG9Akj0LAZZUDXzV9m1S?utm_source=generator", "https://open.spotify.com/embed/track/0pqnGHJpmpxLKifKRmU6WP?utm_source=generator",
		"https://open.spotify.com/embed/track/4VXIryQMWpIdGgYR4TrjT1?utm_source=generator", "https://open.spotify.com/embed/track/5E4mQ2mXblbeuI4tefnMZG?utm_source=generator",
		"https://open.spotify.com/embed/track/2qxmye6gAegTMjLKEBoR3d?utm_source=generator", "https://open.spotify.com/embed/track/4hObp5bmIJ3PP3cKA9K9GY?utm_source=generator",
		"https://open.spotify.com/embed/track/7aJgh6LCvhXJfD7PHjhG70?utm_source=generator", "https://open.spotify.com/embed/track/5w40ZYhbBMAlHYNDaVJIUu?utm_source=generator",
		"https://open.spotify.com/embed/track/40iJIUlhi6renaREYGeIDS?utm_source=generator", "https://open.spotify.com/embed/track/2JvzF1RMd7lE3KmFlsyZD8?utm_source=generator",
		"https://open.spotify.com/embed/track/3RlsVPIIs5KFhLFhxZ4iDF?utm_source=generator", "https://open.spotify.com/embed/track/7N1Vjtzr1lmmCW9iasQ8YO?utm_source=generator",
		"https://open.spotify.com/embed/track/0vfPEfQk0ZCHExTZ007Ryr?utm_source=generator", "https://open.spotify.com/embed/track/5n8Aro6j1bEGIy7Tpo7FV7?utm_source=generator",
		"https://open.spotify.com/embed/track/6ADSaE87h8Y3lccZlBJdXH?utm_source=generator", "https://open.spotify.com/embed/track/5XeFesFbtLpXzIVDNQP22n?utm_source=generator",
		"https://open.spotify.com/embed/track/3Zwu2K0Qa5sT6teCCHPShP?utm_source=generator", "https://open.spotify.com/embed/track/0d28khcov6AiegSCpG5TuT?utm_source=generator",
		"https://open.spotify.com/embed/track/4yugZvBYaoREkJKtbG08Qr?utm_source=generator", "https://open.spotify.com/embed/track/18lR4BzEs7e3qzc0KVkTpU?utm_source=generator",
		"https://open.spotify.com/embed/track/3ZOEytgrvLwQaqXreDs2Jx?utm_source=generator", "https://open.spotify.com/embed/track/4woTEX1wYOTGDqNXuavlRC?utm_source=generator",
		"https://open.spotify.com/embed/track/01VnqAxzbuKVunmItkraw5?utm_source=generator", "https://open.spotify.com/embed/track/0nLiqZ6A27jJri2VCalIUs?utm_source=generator",
		"https://open.spotify.com/embed/track/3AJwUDP919kvQ9QcozQPxg?utm_source=generator", "https://open.spotify.com/embed/track/4cktbXiXOapiLBMprHFErI?utm_source=generator",
		"https://open.spotify.com/embed/track/0yzj6eBs5a6X6d3P5qaQ5J?utm_source=generator", "https://open.spotify.com/embed/track/63T7DJ1AFDD6Bn8VzG6JE8?utm_source=generator",
		"https://open.spotify.com/embed/track/3lPr8ghNDBLc2uZovNyLs9?utm_source=generator", "https://open.spotify.com/embed/track/5OQsiBsky2k2kDKy2bX2eT?utm_source=generator",
		"https://open.spotify.com/embed/track/7BKLCZ1jbUBVqRi2FVlTVw?utm_source=generator",
	}
	r, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Println("Cannot get from URL", err)
		return nil, err
	}
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading json data:", err)
		return nil, err
	}

	var relations models.Relations

	err = json.Unmarshal(data, &relations)
	if err != nil {
		log.Println("Error unmarshalling relation json data:", err)
		return nil, err
	}

	r, err = http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println("Cannot get from URL", err)
		return nil, err
	}
	defer r.Body.Close()

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading json data:", err)
		return nil, err
	}

	var artists models.Artists
	err = json.Unmarshal(data, &artists)
	if err != nil {
		log.Println("Error unmarshalling json data:", err)
		return nil, err
	}
	for i := 0; i < len(artists); i++ {
		artists[i].Relation = relations.Index[i]
		artists[i].MusicLink = musicLinks[i]
	}
	return artists, nil
}
