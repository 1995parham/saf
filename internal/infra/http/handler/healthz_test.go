package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham/saf/internal/http/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type HealthzSuite struct {
	suite.Suite

	engine *fiber.App
}

func (suite *HealthzSuite) SetupSuite() {
	suite.engine = fiber.New()

	handler.Healthz{
		Logger: zap.NewNop(),
		Tracer: trace.NewNoopTracerProvider().Tracer(""),
	}.Register(suite.engine.Group(""))
}

func (suite *HealthzSuite) TestHandler() {
	require := suite.Require()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := suite.engine.Test(req)
	require.NoError(err)
	require.Equal(http.StatusNoContent, resp.StatusCode)
	require.NoError(resp.Body.Close())
}

func TestHealthzSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HealthzSuite))
}
