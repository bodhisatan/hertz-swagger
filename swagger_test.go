package hertzSwagger

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/route"

	"github.com/swaggo/swag"

	"github.com/cloudwego/hertz/pkg/common/test/assert"
	swaggerFiles "github.com/swaggo/files"
)

type mockedSwag struct{}

func (s *mockedSwag) ReadDoc() string {
	return `{
}`
}

func TestWrapHandler(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.GET("/*any", WrapHandler(swaggerFiles.Handler, URL("https://github.com/bodhisatan/hertz-swagger")))

	w := ut.PerformRequest(router, http.MethodGet, "/index.html", nil)
	resp := w.Result()
	assert.DeepEqual(t, http.StatusOK, resp.StatusCode())
}

func TestCustomWrapHandler(t *testing.T) {
	router := route.NewEngine(config.NewOptions([]config.Option{}))
	router.Any("/*any", CustomWrapHandler(&Config{}, swaggerFiles.Handler))

	w1 := ut.PerformRequest(router, http.MethodGet, "/index.html", nil)
	assert.DeepEqual(t, http.StatusOK, w1.Code)
	assert.DeepEqual(t, string(w1.Header().ContentType()), "text/html; charset=utf-8")

	w2 := ut.PerformRequest(router, http.MethodGet, "/doc.json", nil)
	assert.DeepEqual(t, http.StatusInternalServerError, w2.Code)

	doc := &mockedSwag{}
	swag.Register(swag.Name, doc)

	w3 := ut.PerformRequest(router, http.MethodGet, "/doc.json", nil)
	assert.DeepEqual(t, http.StatusOK, w3.Code)
	assert.DeepEqual(t, string(w3.Header().ContentType()), "application/json; charset=utf-8")

	// Perform body rendering validation
	w3Body, err := ioutil.ReadAll(w3.Body)
	assert.Nil(t, err)
	assert.DeepEqual(t, doc.ReadDoc(), string(w3Body))

	w4 := ut.PerformRequest(router, http.MethodGet, "/favicon-16x16.png", nil)
	assert.DeepEqual(t, http.StatusOK, w4.Code)
	assert.DeepEqual(t, string(w4.Header().ContentType()), "image/png")

	w5 := ut.PerformRequest(router, http.MethodGet, "/swagger-ui.css", nil)
	assert.DeepEqual(t, http.StatusOK, w5.Code)
	assert.DeepEqual(t, string(w5.Header().ContentType()), "text/css; charset=utf-8")

	w6 := ut.PerformRequest(router, http.MethodGet, "/swagger-ui-bundle.js", nil)
	assert.DeepEqual(t, http.StatusOK, w6.Code)
	assert.DeepEqual(t, string(w6.Header().ContentType()), "application/javascript")

	assert.DeepEqual(t, http.StatusNotFound, ut.PerformRequest(router, http.MethodGet, "/notfound", nil).Code)

	assert.DeepEqual(t, http.StatusMethodNotAllowed, ut.PerformRequest(router, http.MethodPost, "/index.html", nil).Code)

	assert.DeepEqual(t, http.StatusMethodNotAllowed, ut.PerformRequest(router, http.MethodPut, "/index.html", nil).Code)
}

func TestURL(t *testing.T) {
	cfg := Config{}

	expected := "https://github.com/swaggo/http-swagger"
	configFunc := URL(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.URL)
}

func TestDocExpansion(t *testing.T) {
	var cfg Config

	expected := "list"
	configFunc := DocExpansion(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.DocExpansion)

	expected = "full"
	configFunc = DocExpansion(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.DocExpansion)

	expected = "none"
	configFunc = DocExpansion(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.DocExpansion)
}

func TestDeepLinking(t *testing.T) {
	var cfg Config
	assert.DeepEqual(t, false, cfg.DeepLinking)

	configFunc := DeepLinking(true)
	configFunc(&cfg)
	assert.DeepEqual(t, true, cfg.DeepLinking)

	configFunc = DeepLinking(false)
	configFunc(&cfg)
	assert.DeepEqual(t, false, cfg.DeepLinking)
}

func TestDefaultModelsExpandDepth(t *testing.T) {
	var cfg Config

	assert.DeepEqual(t, 0, cfg.DefaultModelsExpandDepth)

	expected := -1
	configFunc := DefaultModelsExpandDepth(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.DefaultModelsExpandDepth)

	expected = 1
	configFunc = DefaultModelsExpandDepth(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.DefaultModelsExpandDepth)
}

func TestInstanceName(t *testing.T) {
	var cfg Config

	assert.DeepEqual(t, "", cfg.InstanceName)

	expected := swag.Name
	configFunc := InstanceName(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.InstanceName)

	expected = "custom_name"
	configFunc = InstanceName(expected)
	configFunc(&cfg)
	assert.DeepEqual(t, expected, cfg.InstanceName)
}

func TestPersistAuthorization(t *testing.T) {
	var cfg Config
	assert.DeepEqual(t, false, cfg.PersistAuthorization)

	configFunc := PersistAuthorization(true)
	configFunc(&cfg)
	assert.DeepEqual(t, true, cfg.PersistAuthorization)

	configFunc = PersistAuthorization(false)
	configFunc(&cfg)
	assert.DeepEqual(t, false, cfg.PersistAuthorization)
}

func TestOauth2DefaultClientID(t *testing.T) {
	var cfg Config
	assert.DeepEqual(t, "", cfg.Oauth2DefaultClientID)

	configFunc := Oauth2DefaultClientID("default_client_id")
	configFunc(&cfg)
	assert.DeepEqual(t, "default_client_id", cfg.Oauth2DefaultClientID)

	configFunc = Oauth2DefaultClientID("")
	configFunc(&cfg)
	assert.DeepEqual(t, "", cfg.Oauth2DefaultClientID)
}
