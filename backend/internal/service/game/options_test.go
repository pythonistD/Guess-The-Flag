package game

import (
	"testing"

	"github.com/pythonistD/Guess-The-Flag/internal/service/game/storage"
)

type mockCountryStorage struct {
	countries map[int]storage.Country
	order     []int
	next      int
}

func (m *mockCountryStorage) InitCountryStorageState() error { return nil }

func (m *mockCountryStorage) GetRandom() (storage.Country, error) {
	if len(m.order) == 0 {
		return storage.Country{}, nil
	}
	id := m.order[m.next%len(m.order)]
	m.next++
	return m.countries[id], nil
}

func (m *mockCountryStorage) GetByID(id int) (storage.Country, error) {
	return m.countries[id], nil
}

func (m *mockCountryStorage) GetAll() []storage.Country {
	result := make([]storage.Country, 0, len(m.countries))
	for _, c := range m.countries {
		result = append(result, c)
	}
	return result
}

func newMockCountryStorage() *mockCountryStorage {
	return &mockCountryStorage{
		order: []int{1, 2, 3, 4, 5},
		countries: map[int]storage.Country{
			1: {
				Id: 1,
				CommonCountryNames: map[string]storage.CountryName{
					"rus": {Name: "Франция", NormalizedName: "франция", Threshold: 2},
				},
			},
			2: {
				Id: 2,
				CommonCountryNames: map[string]storage.CountryName{
					"rus": {Name: "Германия", NormalizedName: "германия", Threshold: 2},
				},
			},
			3: {
				Id: 3,
				CommonCountryNames: map[string]storage.CountryName{
					"rus": {Name: "Италия", NormalizedName: "италия", Threshold: 2},
				},
			},
			4: {
				Id: 4,
				CommonCountryNames: map[string]storage.CountryName{
					"rus": {Name: "Испания", NormalizedName: "испания", Threshold: 2},
				},
			},
			5: {
				Id: 5,
				CommonCountryNames: map[string]storage.CountryName{
					"rus": {Name: "Польша", NormalizedName: "польша", Threshold: 2},
				},
			},
		},
	}
}

func TestGenerateOptions(t *testing.T) {
	store := newMockCountryStorage()
	used := map[int]struct{}{2: {}}

	options, err := GenerateOptions(store, 1, "rus", used)
	if err != nil {
		t.Fatalf("GenerateOptions returned error: %v", err)
	}
	if len(options) != 4 {
		t.Fatalf("expected 4 options, got %d", len(options))
	}

	seen := make(map[int]struct{})
	hasCorrect := false
	for _, opt := range options {
		if _, ok := seen[opt.CountryId]; ok {
			t.Fatalf("duplicate country_id in options: %d", opt.CountryId)
		}
		seen[opt.CountryId] = struct{}{}
		if opt.CountryId == 1 {
			hasCorrect = true
			if opt.Name != "Франция" {
				t.Fatalf("expected correct option name Франция, got %s", opt.Name)
			}
		}
	}
	if !hasCorrect {
		t.Fatal("correct country not found in options")
	}
}
