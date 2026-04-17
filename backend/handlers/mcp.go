package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type headerTransport struct {
	base   http.RoundTripper
	header http.Header
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	for k, vv := range t.header {
		for _, v := range vv {
			req.Header.Add(k, v)
		}
	}
	if t.base == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.base.RoundTrip(req)
}

type AddressResolveRequest struct {
	Address string `json:"address" binding:"required"`
}

func ResolveAddress(c *gin.Context) {
	var reqBody AddressResolveRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GOOGLE_MAPS_API_KEY environment variable is not set"})
		return
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-client", Version: "v1.0.0"}, nil)

	var customClient *http.Client
	if apiKey != "" {
		customClient = &http.Client{
			Transport: otelhttp.NewTransport(&headerTransport{
				base:   http.DefaultTransport,
				header: http.Header{"x-goog-api-key": []string{apiKey}},
			}),
			Timeout: 120 * time.Second,
		}
	} else {
		customClient = &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   120 * time.Second, // 2 minute timeout for AI operations
		}
	}

	transport := &mcp.StreamableClientTransport{
		Endpoint:   "https://mapstools.googleapis.com/mcp",
		HTTPClient: customClient,
	}

	cs, err := client.Connect(c.Request.Context(), transport, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Maps MCP server"})
		return
	}
	defer cs.Close()

	result, err := cs.CallTool(c.Request.Context(), &mcp.CallToolParams{
		Name:      "search_places",
		Arguments: map[string]any{"textQuery": reqBody.Address},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call MCP tool: " + err.Error()})
		return
	}

	for _, content := range result.Content {
		if textContent, ok := content.(*mcp.TextContent); ok {
			c.Data(http.StatusOK, "application/json", []byte(textContent.Text))
			return
		}
	}

	c.JSON(http.StatusOK, result)
}
