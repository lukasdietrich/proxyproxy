package auto

import (
	"embed"
	"io/fs"
	"net"
	"os"
	"text/template"

	"github.com/spf13/viper"
)

//go:embed templates/*
var templateFs embed.FS

type osRoot struct {
	*os.Root
	templates *template.Template
}

func newOsRoot(folder string) (Root, error) {
	root, err := os.OpenRoot(folder)
	if err != nil {
		return nil, err
	}

	tmpl, err := template.ParseFS(templateFs, "templates/*")
	if err != nil {
		return nil, err
	}

	return &osRoot{root, tmpl}, nil
}

func (r *osRoot) Exists(path string, mode fs.FileMode) (bool, error) {
	info, err := r.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return info.Mode()&mode == mode, nil
}

func (r *osRoot) Template(path string, name string) error {
	data, err := makeDataFromEnv()
	if err != nil {
		return err
	}

	f, err := r.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	return r.templates.ExecuteTemplate(f, name, data)
}

func makeDataFromEnv() (any, error) {
	addr := viper.GetString("autoconfigure.config.addr")
	if addr == "" {
		addr = viper.GetString("http.addr")
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	if host == "" {
		host = "localhost"
	}

	return map[string]any{
		"scheme": "http",
		"host":   host,
		"port":   port,
	}, nil
}
