package args

import (
	"fmt"
)

var (
	Test    = "test"
	Current = "current"
	Use     = "use"
	Get     = "get"
	Set     = "set"
	Help    = "help"
)

func PrintHelp() {
	fmt.Println("Usage of zimg:")
	fmt.Println(" zimg <file>\t上传图片")
	fmt.Println(" zimg test\t检查配置文件, 没有配置文件则创建一份配置文件[~/.zimgrc]")
	fmt.Println(" zimg current\t查看当前使用的图床源")
	fmt.Println(" zimg use\t选择使用的图床源")
	fmt.Println(" zimg get\t选择查看某个图床源配置")
	fmt.Println(" zimg set\t选择设置某个图床源配置")
	fmt.Println(" zimg help\t查看帮助")
}
