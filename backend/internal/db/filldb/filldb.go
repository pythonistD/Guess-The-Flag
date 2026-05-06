package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
	"github.com/pythonistD/Guess-The-Flag/internal/db"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
)

func fillCountriesInDb(db *sqlx.DB) error {
	countries, err := getCountriesFromApi("https://restcountries.com/v3.1/independent?fields=name,translations,flags,cca2")
	if err != nil {
		return fmt.Errorf("error fetching countries from API: %w", err)
	}
	err = createCountries(db, countries)
	if err != nil {
		return fmt.Errorf("error creating countries: %w", err)
	}
	fmt.Println("Countries created successfully")
	return nil
}

func getCountriesFromApi(url string) ([]Country, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected API status %s: %s", resp.Status, string(body))
	}

	var countries []Country
	err = json.NewDecoder(resp.Body).Decode(&countries)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode json response: %w", err)
	}
	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Countries parsed: %d\n", len(countries))
	return countries, nil
}

func createCountries(db *sqlx.DB, countriesSchema []Country) error {
	countryRepo := repo.NewCountriesRepo(db)
	imagesRepo := repo.NewImagesRepo(db)
	countryNamesRepo := repo.NewCountryNamesRepo(db)
	ctx := context.Background()

	for _, country := range countriesSchema {
		imgId, err := addCountryImgToDb(ctx, imagesRepo, country.Flag.SVG)
		if err != nil {
			return fmt.Errorf("error creating country %s Code: %s img: %w", country.Name, country.Code, err)
		}
		countryModel := models.Country{
			Code:        country.Code,
			FlagImageId: imgId,
		}
		countryDb, err := countryRepo.Create(ctx, &countryModel)
		if err != nil {
			return fmt.Errorf("error adding country to db: %w", err)
		}
		err = createCountryNames(ctx, countryNamesRepo, countryDb.CountryId, country.Translations)
		if err != nil {
			return fmt.Errorf("error adding country names to db: %w", err)
		}
	}
	return nil
}

func addCountryImgToDb(ctx context.Context, repo *repo.ImagesRepo, url string) (int, error) {
	flagImgModel, err := createCountryImgModel(url)
	if err != nil {
		return 0, fmt.Errorf("error creating country img model: %w", err)
	}
	flagImgDb, err := repo.Create(ctx, &flagImgModel)
	if err != nil {
		return 0, fmt.Errorf("error adding flag img to db, img hash:%s %w", flagImgModel.ImageHash, err)
	}
	return flagImgDb.ImageId, nil
}

func createCountryNames(ctx context.Context, repo *repo.CountryNamesRepo, countryId int, translations map[string]Translation) error {
	for langCode, t := range translations {
		countryNamesModelOfficial := models.CountryNames{
			LanguageCode:   langCode,
			CountryId:      countryId,
			Name:           t.Official,
			NormalizedName: normalizeName(t.Official),
			Threshold:      calculateThreshold(t.Official),
			IsDisplayName:  true,
		}
		countryNamesModelCommon := models.CountryNames{
			LanguageCode:   langCode,
			CountryId:      countryId,
			Name:           t.Common,
			NormalizedName: normalizeName(t.Common),
			Threshold:      calculateThreshold(t.Common),
			IsDisplayName:  false,
		}
		err := repo.CreateAll(ctx, []models.CountryNames{countryNamesModelOfficial, countryNamesModelCommon})
		if err != nil {
			return fmt.Errorf("error adding country names to db: %w", err)
		}
	}
	return nil
}

func normalizeName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, "!", "")
	name = strings.ReplaceAll(name, "?", "")
	return name
}

func calculateThreshold(name string) int {
	l := utf8.RuneCountInString(strings.TrimSpace(name))
	switch {
	case l <= 4:
		return 1
	case l <= 8:
		return 2
	default:
		return int(math.Ceil(float64(l) * 0.25))
	}
}

func createCountryImgModel(imgUrl string) (models.FlagImage, error) {
	imgSvg, err := downloadData(imgUrl)
	if err != nil {
		log.Printf("Couldn't download svg: %v", err)
		return models.FlagImage{}, err
	}
	return models.FlagImage{
		SvgData:   imgSvg,
		ImageHash: hashSVG(imgSvg),
		FileSize:  len(imgSvg),
	}, nil

}

func downloadData(imgUrl string) (string, error) {
	response, err := http.Get(imgUrl)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func hashSVG(svg string) string {
	h := sha256.Sum256([]byte(svg))
	return hex.EncodeToString(h[:])
}

func main() {
	yamlConfigPath := flag.String("config", "./config.yml", "path to config file")
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
	var existing int
	err = database.Get(&existing, "SELECT COUNT(*) FROM countries")
	if err != nil {
		log.Fatal(err)
	}
	if existing > 0 {
		fmt.Printf("countries already filled (%d rows), skipping\n", existing)
		return
	}
	err = fillCountriesInDb(database)
	if err != nil {
		log.Fatal(err)
	}
}
