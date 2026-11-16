package favorites

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath).
		Query("page", strconv.Itoa(1)).
		Query("page_size", strconv.Itoa(50)).
		Expect(t).
		End().Response
}

func HttpPostFavorite(t *testing.T, token string, advertisementID string) *http.Response {
	request := models.AddFavoriteRequest{AdvertisementID: advertisementID}
	return apiRunner.GetRunner().Auth(token).Create().Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

func HttpDeleteFavorite(t *testing.T, token string, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(fmt.Sprintf("%s/%s", path.FavoritesPath, url.PathEscape(advertisementID))).
		Expect(t).
		End().Response
}
