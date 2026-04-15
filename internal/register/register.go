// Package register provides engine registration functionality.
package register

import (
	"github.com/tritueviet/search-agents/internal/engine"
	"github.com/tritueviet/search-agents/internal/engines/annasarchive"
	"github.com/tritueviet/search-agents/internal/engines/bing"
	"github.com/tritueviet/search-agents/internal/engines/bing_images"
	"github.com/tritueviet/search-agents/internal/engines/bing_news"
	"github.com/tritueviet/search-agents/internal/engines/brave"
	"github.com/tritueviet/search-agents/internal/engines/duckduckgo"
	"github.com/tritueviet/search-agents/internal/engines/duckduckgo_news"
	"github.com/tritueviet/search-agents/internal/engines/duckduckgo_videos"
	"github.com/tritueviet/search-agents/internal/engines/google"
	"github.com/tritueviet/search-agents/internal/engines/grokipedia"
	"github.com/tritueviet/search-agents/internal/engines/mojeek"
	"github.com/tritueviet/search-agents/internal/engines/openlibrary"
	"github.com/tritueviet/search-agents/internal/engines/wikipedia"
	"github.com/tritueviet/search-agents/internal/engines/yahoo"
	"github.com/tritueviet/search-agents/internal/engines/yahoo_news"
	"github.com/tritueviet/search-agents/internal/engines/yandex"
	"github.com/tritueviet/search-agents/internal/httpclient"
)

// DefaultEngines registers all available search engines.
func DefaultEngines(client *httpclient.Client, registry *engine.Registry) {
	// Text engines (9 engines)
	registry.Register(engine.CategoryText, "duckduckgo", func(c *httpclient.Client) engine.SearchEngine {
		return duckduckgo.New(c)
	})
	registry.Register(engine.CategoryText, "bing", func(c *httpclient.Client) engine.SearchEngine {
		return bing.New(c)
	})
	registry.Register(engine.CategoryText, "google", func(c *httpclient.Client) engine.SearchEngine {
		return google.New(c)
	})
	registry.Register(engine.CategoryText, "brave", func(c *httpclient.Client) engine.SearchEngine {
		return brave.New(c)
	})
	registry.Register(engine.CategoryText, "yahoo", func(c *httpclient.Client) engine.SearchEngine {
		return yahoo.New(c)
	})
	registry.Register(engine.CategoryText, "yandex", func(c *httpclient.Client) engine.SearchEngine {
		return yandex.New(c)
	})
	registry.Register(engine.CategoryText, "wikipedia", func(c *httpclient.Client) engine.SearchEngine {
		return wikipedia.New(c)
	})
	registry.Register(engine.CategoryText, "grokipedia", func(c *httpclient.Client) engine.SearchEngine {
		return grokipedia.New(c)
	})
	registry.Register(engine.CategoryText, "mojeek", func(c *httpclient.Client) engine.SearchEngine {
		return mojeek.New(c)
	})

	// Images engines (1 engine)
	registry.Register(engine.CategoryImages, "bing", func(c *httpclient.Client) engine.SearchEngine {
		return bing_images.New(c)
	})

	// Videos engines (1 engine)
	registry.Register(engine.CategoryVideos, "duckduckgo", func(c *httpclient.Client) engine.SearchEngine {
		return duckduckgo_videos.New(c)
	})

	// News engines (3 engines)
	registry.Register(engine.CategoryNews, "bing", func(c *httpclient.Client) engine.SearchEngine {
		return bing_news.New(c)
	})
	registry.Register(engine.CategoryNews, "duckduckgo", func(c *httpclient.Client) engine.SearchEngine {
		return duckduckgo_news.New(c)
	})
	registry.Register(engine.CategoryNews, "yahoo", func(c *httpclient.Client) engine.SearchEngine {
		return yahoo_news.New(c)
	})

	// Books engines (2 engines)
	registry.Register(engine.CategoryBooks, "openlibrary", func(c *httpclient.Client) engine.SearchEngine {
		return openlibrary.New(c)
	})
	registry.Register(engine.CategoryBooks, "annasarchive", func(c *httpclient.Client) engine.SearchEngine {
		return annasarchive.New(c)
	})
}
