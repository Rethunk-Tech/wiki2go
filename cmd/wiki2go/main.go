package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	embed "github.com/13rac1/goldmark-embed"
	log "github.com/AlbinoGeek/logxi/v1"
	"github.com/adrg/frontmatter"
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/hashtag"
	"go.abhg.dev/goldmark/wikilink"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var rootCmd = &cobra.Command{
	Aliases:    []string{"w2g"},
	Example:    `wiki2go build -o /path/to/output .`,
	Long:       `wiki2go creates a static wiki from markdown files.`,
	Run:        rootRun,
	Short:      "wiki2go creates a static wiki from markdown files",
	SuggestFor: []string{"wiki"},
	Use:        "wiki2go [flags] {build|edit|new|serve} [args]",
}

var buildCmd = &cobra.Command{
	Aliases:    []string{"b"},
	Example:    `wiki2go build -o /path/to/output .`,
	Long:       `Builds the static wiki from the markdown files in the specified path. If no path is specified, the current directory is used.`,
	Run:        buildRun,
	Short:      "Builds the static wiki",
	SuggestFor: []string{"compile"},
	Use:        "build [flags] [path]",
}

var editCmd = &cobra.Command{
	Aliases: []string{"e"},
	Example: `wiki2go edit`,
	Long:    `Opens the wiki's content folder in the application specified by EDITOR, otherwise the system default editor.`,
	Run:     editRun,
	Short:   "Edit the wiki",
	Use:     "edit [flags] [path]",
}

var newCmd = &cobra.Command{
	Aliases: []string{"n", "create"},
	Example: `wiki2go new`,
	Long:    `Creates a new markdown file in the wiki's content folder.`,
	Run:     newRun,
	Short:   "Create a new markdown file",
	Use:     "new [flags] [path]",
}

var serveCmd = &cobra.Command{
	Aliases: []string{"s"},
	Example: `wiki2go serve`,
	Long:    `Starts a web server to serve the static wiki.`,
	Run:     serveRun,
	Short:   "Serve the static wiki",
	Use:     "serve [flags] [path]",
}

func setupCommands() error {
	// Build
	buildCmd.Flags().StringArrayP("exclude", "x", []string{},
		"The files, directories, or wildcards to exclude from the build.")
	buildCmd.Flags().StringArrayP("include", "i", []string{},
		"The files, directories, or wildcards to include in the build.")
	buildCmd.Flags().StringP("output", "o", "build",
		"The output directory for the static wiki.")
	buildCmd.Flags().BoolP("watch", "w", false,
		"Watch the content directory for changes and rebuild the wiki.")
	rootCmd.AddCommand(buildCmd)

	// Edit
	editCmd.Flags().StringP("editor", "e", "",
		"The editor to use for editing the wiki files.")
	rootCmd.AddCommand(editCmd)

	// New
	rootCmd.AddCommand(newCmd)

	// Serve
	serveCmd.Flags().StringP("hostname", "H", "localhost",
		"The host to serve the static wiki on.")
	serveCmd.Flags().Int16P("port", "P", 8123,
		"The port to serve the static wiki on.")

	buildCmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		// skip "output" flag
		if flag.Name == "output" {
			return
		}

		serveCmd.Flags().AddFlag(flag)
	})

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	rootCmd.PersistentFlags().Lookup("help").Hidden = true
	rootCmd.AddCommand(serveCmd)
	return nil
}

func main() {
	if err := setupCommands(); err != nil {
		log.Fatal("failed to setup commands", "error", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("failed to execute root command", "error", err)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	log.Fatal("failed to access feature \"rootRun\" - not implemented")
}

func editRun(cmd *cobra.Command, args []string) {
	log.Fatal("failed to access feature \"editRun\" - not implemented")
}

func newRun(cmd *cobra.Command, args []string) {
	log.Fatal("failed to access feature \"newRun\" - not implemented")
}

//
//
//

func getPaths(cmd *cobra.Command) ([]string, []string) {
	includePaths, err := cmd.Flags().GetStringArray("include")
	if err != nil {
		log.Fatal("failed to retrieve include paths flag", "error", err)
	}

	for i, includePath := range includePaths {
		// Resolve to absolute and sure exists
		if includePath, err = filepath.Abs(includePath); err != nil {
			log.Fatal("failed to resolve include path", "error", err, "path", includePath)
		}

		if _, err = os.Stat(includePath); os.IsNotExist(err) {
			log.Fatal("include path does not exist", "path", includePath)
		}

		if err != nil {
			log.Fatal("failed to stat include path", "error", err, "path", includePath)
		}

		includePaths[i] = includePath
	}

	excludePaths, err := cmd.Flags().GetStringArray("exclude")
	if err != nil {
		log.Fatal("failed to retrieve exclude paths flag", "error", err)
	}

	return includePaths, excludePaths
}

//
//
//

func buildRun(cmd *cobra.Command, args []string) {
	includePaths, excludePaths := getPaths(cmd)
	log.Fatal("failed to access feature \"buildRun\" - not implemented",
		"includePaths", includePaths, "excludePaths", excludePaths)
	// build(includePaths, excludePaths)
}

var englishTitleCase = cases.Title(language.English)

func serveRun(cmd *cobra.Command, args []string) {
	e := echo.New()
	e.HideBanner = true
	server = &Server{
		CSS:        "",
		Echo:       e,
		KnownFiles: make([]*KnownFile, 0),
	}

	includePaths, excludePaths := getPaths(cmd)
	// build(includePaths, excludePaths)

	host, err := cmd.Flags().GetString("hostname")
	if err != nil {
		log.Fatal("failed to retrieve hostname flag", "error", err)
	}

	port, err := cmd.Flags().GetInt16("port")
	if err != nil {
		log.Fatal("failed to retrieve port flag", "error", err)
	}

	watchMode, err := cmd.Flags().GetBool("watch")
	if err != nil {
		log.Fatal("failed to retrieve watch flag", "error", err)
	}

	if watchMode {
		watch(includePaths, excludePaths)
		defer watcher.Close()
	}

	e.GET("/", func(c echo.Context) error {
		return server.handleGET(c, "home", "")
	})

	e.GET("/:page_name", func(c echo.Context) error {
		return server.handleGET(c, c.Param("page_name"), "")
	})

	e.GET("/:page_name/*", func(c echo.Context) error {
		return server.handleGET(c, c.Param("page_name"), c.Param("*"))
	})

	// Start server and block
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}

func makeNameCanonical(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	return strings.ToLower(s)
}

func makeNamePretty(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	return englishTitleCase.String(s)
}

type KnownFile struct {
	// Path is the absolute path to the file on disk
	Path string

	// Slug is the canonical name of the file relative to the server root
	Slug string

	// Title is the pretty name of the file, suitable for display
	Title string
}

var server *Server

type Server struct {
	CSS        string
	Echo       *echo.Echo
	KnownFiles []*KnownFile
}

func (s *Server) handleGET(c echo.Context, desiredPage, args string) error {
	schema := "http"
	if c.Request().TLS != nil {
		schema = "https"
	}

	canonicalUrl := fmt.Sprintf("%s://%s/%s",
		schema, s.Echo.ListenerAddr().String(),
		makeNameCanonical(desiredPage))

	var resolved *KnownFile
	for _, knownFile := range s.KnownFiles {
		if knownFile.Slug == desiredPage {
			resolved = knownFile
			break
		}
	}

	if resolved == nil {
		return c.JSON(404, map[string]interface{}{
			"error": "page not found",
			"request": map[string]interface{}{
				"args": args,
				"page": desiredPage,
			},
			"url": canonicalUrl,
		})
	}

	return s.render(c, resolved)
}

type r struct{}

var _hash = []byte("#")

func (r) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	dest := make([]byte, 512)

	var i int
	if len(n.Target) > 0 {
		i += copy(dest, []byte(makeNameCanonical(string(n.Target))))
	}

	if len(n.Fragment) > 0 {
		i += copy(dest[i:], _hash)
		i += copy(dest[i:], n.Fragment)
	}

	return dest[:i], nil
}

func (s *Server) render(c echo.Context, file *KnownFile) error {
	f, err := os.Open(file.Path)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"error": "failed to open file",
			"file":  file.Path,
		})
	}
	defer f.Close()

	var meta map[string]interface{}
	content, err := frontmatter.Parse(f, &meta)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"error": "failed to parse frontmatter",
			"file":  file.Path,
		})
	}

	var r r

	// Render the template
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.DefinitionList,
			extension.Footnote,
			extension.GFM,
			emoji.Emoji,
			embed.New(),
			&hashtag.Extender{},
			&anchor.Extender{},
			// &toc.Extender{
			// 	MaxDepth: 4,
			// 	MinDepth: 2,
			// 	Title:    "Contents",
			// 	TitleID:  "toc-header",
			// },
			&wikilink.Extender{
				Resolver: r,
			},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
			html.WithXHTML(),
		),
	)

	buf := bytes.Buffer{}
	if err := md.Convert(content, &buf); err != nil {
		return c.JSON(500, map[string]interface{}{
			"error": "failed to render markdown",
			"file":  file.Path,
		})
	}

	buf.WriteString("<style>")
	buf.WriteString(s.CSS)
	buf.WriteString("</style>")

	return c.HTMLBlob(200, buf.Bytes())
}

//
//
//

// func build(includePaths, excludePaths []string) {
// 	log.Warn("failed to access feature \"build\" - not implemented",
// 		"includePaths", includePaths, "excludePaths", excludePaths)
// }

func scanForKnownFiles(path string) {
	if server == nil {
		return
	}

	// replace with filepath.WalkDir
	filepath.WalkDir(path, func(p string, i os.DirEntry, e error) error {
		if e != nil {
			return e
		}

		// Skip root
		if p == path {
			return nil
		}

		// Recurse into directories
		if i.IsDir() {
			scanForKnownFiles(p)
			return nil
		}

		// If it's a "_root.css" file, load it into the renderer
		if filepath.Base(p) == "_root.css" {
			log.Info("Loading CSS", "path", p)

			f, err := os.Open(p)
			if err != nil {
				return err
			}
			defer f.Close()

			// Read the CSS
			css, err := io.ReadAll(f)
			if err != nil {
				return err
			}

			// Set the CSS
			server.CSS = string(css)
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(p) != ".md" {
			return nil
		}

		log.Info("Loading Markdown File", "path", p)
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		// Parse frontmatter and/or document title
		var meta map[string]interface{}
		_, err = frontmatter.Parse(f, &meta)
		if err != nil {
			return err
		}

		// Determine the title
		var title string
		if t, ok := meta["title"]; ok {
			title = t.(string)
		} else {
			title = strings.TrimSuffix(filepath.Base(p), filepath.Ext(p))
		}

		// Add to known files if it doesn't already exist
		for _, knownFile := range server.KnownFiles {
			if knownFile.Path == p {
				// Update the title
				log.Info("Updating Known File", "path", p, "title", title)
				knownFile.Slug = makeNameCanonical(title)
				knownFile.Title = makeNamePretty(title)
				return nil
			}
		}

		log.Info("Adding Known File", "path", p, "title", title)
		server.KnownFiles = append(server.KnownFiles, &KnownFile{
			Path:  p,
			Slug:  makeNameCanonical(title),
			Title: makeNamePretty(title),
		})
		return nil
	})
}

var watcher *fsnotify.Watcher

func watch(includePaths, excludePaths []string) {
	var err error
	if watcher, err = fsnotify.NewWatcher(); err != nil {
		log.Fatal("failed to create watcher", "error", err)
	}

	for _, includePath := range includePaths {
		if err := filepath.WalkDir(includePath, func(p string, i os.DirEntry, e error) error {
			if e != nil {
				return e
			}

			if !i.IsDir() {
				return nil
			}

			// Ensure p is absolute
			if p, e = filepath.Abs(p); e != nil {
				return e
			}

			// Exclude Paths
			for _, excludePath := range excludePaths {
				if strings.Contains(p, excludePath) {
					log.Warn("Excluding Watch", "path", p, "matched", excludePath)
					return nil
				}
			}

			log.Info("Watching Subdirectory", "path", p, "root", includePath)
			scanForKnownFiles(p)
			return watcher.Add(p)
		}); err != nil {
			log.Fatal("failed to walk path", "error", err, "path", includePath)
		}

		// Ensure the root path is watched
		if err := watcher.Add(includePath); err != nil {
			log.Fatal("failed to watch path", "error", err, "path", includePath)
		}
	}

	// Process Events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Exclude paths
				for _, excludePath := range excludePaths {
					if strings.Contains(event.Name, excludePath) {
						continue
					}
				}

				log.Warn("Event", "Name", event.Name, "Op", event.Op)
				scanForKnownFiles(filepath.Dir(event.Name))

				// build(includePaths, excludePaths)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Error("Watcher Error", "error", err)
			}
		}
	}()
}
