package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/config"
	"github.com/1995parham/saf/internal/infra/http/handler"
	"github.com/1995parham/saf/internal/infra/http/request"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type EventSuite struct {
	suite.Suite

	engine *fiber.App
}

func (suite *EventSuite) SetupSuite() {
	suite.engine = fiber.New()

	cfg := config.Provide()

	cmq, err := cmq.Provide(nil, cfg.NATS, zap.NewNop())
	suite.Require().NoError(err)

	suite.Require().NoError(cmq.Streams())

	handler.Event{
		CMQ:    cmq,
		Logger: zap.NewNop(),
		Tracer: trace.NewNoopTracerProvider().Tracer(""),
	}.Register(suite.engine.Group(""))
}

func (suite *EventSuite) TestHandler() {
	require := suite.Require()

	payload, err := json.Marshal(request.Event{
		Subject: "hello",
		ID:      "",
		Payload: []byte("from the otherside"),
	})
	require.NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(payload))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	resp, err := suite.engine.Test(req)
	require.NoError(err)
	require.Equal(http.StatusOK, resp.StatusCode)
	require.NoError(resp.Body.Close())
}

func TestEventSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EventSuite))
}
