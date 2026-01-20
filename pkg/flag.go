package pkg

import (
	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitFlag() {
	viper.SetConfigFile("flag.toml")
	viper.WatchConfig()

	viper.OnConfigChange(func(_ fsnotify.Event) {

		Log.Info("[FLAG]: Change detected")
		cmd.Env.Environment = viper.GetString("env")
		Log, _ = InitLogger(cmd.Env.Environment)
	})
}
