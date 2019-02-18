package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/macaron.v1"
  "github.com/go-macaron/session"
  
	"github.com/99designs/gqlgen/handler"
	"github.com/KellyLSB/demondin/graphql"
)

func main() {
	// Load .env file from repo
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load a macaron daemon
	m := macaron.Classic()
	// Injected automatically by macaron.Classic()
	//m.Use(macaron.Logger())
	//m.Use(macaron.Recovery())
	//m.Use(macaron.Static("public"))
	
	// Handle Session Cookies
	m.Use(session.Sessioner(session.Options{
    // Name of provider. Default is "memory".
    Provider:       "memory",
    // Provider configuration, it's corresponding to provider.
    ProviderConfig: "",
    // Cookie name to save session ID. Default is "MacaronSession".
    CookieName:     "demondin",
    // Cookie path to store. Default is "/".
    CookiePath:     "/",
    // GC interval time in seconds. Default is 3600.
    Gclifetime:     3600,
    // Max life time in seconds. Default is whatever GC interval time is.
    Maxlifetime:    3600,
    // Use HTTPS only. Default is false.
    Secure:         false,
    // Cookie life time. Default is 0.
    CookieLifeTime: 0,
    // Cookie domain name. Default is empty.
    Domain:         "",
    // Session ID length. Default is 16.
    IDLength:       16,
    // Configuration section name. Default is "session".
    Section:        "session",
  }))


	// Handle Templating
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Extensions: []string{".tmpl", ".html"},
		Directory:  "templates/default",
		IndentXML:  true,
		IndentJSON: true,

		Funcs: []template.FuncMap{map[string]interface{}{
			"AppName": func() string {
				return "DemonDin"
			},
			"AppVer": func() string {
				return "0.0.1"
			},
		}},
	}))

	// Routes
	m.Group("/shop", func() {
		m.Group("/keeper", func() {
			m.Get("/", func(ctx *macaron.Context) {
        ctx.HTML(200, "shop/keeper/items")
			})
		})

		m.Get("/", func(ctx *macaron.Context) {
			ctx.HTML(200, "shop/items")
		})
	})

	m.Get("/", myHandler)
	
	// Graphql
	m.Get("/playground", handler.Playground("GraphQL playground", "/query"))
	m.Get("/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}})))

	// Start our macaron daemon
	log.Println("Server is running...")
	log.Printf("connect to http://%s/playground for GraphQL playground", os.Getenv("HOSTPORT"))
	log.Println(http.ListenAndServe(os.Getenv("HOSTPORT"), m))
}

func myHandler(ctx *macaron.Context) string {
	return "the request path is: " + ctx.Req.RequestURI
}
