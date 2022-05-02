package azuremonitor

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/grafana/grafana-azure-sdk-go/azcredentials"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"

	"github.com/grafana/grafana/pkg/infra/httpclient"
	"github.com/grafana/grafana/pkg/setting"
	"github.com/grafana/grafana/pkg/tsdb/azuremonitor/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpClient_AzureCredentials(t *testing.T) {
	model := types.DatasourceInfo{
		Credentials: &azcredentials.AzureManagedIdentityCredentials{},
	}

	cfg := &setting.Cfg{}
	provider := &fakeHttpClientProvider{}

	t.Run("should have Azure middleware when scopes provided", func(t *testing.T) {
		route := types.AzRoute{
			Scopes: []string{"https://management.azure.com/.default"},
		}

		_, err := newHTTPClient(route, model, cfg, provider)
		require.NoError(t, err)

		require.NotNil(t, provider.opts)
		require.NotNil(t, provider.opts.Middlewares)
		assert.Len(t, provider.opts.Middlewares, 1)
	})

	t.Run("should not have Azure middleware when scopes are not provided", func(t *testing.T) {
		route := types.AzRoute{
			Scopes: []string{},
		}

		_, err := newHTTPClient(route, model, cfg, provider)
		require.NoError(t, err)

		assert.NotNil(t, provider.opts)

		if provider.opts.Middlewares != nil {
			assert.Len(t, provider.opts.Middlewares, 0)
		}
	})
}

type fakeHttpClientProvider struct {
	httpclient.Provider

	opts sdkhttpclient.Options
}

func (p *fakeHttpClientProvider) New(opts ...sdkhttpclient.Options) (*http.Client, error) {
	p.opts = opts[0]
	return nil, nil
}

func (p *fakeHttpClientProvider) GetTransport(opts ...sdkhttpclient.Options) (http.RoundTripper, error) {
	p.opts = opts[0]
	return nil, nil
}

func (p *fakeHttpClientProvider) GetTLSConfig(opts ...sdkhttpclient.Options) (*tls.Config, error) {
	p.opts = opts[0]
	return nil, nil
}
