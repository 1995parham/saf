package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham/saf/internal/infra/http/handler"
	"github.com/1995parham/saf/internal/infra/logger"
	"github.com/1995parham/saf/internal/infra/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

type HealthzSuite struct {
	suite.Suite

	engine *fiber.App
	app    *fxtest.App
}

func (suite *HealthzSuite) SetupSuite() {
	suite.app = fxtest.New(suite.T(),
		fx.Provide(logger.Provide),
		fx.Provide(telemetry.ProvideNull),
		fx.Invoke(func(logger *zap.Logger, _ telemetry.Telemetery) {
			suite.engine = fiber.New()

			handler.Healthz{
				Logger: logger,
				Tracer: otel.GetTracerProvider().Tracer(""),
			}.Register(suite.engine.Group(""))
		}),
	).RequireStart()
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
