// Command mode parses a toml to generate.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kaneshin/genex"
)

var (
	pkg    = flag.String("pkg", "", "Package name to use in the generated code. (default \"main\")")
	srcDir = flag.String("path", "", "output file directory")
	output = flag.String("output", "", "output file name; default srcdir/routes_gen.go")
)

const (
	NAME   = "routes"
	OUTPUT = "routes_gen.go"
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(fmt.Sprintf("gen-%s: ", NAME))
	flag.Usage = Usage
	flag.Parse()
	if len(*pkg) == 0 {
		*pkg = "main"
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(3)
	}

	if len(*srcDir) > 0 {
		if d, err := filepath.Abs(*srcDir); err == nil {
			*srcDir = d
		}
	}

	if len(*srcDir) == 0 {
		*srcDir = filepath.Dir(os.Args[0])
	}

	os.Exit(run())
}

func run() int {
	var (
		args = flag.Args()
		mg   genex.Generator
		g    genex.Generator
	)

	// Print the header and package clause.
	mg.Printf("// Code generated by gen-routes.\n// DO NOT EDIT\n")
	mg.Printf("\npackage %s\n", *pkg)

	// Print the header and package clause.
	g.Printf("// Code generated by gen-routes.\n")
	// Set Routes from args
	list, filepaths := MustParseGlobs(args[:])
	for _, fp := range filepaths {
		g.Printf("// %s\n", fp)
	}
	g.Printf("// DO NOT EDIT\n")
	g.Printf("\npackage %s\n", *pkg)

	// Generate two source files.
	params := map[string]interface{}{}
	type Element struct {
		// Essentials
		Method, Path, Controller string
		// Controller
		ControllerPack, ControllerName, ControllerFunc string
		// Param
		Param []string
	}
	data := []Element{}
	statics := []Static{}
	htmls := []HTML{}
	for _, routes := range list {
		for _, route := range routes.Endpoint.Routes {
			elm := Element{}
			elm.Param = []string{}
			// Essentials
			elm.Method = route.Method
			elm.Path = path.Join("/", route.Path)
			elm.Controller = route.Controller

			// Controller
			_, t, fn := parsePackage(route.Controller)
			elm.ControllerFunc = fn
			_, elm.ControllerPack, elm.ControllerName = parsePackage(t)

			// Param
			elm.Param = parsePatternPath(elm.Path)
			// To be c.Param("%s")
			for i := 0; i < len(elm.Param); i++ {
				elm.Param[i] = fmt.Sprintf("c.Param(\"%s\")", elm.Param[i])
			}

			// append
			data = append(data, elm)
		}
		statics = append(statics, routes.Statics...)
		htmls = append(htmls, routes.HTML...)
	}
	params["data"] = data
	params["package"] = func() (ret []string) {
		ins := map[string]bool{}
		for _, routes := range list {
			for _, v := range routes.Packages {
				if ok, found := ins[v.Name]; found && ok {
					continue
				}
				ins[v.Name] = true
				ret = append(ret, v.Name)
			}
		}
		return
	}()
	params["static"] = statics
	params["html"] = htmls
	params["appName"] = filepath.Base(strings.Replace(*srcDir, "/internal/app", "", 1))

	// Write to routes_gen.go
	outputName := *output
	if outputName == "" {
		baseName := OUTPUT
		outputName = filepath.Join(*srcDir, strings.ToLower(baseName))
	}
	if !strings.HasSuffix(outputName, ".go") {
		outputName = outputName + ".go"
	}
	if err := execute(g, routesSource, outputName, params); err != nil {
		log.Fatalf("writing output: %s", err)
	}

	// Write to main_gen.go
	outputName = filepath.Join(*srcDir, "main_gen.go")
	var src string
	if v := os.Getenv("SOURCE_GAE"); v != "" {
		src = gaeSource
	} else {
		src = mainSource
	}
	if err := execute(mg, src, outputName, params); err != nil {
		log.Fatalf("writing output: %s", err)
	}

	return 0
}

func execute(gen genex.Generator, src, name string, params map[string]interface{}) error {
	// Format the output.
	funcMap := map[string]interface{}{
		"join": func(s []string) string {
			return strings.Join(s, ",")
		},
	}
	t := template.Must(template.New("").Funcs(funcMap).Parse(src))
	var buf bytes.Buffer
	t.Execute(&buf, params)
	gen.Printf("%s", buf.String())
	formattedSrc := gen.Format()

	// Write to file.
	return ioutil.WriteFile(name, formattedSrc, 0644)
}

const mainSource = `
import (
	"flag"
	"strings"
	"path/filepath"

	"github.com/gophergala2016/source/internal"
)

var _ = filepath.Separator

const (
	appName = "{{.appName}}"
)

var (
	port = flag.String("port", "", "Port number to listen the application. (default \"8888\")")
	sock = flag.String("sock", "", "Path to a UNIX socket to listen the application.")
)

func main() {
	flag.Parse()
	if len(*port) == 0 {
		*port = ":8888"
	}

	if !strings.HasPrefix(*port, ":") {
		*port = ":"+*port
	}

	// Initialize App
	internal.Init(appName, *port)

	router := router()
{{range .static}}
	router.Static("{{.Path}}", internal.JoinPath("{{.RelativePath}}")){{end}}
{{range .html}}
	router.LoadHTMLGlob(filepath.Clean(internal.JoinPath("{{.RelativePath}}")),
		filepath.Clean(internal.JoinPath("{{.Pattern}}"))){{end}}

	switch {
	case len(*sock) != 0:
		router.RunUnix(*sock)
	default:
		router.Run(*port)
	}
}
`

const gaeSource = `
import (
	"flag"
	"strings"
	"path/filepath"
	"net/http"

	"github.com/gophergala2016/source/internal"
	"github.com/gophergala2016/source/core/foundation"
	"github.com/gophergala2016/source/core/config"
)

var _ = filepath.Separator

const (
	appName = "{{.appName}}"
)

var (
	port = flag.String("port", "", "Port number to listen the application. (default \"8888\")")
	sock = flag.String("sock", "", "Path to a UNIX socket to listen the application.")
)

func init() {
	foundation.SetMode(foundation.ProdMode)
	config.Load()
	flag.Parse()
	if len(*port) == 0 {
		*port = ":8888"
	}

	if !strings.HasSuffix(*port, ":") {
		*port = ":"+*port
	}

	// Initialize App
	internal.Init(appName, *port)

	router := router()

	names, bodies := []string{}, []string{}
	for filename, _ := range _bindata {
		if !strings.HasSuffix(filename, "html") {
			continue
		}
		b, err := Asset(filename)
		if err != nil {
			panic(err)
		}
		names = append(names, filename)
		bodies = append(bodies, string(b))
	}
	router.LoadHTML(names, bodies)

	http.Handle("/", router.DefaultMux())
}
`

const routesSource = `
import (
	"github.com/gin-gonic/gin"

	"github.com/gophergala2016/source/core/foundation"
	"github.com/gophergala2016/source/core/net/context/accessor"
{{range .package}}
	"{{.}}"{{end}}
)

func router() *foundation.Engine {
	router := foundation.Default()

{{range .data}}
	// {{.Method}} {{.Path}}
	// {{.Controller}}
	router.{{.Method}}("{{.Path}}", func(c *gin.Context) {
		controller := &{{printf "%s." .ControllerPack}}{{.ControllerName}}{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "{{.ControllerFunc | printf "%s.%s" .ControllerName}}")
		controller.SetContext(ctx)
		controller.{{.ControllerFunc}}({{join .Param}})
	})
{{end}}

	return router
}
`
