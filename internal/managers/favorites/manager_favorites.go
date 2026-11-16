package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	httpFavorites "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

func GetFavorites(t *testing.T, token string) models.FavoritesListResult {
	response := httpFavorites.HttpGetFavorites(t, token)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать список избранного")

	result := models.FavoritesListResult{StatusCode: response.StatusCode}
	if response.StatusCode == http.StatusOK {
		require.NoError(t, json.Unmarshal(body, &result.List), "не удалось разобрать список избранного")
	} else if len(body) > 0 {
		require.NoError(t, json.Unmarshal(body, &result.Error), "не удалось разобрать ошибку списка избранного")
	}
	return result
}

func AddFavorite(t *testing.T, token string, advertisementID string) models.FavoriteResult {
	return parseFavoriteResponse(t, httpFavorites.HttpPostFavorite(t, token, advertisementID))
}

func DeleteFavorite(t *testing.T, token string, advertisementID string) models.FavoriteResult {
	return parseFavoriteResponse(t, httpFavorites.HttpDeleteFavorite(t, token, advertisementID))
}

func parseFavoriteResponse(t *testing.T, response *http.Response) models.FavoriteResult {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ избранного")

	result := models.FavoriteResult{StatusCode: response.StatusCode}
	if response.StatusCode == http.StatusCreated {
		require.NoError(t, json.Unmarshal(body, &result.Favorite), "не удалось разобрать добавленное избранное")
	} else if len(body) > 0 {
		require.NoError(t, json.Unmarshal(body, &result.Error), "не удалось разобрать ошибку избранного")
	}
	return result
}
