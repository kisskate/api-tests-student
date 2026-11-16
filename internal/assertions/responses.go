package assertions

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/utils"
)

func APIError(t *testing.T, expectedStatus int, actualStatus int, expectedErrorCode string, actualErrorCode string, expectedMessage string, actualMessage string) {
	t.Helper()
	utils.LogWithLabelAndTimeStamp("API error", "status="+strconv.Itoa(actualStatus)+", error="+actualErrorCode+", message="+actualMessage)
	require.Equalf(t, expectedStatus, actualStatus, "неверный HTTP status code: error=%q, message=%q", actualErrorCode, actualMessage)
	require.Equal(t, expectedErrorCode, actualErrorCode, "неверный код ошибки API")
	require.Equal(t, expectedMessage, actualMessage, "неверное сообщение ошибки API")
}
