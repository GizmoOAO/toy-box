package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cvcd",
		Short: "一个小工具, 用于修改VRChat的缓存目录",
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					os.Exit(2)
				}
			}()
			var yn string
			var err error
			if configFile == "" && cfd != "" {
				configFile = cfd
				fmt.Printf("寻找到VRChat配置文件 %s\n确认使用该配置文件吗? (Y/n) ", configFile)
				fmt.Scan(&yn)
				if strings.ToUpper(yn) != "Y" {
					panic(`请使用运行参数"--config"指定配置文件.`)
				}
			}
			data, err := ioutil.ReadFile(configFile)
			if err != nil {
				panic(err)
			}
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err != nil {
				panic(err)
			}
			newCacheDir := cacheDirectory
			if cacheDirectory == "" {
				if nowCacheDir, ok := config["cache_directory"]; ok {
					fmt.Println("当前缓存目录:", nowCacheDir)
				}
				fmt.Println("请输入新的缓存目录:")
				fmt.Scan(&newCacheDir)
			}
			newCacheDir, err = filepath.Abs(newCacheDir)
			if err != nil {
				panic(err)
			}
			info, err := os.Stat(newCacheDir)
			if err != nil {
				panic(err)
			} else if !info.IsDir() {
				panic(`指定的目标不是一个有效的文件夹, 请使用运行参数"--cache"指定缓存目录.`)
			}
			config["cache_directory"] = newCacheDir
			if data, err = json.Marshal(&config); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(configFile, data, 0655); err != nil {
				panic(err)
			}
			if cacheDirectory == "" {
				fmt.Println("修改成功, VRChat下一次启动后生效. (Windows可直接关闭本程序窗口)")
				fmt.Scan(&yn)
			}
		},
	}

	cacheDirectory string
	configFile     string
)

var cfd string

func init() {
	cfd = filepath.Join(appdata, "LocalLow\\VRChat\\VRChat\\config.json")
	if _, err := os.Stat(cfd); err != nil {
		cfd = ""
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "VRChat config.json")
	rootCmd.PersistentFlags().StringVarP(&cacheDirectory, "cache", "d", "", "Cache Directory")
}

func Execute() {
	cobra.MousetrapHelpText = ""
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
