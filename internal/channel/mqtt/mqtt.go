package mqtt

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/1995parham/saf/internal/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Config struct {
	URL       string        `koanf:"url"`
	ClientID  string        `koanf:"clientid"`
	Username  string        `koanf:"username"`
	Password  string        `koanf:"password"`
	KeepAlive time.Duration `koanf:"keepalive"`
}

const (
	PingTimeout       = 1 * time.Second
	ReconnectInterval = 10 * time.Second
	DisconnectTimeout = 250
)

func (p *MQTT) options() *mqtt.ClientOptions {
	clientID := p.cfg.ClientID
	if clientID == "" {
		var err error

		clientID, err = os.Hostname()
		if err != nil {
			p.logger.Fatal("hostname fetching failed, specify a client id", zap.Error(err))
		}
	}

	opts := mqtt.NewClientOptions().AddBroker(p.cfg.URL).SetClientID(clientID)

	if p.cfg.Username != "" {
		opts.SetUsername(p.cfg.Username)
	}

	if p.cfg.Password != "" {
		opts.SetPassword(p.cfg.Password)
	}

	opts.SetKeepAlive(p.cfg.KeepAlive)
	opts.SetPingTimeout(PingTimeout)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(ReconnectInterval)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		p.logger.Error("mqtt connection lost", zap.Error(err))
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		p.logger.Debug("mqtt reconnect")
	})

	return opts
}

// MQTT is a plugin for the saf app. this plugin consumes event
// and publish them over mqtt.
type MQTT struct {
	ch     <-chan model.ChanneledEvent
	logger *zap.Logger
	tracer trace.Tracer
	cfg    Config
	client mqtt.Client
}

func (p *MQTT) Init(logger *zap.Logger, tracer trace.Tracer, cfg interface{}) {
	p.logger = logger
	p.tracer = tracer

	// nolint: exhaustivestruct
	dc := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc()),
		Metadata:         nil,
		Result:           &p.cfg,
		WeaklyTypedInput: true,
		TagName:          "koanf",
	}

	d, err := mapstructure.NewDecoder(dc)
	if err != nil {
		p.logger.Fatal("failed to create mapstructure decoder", zap.Error(err))
	}

	if err := d.Decode(cfg); err != nil {
		p.logger.Fatal("failed to decode configuration", zap.Error(err))
	}

	p.client = mqtt.NewClient(p.options())

	if token := p.client.Connect(); token.Wait() && token.Error() != nil {
		p.logger.Fatal("mqtt connection failed", zap.Error(token.Error()))
	}
}

func (p *MQTT) Run() {
	for i := 0; i < 10*runtime.GOMAXPROCS(0); i++ {
		go func() {
			for e := range p.ch {
				_, span := p.tracer.Start(e.Context, "channels.mqtt")

				p.client.Publish(fmt.Sprintf("saf/%s", e.Subject), 1, true, e.Payload)

				span.End()
			}
		}()
	}
}

func (p *MQTT) Name() string {
	return "mqtt"
}

func (p *MQTT) SetChannel(c <-chan model.ChanneledEvent) {
	p.ch = c
}
