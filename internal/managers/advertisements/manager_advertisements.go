package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	httpAdvertisements "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

func CreateAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest, photoPaths []string) models.AdvertisementResult {
	return parseAdvertisementResponse(t, httpAdvertisements.HttpPostAdvertisement(t, token, request, photoPaths))
}

func GetAdvertisement(t *testing.T, token string, id string) models.AdvertisementResult {
	return parseAdvertisementResponse(t, httpAdvertisements.HttpGetAdvertisement(t, token, id))
}

func GetAdvertisementPhotos(t *testing.T, token string, id string) models.PhotosResult {
	response := httpAdvertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ с фотографиями")

	result := models.PhotosResult{StatusCode: response.StatusCode}
	if response.StatusCode == http.StatusOK {
		require.NoError(t, json.Unmarshal(body, &result.Photos), "не удалось разобрать список фотографий")
	} else if len(body) > 0 {
		require.NoError(t, json.Unmarshal(body, &result.Error), "не удалось разобрать ошибку фотографий")
	}
	return result
}

func DeleteAdvertisement(t *testing.T, token string, id string) models.AdvertisementResult {
	return parseAdvertisementResponse(t, httpAdvertisements.HttpDeleteAdvertisement(t, token, id))
}

func GetAdvertisementsList(t *testing.T, token string, textForSearch string) models.AdvertisementsList {
	response := httpAdvertisements.HttpGetAdvertisements(t, token, textForSearch)
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать список объявлений")

	var advertisementsList models.AdvertisementsList
	err = json.Unmarshal(body, &advertisementsList)
	require.NoError(t, err, "не удалось разобрать список объявлений")

	return advertisementsList
}

func parseAdvertisementResponse(t *testing.T, response *http.Response) models.AdvertisementResult {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "не удалось прочитать ответ объявления")

	result := models.AdvertisementResult{StatusCode: response.StatusCode}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusCreated {
		require.NoError(t, json.Unmarshal(body, &result.Advertisement), "не удалось разобрать объявление")
	} else if len(body) > 0 {
		require.NoError(t, json.Unmarshal(body, &result.Error), "не удалось разобрать ошибку объявления")
	}
	return result
}
