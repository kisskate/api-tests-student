package assertions

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Title), strings.ToLower(title)) {
			return
		}
	}

	require.Failf(t, "объявление не найдено", "Ни одно название не содержит %q", title)
}

func ContainsAdvertisementByID(t *testing.T, items []models.Advertisement, id string) {
	t.Helper()
	for _, item := range items {
		if item.ID == id {
			return
		}
	}
	require.Failf(t, "объявление не найдено", "В списке нет объявления с ID %q", id)
}
