package cmdo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scrapli/scrapligo/driver/base"
)

type responseWriter interface {
	WriteResponse(r *base.MultiResponse, name string, d device, appCfg *appCfg) error
}

func (app *appCfg) newResponseWriter(f string) (responseWriter, error) {
	switch f {
	case "file":
		parentDir := "outputs"
		switch app.timestamp {
		case true:
			parentDir = parentDir + "_" + time.Now().Format(time.RFC3339)
		}
		return &fileWriter{
			parentDir,
		}, nil
	case "stdout":
		return &consoleWriter{}, nil
	}

	return nil, nil
}

// consoleWriter writes the scrapli responses to the console
type consoleWriter struct{}

func (w *consoleWriter) WriteResponse(r *base.MultiResponse, name string, d device, appCfg *appCfg) error {
	color.Green("\n**************************\n%s\n**************************\n", name)
	for idx, cmd := range d.SendCommands {
		c := color.New(color.Bold)
		c.Printf("\n-- %s:\n", cmd)
		fmt.Println(r.Responses[idx].Result)
	}
	return nil
}

// fileWriter writes the scrapli responses to the files on disk
type fileWriter struct {
	dir string // output dir name
}

func (w *fileWriter) WriteResponse(r *base.MultiResponse, name string, d device, appCfg *appCfg) error {

	outDir := path.Join(w.dir, name)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	for idx, cmd := range d.SendCommands {
		c := sanitizeCmd(cmd)
		rb := []byte(r.Responses[idx].Result)
		if err := ioutil.WriteFile(path.Join(outDir, c), rb, 0755); err != nil {
			return err
		}
	}
	return nil
}

func sanitizeCmd(s string) string {
	r := strings.NewReplacer(
		"/", "-",
		`\`, "-",
		`"`, ``,
		` `, `-`)

	return r.Replace(s)
}