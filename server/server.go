package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/vlasad/redislike/cache"
)

type Server struct {
	cache *cache.Cache
}

func New() Server {
	return Server{
		cache: cache.New(),
	}
}

func (s *Server) Start(port int) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/keys", s.keys)
	e.DELETE("/remove/:key", s.remove)
	e.POST("/ttl/:key", s.setTTL)

	e.POST("/set", s.set)
	e.GET("/get/:key", s.get)

	e.POST("/push", s.push)
	e.GET("/pop/:key", s.pop)

	e.POST("/hset", s.hset)
	e.GET("/hget/:key/:field", s.hget)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

// curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/keys
func (s *Server) keys(c echo.Context) error {
	return response(c, http.StatusOK, s.cache.Keys(), nil)
}

// curl -i -w "\n" -X DELETE -H 'Content-Type: application/json' localhost:8080/remove/abc
func (s *Server) remove(c echo.Context) error {
	s.cache.Remove(c.Param("key"))
	return response(c, http.StatusOK, "ok", nil)
}

// Set TTL in seconds
// curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"value":10}' localhost:8080/ttl/abc
func (s *Server) setTTL(c echo.Context) error {
	data := &struct {
		Value int `json:"value"`
	}{}
	if err := c.Bind(data); err != nil {
		return response(c, http.StatusBadRequest, nil, err)
	}

	err := s.cache.SetTTL(c.Param("key"), time.Duration(data.Value)*time.Second)
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}

	return response(c, http.StatusOK, "ok", nil)
}

// curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"abc", "value":"test value"}' localhost:8080/set
func (s *Server) set(c echo.Context) error {
	data := &struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}

	if err := c.Bind(data); err != nil {
		return response(c, http.StatusBadRequest, nil, err)
	}

	s.cache.Set(data.Key, data.Value)

	return response(c, http.StatusOK, "ok", nil)
}

// curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/get/abc
func (s *Server) get(c echo.Context) error {
	value, err := s.cache.Get(c.Param("key"))
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}
	return response(c, http.StatusOK, value, nil)
}

// curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"list", "value":["v1","v2"]}' localhost:8080/push
func (s *Server) push(c echo.Context) error {
	data := &struct {
		Key   string   `json:"key"`
		Value []string `json:"value"`
	}{}
	if err := c.Bind(data); err != nil {
		return response(c, http.StatusBadRequest, nil, err)
	}

	err := s.cache.Push(data.Key, data.Value...)
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}

	return response(c, http.StatusOK, "ok", nil)
}

// curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/pop/list
func (s *Server) pop(c echo.Context) error {
	value, err := s.cache.Pop(c.Param("key"))
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}
	return response(c, http.StatusOK, value, nil)
}

// curl -i -w "\n" -X POST -H 'Content-Type: application/json' -d '{"key":"dict", "field":"f1", "value": "v1"}' localhost:8080/hset
func (s *Server) hset(c echo.Context) error {
	data := &struct {
		Key   string `json:"key"`
		Field string `json:"field"`
		Value string `json:"value"`
	}{}
	if err := c.Bind(data); err != nil {
		return response(c, http.StatusBadRequest, nil, err)
	}

	err := s.cache.Hset(data.Key, data.Field, data.Value)
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}

	return response(c, http.StatusOK, "ok", nil)
}

// curl -i -w "\n" -X GET -H 'Content-Type: application/json' localhost:8080/hget/dict/f1
func (s *Server) hget(c echo.Context) error {
	value, err := s.cache.Hget(c.Param("key"), c.Param("field"))
	if err != nil {
		return response(c, http.StatusInternalServerError, nil, err)
	}
	return response(c, http.StatusOK, value, nil)
}

func response(c echo.Context, status int, value interface{}, err error) error {
	data := struct {
		Value interface{} `json:"value,omitempty"`
		Error error       `json:"error,omitempty"`
	}{
		value, err,
	}

	return c.JSON(status, data)
}
