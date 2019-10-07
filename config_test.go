package broadcast

import (
	"encoding/json"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCfg struct{ cfg string }

func (cfg *mockCfg) Get(name string) service.Config {
	if name == "same" || name == "jobs" {
		return cfg
	}

	return nil
}
func (cfg *mockCfg) Unmarshal(out interface{}) error { return json.Unmarshal([]byte(cfg.cfg), out) }

func Test_Config_Hydrate_Error(t *testing.T) {
	cfg := &mockCfg{cfg: `{"dead`}
	c := &Config{}

	assert.Error(t, c.Hydrate(cfg))
}

func Test_Config_Hydrate_OK(t *testing.T) {
	cfg := &mockCfg{cfg: `{"path":"/path"}`}
	c := &Config{}

	assert.NoError(t, c.Hydrate(cfg))
}

func Test_Config_Redis_Error(t *testing.T) {
	cfg := &mockCfg{cfg: `{"path":"/path","redis":{}}`}
	c := &Config{}

	assert.Error(t, c.Hydrate(cfg))
}

func Test_Config_Redis_OK(t *testing.T) {
	cfg := &mockCfg{cfg: `{"path":"/path","redis":{"addr":"localhost:6379"}}`}
	c := &Config{}

	assert.NoError(t, c.Hydrate(cfg))
}
