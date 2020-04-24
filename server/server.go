package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/polisgo2020/search-K1ta/database"
	"github.com/polisgo2020/search-K1ta/revindex"
	"github.com/polisgo2020/search-K1ta/server/templates"
	"github.com/sirupsen/logrus"
	"net/http"
)

type App struct {
	*database.DB
}

type config struct {
	Addr string `env:"POLISGO_ADDR" envDefault:"localhost:8080"`
}

func (a *App) index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (a *App) search(c echo.Context) error {
	logrus.Infoln(c.Request().RemoteAddr, "Phrase:", c.QueryParam("phrase"))
	res, err := revindex.FindInDb(c.QueryParam("phrase"), a.DB)
	if err != nil {
		logrus.Error(c.Request().RemoteAddr, "Error:", err)
		return c.Render(http.StatusInternalServerError, "index.html", res)
	}
	logrus.Infoln(c.Request().RemoteAddr, "Result:", res)
	return c.Render(http.StatusOK, "index.html", res)
}

func Start(addr string, db *database.DB) error {
	// configure logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		PadLevelText:  true,
	})

	// create server
	e := echo.New()
	e.Pre(middleware.AddTrailingSlash())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logrus.Infoln(c.Request().RemoteAddr, c.Request().Method, c.Request().URL)
			return next(c)
		}
	})
	e.Use(middleware.Recover())
	app := App{db}

	// add page renderer
	renderer, err := templates.Init()
	if err != nil {
		return err
	}
	e.Renderer = renderer

	// add routes
	e.Add(echo.GET, "/", app.index)
	e.Add(echo.GET, "/search/", app.search)
	e.Static("/static", "server/static")

	// start server
	err = e.Start(addr)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
