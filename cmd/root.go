package cmd

import (
	appInit "first/init"
	"first/init/mysql"
	"first/internal/app/install"
	"first/internal/app/migrate"
	"first/internal/app/router"
	"first/internal/app/service"
	"first/internal/defines"
	"first/internal/model"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "./first",
	Short: "main",
	Long:  `main`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		engine := gin.Default()

		engine.Static("/public", "./public")
		router.Routers(engine)

		tableExist := mysql.NewDB().Migrator().HasTable(model.Version{}.TableName())
		if !tableExist {
			install.Install()
		}
		var currentVerion string
		mysql.NewDB().Model(&model.Version{}).Select("version").Take(&currentVerion)
		if currentVerion != defines.DBVersion {
			fmt.Println("start migrating……")
			migrate.Migrate()
			fmt.Println("migrate successfully")
		}

		//定时任务
		err := service.CronInit()
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("当前时间:", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("AppConfig:%+v", appInit.AppConfig)
		// gin启动

		// 此处初始化一些操作
		// logic.Common{}.InitJob()

		err = engine.Run(":" + appInit.AppConfig.ServicePort)
		if err != nil {
			panic(err.Error())

		}

	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.admin.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("env", "e", true, "aim to enviroment") // -e 绑定 env
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigName(".admin")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
