package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var Root = &cobra.Command{}
var AppName = "app"

func init() {
	log.Println("init")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	Root.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("Config file (default is $HOME/.%s.yaml)", AppName))
	Root.PersistentFlags().String("pid", fmt.Sprintf("/var/run/%s.pid", AppName), fmt.Sprintf("Pid file (default is /var/run/%s.pid)", AppName))
	Root.PersistentFlags().String("loglevel", "info", "Log Level (default: info)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Println("initconfig")
	if cfgFile != "" {
		// enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		path := absPathify("$HOME")
		if _, err := os.Stat(filepath.Join(path, fmt.Sprintf(".%s.yml", AppName))); err != nil {
			_, _ = os.Create(filepath.Join(path, fmt.Sprintf(".%s.yml", AppName)))
		}

		viper.SetConfigType("yaml")
		viper.SetConfigName(fmt.Sprintf(".%s", AppName)) // name of config file (without extension)
		viper.AddConfigPath("$HOME")                     // adding home directory as first search path
	}

	//viper.SetDefault("LOG_LEVEL", "info")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		fmt.Println("")
	}
}

func absPathify(inPath string) string {
	if strings.HasPrefix(inPath, "$HOME") {
		inPath = userHomeDir() + inPath[5:]
	}

	if strings.HasPrefix(inPath, "$") {
		end := strings.Index(inPath, string(os.PathSeparator))
		inPath = os.Getenv(inPath[1:end]) + inPath[end:]
	}

	if filepath.IsAbs(inPath) {
		return filepath.Clean(inPath)
	}

	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}
	return ""
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
