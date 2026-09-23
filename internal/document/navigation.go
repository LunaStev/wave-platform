package document

import (
	"sort"
	"strings"
)

// ProjectForPath uses document paths, not groups or titles, as the namespace.
func ProjectForPath(path string) string {
	if strings.HasPrefix(path, "whale/") {
		return "whale"
	}
	return "wave"
}

// Navigation includes untranslated English paths, with translations replacing
// the corresponding English entries before metadata ordering is applied.
func (repository *Repository) Navigation(locale, project string) ([]Summary, error) {
	translated, err := repository.Summaries(locale)
	if err != nil {
		return nil, err
	}
	english := []Summary{}
	if locale != "en" {
		english, err = repository.Summaries("en")
		if err != nil {
			return nil, err
		}
	}
	byPath := make(map[string]Summary)
	for _, items := range [][]Summary{english, translated} {
		for _, item := range items {
			if ProjectForPath(item.Path) == project {
				byPath[item.Path] = item
			}
		}
	}
	result := make([]Summary, 0, len(byPath))
	for _, item := range byPath {
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].GroupOrder != result[j].GroupOrder {
			return result[i].GroupOrder < result[j].GroupOrder
		}
		if result[i].Order != result[j].Order {
			return result[i].Order < result[j].Order
		}
		return result[i].Path < result[j].Path
	})
	return result, nil
}
