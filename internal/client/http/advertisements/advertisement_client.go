package advertisements

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

func HttpPostAdvertisement(t *testing.T, token string, request models.CreateAdvertisementRequest, photoPaths []string) *http.Response {
	runner := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath).
		MultipartFormData("title", request.Title).
		MultipartFormData("description", request.Description).
		MultipartFormData("price", strconv.FormatFloat(request.Price, 'f', -1, 64)).
		MultipartFormData("quantity", strconv.Itoa(request.Quantity))
	if len(photoPaths) > 0 {
		runner.MultipartFile("photos", photoPaths...)
	}
	return runner.Expect(t).End().Response
}

func HttpGetAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

func HttpGetAdvertisementPhotos(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Get(fmt.Sprintf("%s/%s/photos", path.AdvertisementsPath, id)).
		Expect(t).
		End().Response
}

func HttpDeleteAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}
