package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/inconshreveable/mousetrap"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "cvcd",
		Short: "一个小工具, 用于修改VRChat的缓存目录",
		Run: func(cmd *cobra.Command, args []string) {
			startedByExplorer := mousetrap.StartedByExplorer()
			defer func() {
				if err := recover(); err != nil {
					msg := fmt.Sprint(err)
					if _, ok := err.(error); ok && startedByExplorer {
						msg = fmt.Sprintf("出现错误: %v\n", msg)
					}
					fmt.Print(msg)
					if startedByExplorer {
						scanner(true)
					}
					os.Exit(2)
				}
			}()
			if configFile == "" && cfd != "" {
				configFile = cfd
				fmt.Printf("寻找到VRChat配置文件 %s\n确认使用该配置文件吗? (Y/n) ", configFile)
				var yn string
				yn = scanner(false)
				if strings.ToUpper(yn) != "Y" {
					panic(`如果需要指定配置, 请使用参数"--config"指定.`)
				}
			}
			if _, err := os.Stat(configFile); err != nil {
				if err := ioutil.WriteFile(configFile, []byte("{}"), 0600); err != nil {
					panic(err)
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
				newCacheDir = scanner(false)
				if err != nil {
					panic(err)
				}
			}
			newCacheDir = strings.TrimSpace(newCacheDir)
			newCacheDir, err = filepath.Abs(newCacheDir)
			if err != nil {
				panic(err)
			}
			info, err := os.Stat(newCacheDir)
			if err != nil {
				panic(err)
			} else if !info.IsDir() {
				panic(errors.New(`指定的目标不是一个有效的文件夹, 可使用运行参数"--cache"指定缓存目录.`))
			}
			config["cache_directory"] = newCacheDir
			if data, err = json.Marshal(&config); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(configFile, data, 0655); err != nil {
				panic(err)
			}
			if startedByExplorer {
				panic("修改成功, VRChat下一次启动后生效. 按任意键关闭")
			}
		},
	}

	cacheDirectory string
	configFile     string
)

var cfd string

func Execute() {
	cobra.MousetrapHelpText = ""
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func scanner(sp bool) (s string) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s = sc.Text()
		if s != "" || (s == "" && sp) {
			return
		}
	}
	return
}
