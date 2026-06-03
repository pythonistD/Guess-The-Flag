package filldb

type Country struct {
	Name         CountryName            `json:"name"`
	Flag         Flag                   `json:"flags"`
	Translations map[string]Translation `json:"translations"`
	Code         string                 `json:"cca2"`
}

type CountryName struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

type Flag struct {
	SVG string `json:"svg"`
	PNG string `json:"png"`
}

type Translation struct {
	Official string `json:"official"`
	Common   string `json:"common"`
}
