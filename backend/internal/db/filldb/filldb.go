package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
	"github.com/pythonistD/Guess-The-Flag/internal/db"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
	"log"
	"net/http"
)

func fillCountriesInDb(db *sqlx.DB) error {
	countries, _ := getCountriesFromApi("https://www.apicountries.com/countries")
	fmt.Printf("Countries: %v", countries)
	countriesDb, _ := transformFromJsonToModel(countries)
	countryRepo := repo.NewCountriesRepo(db)
	err := addCountriesToDB(countriesDb, countryRepo)
	if err != nil {
		log.Print(err)
	}
	return nil
}

func getCountriesFromApi(url string) ([]Country, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var countries []Country
	err = json.NewDecoder(resp.Body).Decode(&countries)
	if err != nil {
		fmt.Print("Couldn't decode json response")
		return nil, err
	}
	fmt.Printf("Response status: %s", resp.Status)
	fmt.Printf("Countries parsed: %v", countries)
	return countries, nil
}

func transformFromJsonToModel(countriesSchema []Country) ([]models.Country, error) {
	countriesDb := make([]models.Country, 0, len(countriesSchema))
	for _, c := range countriesSchema {
		countryNew := models.Country{
			Name:    c.Name,
			Code:    c.Code,
			FlagUrl: c.Flag.PNG,
		}
		countriesDb = append(countriesDb, countryNew)
	}
	return countriesDb, nil
}

func addCountriesToDB(countries []models.Country, countriesRepo *repo.CountriesRepo) error {
	ctx := context.Background()
	err := countriesRepo.CreateAll(ctx, countries)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	yamlConfigPath := flag.String("config", "C:\\Users\\User\\GolandProjects\\Guess-The-Flag\\backend\\config.yml", "path to config file")
	flag.Parse()
	cfg, err := config.LoadConfigFromFile(*yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.NewPostgres(cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = fillCountriesInDb(database)
	if err != nil {
		log.Fatal(err)
	}
}
