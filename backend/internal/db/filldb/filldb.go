package filldb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/db/models"
	"github.com/pythonistD/Guess-The-Flag/internal/db/repo"
)

func FillCountriesInDb(db *sqlx.DB) error {
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
		if err != nil && errors.Is(err, svgAlreadyExists) {
			fmt.Printf("svg with the same hash for country: %s already exists in DB", country.Name)
			continue
		}
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
	res, err := repo.GetByHash(ctx, flagImgModel.ImageHash)
	if res != nil {
		return 0, svgAlreadyExists
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

// Завести мапу и проверять, встречался ли hash уже
// скипать, если встречался
// проверять, что svg имеет тэги svg, и убирать height и weight

var (
	svgAlreadyExists = errors.New("svg already exists: equal hash found")
)

func createCountryImgModel(imgUrl string) (models.FlagImage, error) {
	imgSvg, err := downloadData(imgUrl)
	if err != nil {
		log.Printf("Couldn't download svg: %v", err)
		return models.FlagImage{}, err
	}
	imgSvg, err = normalizeSVG(imgSvg)
	if err != nil {
		return models.FlagImage{}, fmt.Errorf("failed to normalize svg: %w", err)
	}
	h := hashSVG(imgSvg)
	return models.FlagImage{
		SvgData:   imgSvg,
		ImageHash: h,
		FileSize:  len(imgSvg),
	}, nil

}

func normalizeSVG(svg string) (string, error) {
	svg = strings.TrimSpace(svg)
	openTag := regexp.MustCompile(`^<svg`)
	closeTag := regexp.MustCompile(`</svg>$`)
	if !openTag.MatchString(svg) || !closeTag.MatchString(svg) {
		return "", fmt.Errorf("incorrect svg: doesnt contains svg tags")
	}

	// Работаем только с открывающим тегом <svg ...>: вытаскиваем атрибуты,
	// чтобы не зацепить случайные width/height у вложенных элементов.
	openSvgRe := regexp.MustCompile(`(?s)^<svg\b([^>]*)>`)
	openMatch := openSvgRe.FindStringSubmatchIndex(svg)
	if openMatch == nil {
		return "", fmt.Errorf("incorrect svg: cannot find <svg> open tag")
	}
	attrs := svg[openMatch[2]:openMatch[3]]

	widthRe := regexp.MustCompile(`\s+width="([^"]*)"`)
	heightRe := regexp.MustCompile(`\s+height="([^"]*)"`)
	viewBoxRe := regexp.MustCompile(`\s+viewBox="([^"]*)"`)

	widthMatch := widthRe.FindStringSubmatch(attrs)
	heightMatch := heightRe.FindStringSubmatch(attrs)
	hasViewBox := viewBoxRe.MatchString(attrs)

	newAttrs := widthRe.ReplaceAllString(attrs, "")
	newAttrs = heightRe.ReplaceAllString(newAttrs, "")

	// Если viewBox отсутствует, но были width и height — синтезируем viewBox,
	// чтобы у SVG сохранилось intrinsic aspect ratio для отрисовки в <img>.
	if !hasViewBox && widthMatch != nil && heightMatch != nil {
		w := stripUnits(widthMatch[1])
		h := stripUnits(heightMatch[1])
		if w != "" && h != "" {
			newAttrs += fmt.Sprintf(` viewBox="0 0 %s %s"`, w, h)
		}
	}

	newOpenTag := "<svg" + newAttrs + ">"
	normalizedSvg := newOpenTag + svg[openMatch[1]:]
	return normalizedSvg, nil
}

// stripUnits убирает CSS-единицы из значений атрибутов width/height,
// оставляя голое число (px/pt/% и т.п. в viewBox недопустимы).
func stripUnits(v string) string {
	v = strings.TrimSpace(v)
	unitsRe := regexp.MustCompile(`(?i)(px|pt|pc|mm|cm|in|em|ex|rem|%)$`)
	return strings.TrimSpace(unitsRe.ReplaceAllString(v, ""))
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
