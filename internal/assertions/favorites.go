package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites/models"
)

func ContainsFavoriteByAdvertisementID(t *testing.T, items []models.Favorite, advertisementID string) {
	t.Helper()
	for _, item := range items {
		if item.Advertisement.ID == advertisementID {
			return
		}
	}
	require.Failf(t, "избранное не содержит объявление", "В избранном нет объявления с ID %q", advertisementID)
}

func DoesNotContainFavoriteByAdvertisementID(t *testing.T, items []models.Favorite, advertisementID string) {
	t.Helper()
	for _, item := range items {
		if item.Advertisement.ID == advertisementID {
			require.Failf(t, "объявление осталось в избранном", "В избранном найдено объявление с ID %q", advertisementID)
			return
		}
	}
}
