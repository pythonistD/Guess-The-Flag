package main

type Country struct {
	Name string `json:"name"`
	Code string `json:"numericCode"`
	Flag Flag   `json:"flags"`
}

type Flag struct {
	SVG string `json:"svg"`
	PNG string `json:"png"`
}
