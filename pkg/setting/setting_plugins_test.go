package setting

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPluginSettings(t *testing.T) {
	cfg := NewCfg()
	sec, err := cfg.Raw.NewSection("plugin")
	require.NoError(t, err)
	_, err = sec.NewKey("key", "value")
	require.NoError(t, err)

	sec, err = cfg.Raw.NewSection("plugin.plugin")
	require.NoError(t, err)
	_, err = sec.NewKey("key1", "value1")
	require.NoError(t, err)
	_, err = sec.NewKey("key2", "value2")
	require.NoError(t, err)

	sec, err = cfg.Raw.NewSection("plugin.plugin2")
	require.NoError(t, err)
	_, err = sec.NewKey("key3", "value3")
	require.NoError(t, err)
	_, err = sec.NewKey("key4", "value4")
	require.NoError(t, err)

	sec, err = cfg.Raw.NewSection("other")
	require.NoError(t, err)
	_, err = sec.NewKey("keySomething", "whatever")
	require.NoError(t, err)

	ps := extractPluginSettings(cfg.Raw.Sections())
	require.Len(t, ps, 2)
	require.Len(t, ps["plugin"], 2)
	require.Equal(t, ps["plugin"]["key1"], "value1")
	require.Equal(t, ps["plugin"]["key2"], "value2")
	require.Len(t, ps["plugin2"], 2)
	require.Equal(t, ps["plugin2"]["key3"], "value3")
	require.Equal(t, ps["plugin2"]["key4"], "value4")
}

func Test_readPluginSettings(t *testing.T) {
	t.Run("should parse plugin ids", func(t *testing.T) {
		cfg := NewCfg()
		sec, err := cfg.Raw.NewSection("plugins")
		require.NoError(t, err)
		_, err = sec.NewKey("disable_plugins", "plugin1,plugin2")
		require.NoError(t, err)

		_, err = sec.NewKey("plugin_catalog_hidden_plugins", "plugin3")
		require.NoError(t, err)

		_, err = sec.NewKey("hide_angular_deprecation", "a,b,c")
		require.NoError(t, err)

		err = cfg.readPluginSettings(cfg.Raw)
		require.NoError(t, err)
		require.Equal(t, []string{"plugin1", "plugin2"}, cfg.DisablePlugins)
		require.Equal(t, []string{"plugin3", "plugin1", "plugin2"}, cfg.PluginCatalogHiddenPlugins)
		require.Equal(t, []string{"a", "b", "c"}, cfg.HideAngularDeprecation)
	})
}
