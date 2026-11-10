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
