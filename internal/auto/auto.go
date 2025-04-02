package auto

import (
	"io/fs"
	"log/slog"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("autoconfigure.enabled", false)
	viper.SetDefault("autoconfigure.root", "/")
	viper.SetDefault("autoconfigure.config.addr", "")
}

type Root interface {
	Exists(path string, mode fs.FileMode) (bool, error)
	Template(path, name string) error
}

type configureFunc func(Root) error

func ConfigureFromEnv() error {
	folder := viper.GetString("autoconfigure.root")

	root, err := newOsRoot(folder)
	if err != nil {
		return err
	}

	return Configure(root)
}

func Configure(root Root) error {
	if enabled := viper.GetBool("autoconfigure.enabled"); !enabled {
		slog.Debug("autoconfigure is disabled")
		return nil
	}

	for _, configure := range []configureFunc{
		configureProfile,
		configureApt,
	} {
		if err := configure(root); err != nil {
			return err
		}
	}

	return nil
}

func configureProfile(root Root) error {
	const (
		profileDir  = "etc/profile.d"
		profileFile = profileDir + "/99-proxyproxy.sh"
	)

	if ok, err := root.Exists(profileDir, fs.ModeDir); err != nil || !ok {
		return err
	}

	slog.Info("configuring profile", slog.String("path", profileFile))
	return root.Template(profileFile, "profile.sh")
}

func configureApt(root Root) error {
	const (
		aptConfDir  = "etc/apt/apt.conf.d"
		aptConfFile = aptConfDir + "/99-proxyproxy.conf"
	)

	if ok, err := root.Exists(aptConfDir, fs.ModeDir); err != nil || !ok {
		return err
	}

	slog.Info("configuring apt", slog.String("path", aptConfFile))
	return root.Template(aptConfFile, "apt.conf")
}
