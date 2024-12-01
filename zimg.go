package main

import (
	"errors"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strings"
	"zimg/args"
	"zimg/config"
	"zimg/utils"
	"zimg/view"
)

func main() {
	var input []string
	if len(os.Args) > 1 {
		input = os.Args[1:]
	} else {
		args.PrintHelp()
		os.Exit(0)
	}

	if input[0] == args.Help {
		args.PrintHelp()
		os.Exit(0)
	}

	// 加载配置文件
	config.Load()

	if input[0] == args.Test {
		config.TestHeader()
		config.TestBody()
		fmt.Println("OK...")
		os.Exit(0)
	}

	// 校验配置文件
	config.TestHeader()

	if input[0] == args.Use {
		p := tea.NewProgram(view.NewUse())
		_, err := p.Run()
		utils.Boom(err)
		current, _ := config.Get(config.Current)
		fmt.Println("\n图床源已切换为: " + current)
		os.Exit(0)
	} else if input[0] == args.Current {
		current, _ := config.Get(config.Current)
		if current == "" {
			fmt.Println("请使用[use]命令选择图床源")
			os.Exit(0)
		}
		group := config.GetGroup(current)
		builder := strings.Builder{}
		builder.WriteString(fmt.Sprintf("当前图床源是[%s], 配置如下:\n", current))
		for i, e := range group {
			if i != len(group)-1 {
				builder.WriteString(fmt.Sprintf("  %s\n", e))
			} else {
				builder.WriteString(fmt.Sprintf("  %s", e))
			}
		}
		fmt.Println(builder.String())
		os.Exit(0)
	} else if input[0] == args.Get {
		p := tea.NewProgram(view.NewGet())
		m, err := p.Run()
		utils.Boom(err)
		if m, ok := m.(view.Get); ok && m.Result != "" {
			fmt.Println()
			fmt.Println(m.Result)
		}
		os.Exit(0)
	} else if input[0] == args.Set {
		p := tea.NewProgram(view.NewSet())
		_, err := p.Run()
		utils.Boom(err)
		os.Exit(0)
	}

	// 校验配置文件
	config.TestBody()

	current, ok := config.Get(config.Current)
	// 默认图床设置为「github」
	if !ok {
		current = config.Current
	}

	var url string
	if current == config.Github {
		url = uploadGithub(input[0])
	} else if current == config.Gitee { // TODO gitee
		utils.Boom(errors.New(fmt.Sprintf("[%s] not supported", current)))
	} else {
		utils.Boom(errors.New(fmt.Sprintf("[%s] not supported", current)))
	}

	fmt.Println(url)
	os.Exit(0)

}
