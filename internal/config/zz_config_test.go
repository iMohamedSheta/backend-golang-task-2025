package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetAndGet(t *testing.T) {
	cfg := New()

	cfg.Set("foo", "bar")
	val, err := cfg.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "bar", val)

	cfg.Set("x.y", "dot-value")
	val, err = cfg.Get("x.y")
	assert.NoError(t, err)
	assert.Equal(t, "dot-value", val)
}

func TestConfig_NestedKeys(t *testing.T) {
	cfg := New()
	cfg.Set("parent", map[string]any{"child": "value"})

	val, err := cfg.Get("parent.child")
	assert.NoError(t, err)
	assert.Equal(t, "value", val)
}

func TestConfig_Set_NestedKey(t *testing.T) {
	cfg := New()
	cfg.Set("a.b.c", "nested-value")

	val, err := cfg.Get("a.b.c")
	assert.NoError(t, err)
	assert.Equal(t, "nested-value", val)
}

func TestConfig_Set_OverwriteNested(t *testing.T) {
	cfg := New()
	cfg.Set("a.b.c", "v1")
	cfg.Set("a.b.c", "v2")

	val, err := cfg.Get("a.b.c")
	assert.NoError(t, err)
	assert.Equal(t, "v2", val)
}

func TestConfig_GetWithDefault(t *testing.T) {
	cfg := New()
	cfg.Set("key", "value")

	assert.Equal(t, "value", cfg.GetWithDefault("key", "fallback"))
	assert.Equal(t, "fallback", cfg.GetWithDefault("missing", "fallback"))
}

func TestConfig_GetTypedValues(t *testing.T) {
	cfg := New()
	cfg.Set("bool", true)
	cfg.Set("int", 42)
	cfg.Set("map", map[string]any{"x": 1})
	cfg.Set("duration", time.Second)
	cfg.Set("strings", []string{"a", "b"})

	assert.True(t, cfg.GetBool("bool", false))
	assert.Equal(t, 42, cfg.GetInt("int", 0))
	assert.Equal(t, 1, cfg.GetMap("map", nil)["x"])
	assert.Equal(t, time.Second, cfg.GetDuration("duration", 0))
	assert.Equal(t, []string{"a", "b"}, cfg.GetArrayOfStrings("strings", nil))
}

func TestConfig_Get_InvalidPath(t *testing.T) {
	cfg := New()
	cfg.Set("parent", "not a map")

	_, err := cfg.Get("parent.child")
	assert.Error(t, err)
}

func TestEnv(t *testing.T) {
	t.Setenv("EXISTING_ENV", "123")
	assert.Equal(t, "123", Env("EXISTING_ENV", "default"))
	assert.Equal(t, "default", Env("MISSING_ENV", "default"))
	assert.Nil(t, Env("MISSING_ENV", nil))
}

func TestRegisterAndApplyLoaders(t *testing.T) {
	cfg := New()

	Register(func(c *Config) {
		c.Set("loaded", true)
	})

	ApplyRegisteredLoaders(cfg)

	val, err := cfg.Get("loaded")
	assert.NoError(t, err)
	assert.Equal(t, true, val)
}
