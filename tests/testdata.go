package tests

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

func CreateTestPhoto(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "testboard-photo-*.png")
	require.NoError(t, err, "не удалось создать временный файл фотографии")

	photo := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for y := 0; y < 480; y++ {
		for x := 0; x < 640; x++ {
			photo.SetRGBA(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: 128, A: 255})
		}
	}
	require.NoError(t, png.Encode(file, photo), "не удалось записать PNG-фотографию")
	require.NoError(t, file.Close(), "не удалось закрыть файл фотографии")
	return file.Name()
}

func NewTestAdvertisementRequest() models.CreateAdvertisementRequest {
	return models.CreateAdvertisementRequest{
		Title:       "api-test-" + utils.RandomString(12),
		Description: "TestBoard API integration test advertisement",
		Price:       14999,
		Quantity:    2,
	}
}

func CreateTestAdvertisement(t *testing.T, token string) (models.Advertisement, func()) {
	t.Helper()
	Precondition("Создаем временное объявление для тестового сценария")

	request := NewTestAdvertisementRequest()
	result := advertisements.CreateAdvertisement(t, token, request, []string{CreateTestPhoto(t)})
	require.Equal(t, 201, result.StatusCode, "не удалось создать тестовое объявление: %s", result.Error.Message)
	require.NotEmpty(t, result.Advertisement.ID, "API создал объявление без ID")

	deleted := false
	t.Cleanup(func() {
		if deleted {
			return
		}
		Precondition("Удаляем временное объявление после теста")
		deleteResult := advertisements.DeleteAdvertisement(t, token, result.Advertisement.ID)
		if deleteResult.StatusCode != 204 && deleteResult.StatusCode != 404 {
			t.Errorf("не удалось очистить тестовое объявление: status=%d, error=%s", deleteResult.StatusCode, deleteResult.Error.Message)
		}
	})

	return result.Advertisement, func() { deleted = true }
}
