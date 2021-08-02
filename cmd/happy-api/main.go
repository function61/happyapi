package main

import (
	"embed"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/function61/gokit/app/aws/lambdautils"
	"github.com/function61/gokit/app/dynversion"
	"github.com/function61/gokit/log/logex"
	"github.com/function61/gokit/net/http/httputils"
	"github.com/function61/gokit/os/osutil"
	"github.com/function61/happy-api/pkg/turbocharger/turbochargerapp"
	"github.com/function61/happy-api/static"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

//go:embed item.html
var templates embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())

	// AWS Lambda doesn't support giving argv, so we use an ugly hack to detect when
	// we're in Lambda
	if lambdautils.InLambda() {
		lambda.StartHandler(lambdautils.NewLambdaHttpHandlerAdapter(httpHandler()))
		return
	}

	app := &cobra.Command{
		Use:     os.Args[0],
		Short:   "Happiness as a service",
		Version: dynversion.Version,
		Args:    cobra.NoArgs,
		Run: func(_ *cobra.Command, args []string) {
			srv := &http.Server{
				Addr:    ":80",
				Handler: httpHandler(),
			}

			osutil.ExitIfError(httputils.CancelableServer(
				osutil.CancelOnInterruptOrTerminate(logex.StandardLogger()),
				srv,
				func() error { return srv.ListenAndServe() }))

		},
	}

	app.AddCommand(newEntry())
	app.AddCommand(turbochargerapp.StaticFilesExportEntrypoint(static.Files))

	osutil.ExitIfError(app.Execute())
}

func httpHandler() http.Handler {
	uiTpl, err := template.ParseFS(templates, "item.html")
	if err != nil {
		panic(err)
	}

	happiness, err := static.Files.ReadDir("images")
	if err != nil {
		panic(err)
	}

	routes := mux.NewRouter()

	redirectToRandomItem := func(w http.ResponseWriter, r *http.Request) {
		idx := randBetween(0, len(happiness)-1)

		http.Redirect(w, r, "/happy/"+fileIdFromFilename(happiness[idx].Name()), http.StatusFound)
	}

	routes.PathPrefix("/happy/static").Handler(turbochargerapp.FileHandler("/happy/static", static.Files))

	routes.HandleFunc("/happy", redirectToRandomItem)
	routes.HandleFunc("/happy/", redirectToRandomItem)

	routes.HandleFunc("/happy/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		attribution, err := findAttributionFromExifArtist(id)
		if err != nil { // assuming error is ErrNotExist
			if os.IsNotExist(err) {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")

		_ = uiTpl.Execute(w, struct {
			ImgSrc      string
			Attribution string
		}{
			ImgSrc:      "/happy/static/images/" + id + ".jpg",
			Attribution: attribution,
		})
	})

	return routes
}
