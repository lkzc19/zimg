package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"zimg/mapx"
	"zimg/utils"
)

var Zimgrc = mapx.NewSliceMap[string, string]()

var filePath = filepath.Clean(filepath.Join(getHome(), ".zimgrc"))

func getHome() string {
	currentUser, err := user.Current()
	utils.Boom(err)
	return currentUser.HomeDir
}

func Load() {
	file, err := os.Open(filePath)
	// 没有配置文件则创建配置文件
	if os.IsNotExist(err) {
		Zimgrc = template()
		Flush()
		fmt.Println("配置文件已初始化完毕, 请使用[set]命令设置配置文件")
		os.Exit(0)
	}
	utils.Boom(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		i++
		line := scanner.Text()
		line2map := strings.TrimRight(line, " \t\n")
		_line := strings.TrimSpace(line)
		if _line == "" || strings.HasPrefix(_line, "#") {
			Zimgrc.Set(fmt.Sprintf("#%o", i), line2map)
			continue
		}
		kv := strings.SplitN(line2map, "=", 2)
		if len(kv) != 2 {
			continue
		}
		Zimgrc.Set(kv[0], kv[1])
	}
}

func Flush() {
	file, err := os.Create(filePath)
	utils.Boom(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, k := range Zimgrc.Keys() {
		v, _ := Get(k)
		line := fmt.Sprintf("%s=%s", k, v)
		if strings.HasPrefix(k, "#") {
			line = v
		}
		_, err := writer.WriteString(line + "\n")
		utils.Boom(err)
	}

	err = writer.Flush()
	utils.Boom(err)
}

func Get(key string) (string, bool) {
	value, ok := Zimgrc.Get(key)
	return value, ok
}

func Set(key, value string) {
	Zimgrc.Set(key, value)
}

func template() *mapx.SliceMap[string, string] {
	t := mapx.NewSliceMap[string, string]()
	t.Set(Current, Github) // 默认图床源
	t.Set("# 1", "")
	t.Set("# 2", "# github")
	t.Set(GithubOwner, "")
	t.Set(GithubRepo, "")
	t.Set(GithubBucket, "default")
	t.Set(GithubToken, "")
	t.Set("# 3", "")
	t.Set("# 4", "# gitee")
	t.Set(GiteeOwner, "")
	t.Set(GiteeRepo, "")
	t.Set(GiteeBucket, "default")
	t.Set(GiteeToken, "")
	t.Set("# 5", "")
	return t
}

func TestHeader() {
	current, _ := Zimgrc.Get(Current)
	if current == "" {
		fmt.Println("请使用[use]命令选择图床源")
		os.Exit(0)
	}
}

func TestBody() {
	current, _ := Zimgrc.Get(Current)
	if current == Github {
		str, _ := Zimgrc.Get(GithubOwner)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GithubRepo)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GithubBucket)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GithubToken)
		testEmpty(current, str)
	} else if current == Gitee {
		str, _ := Zimgrc.Get(GiteeOwner)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GiteeRepo)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GiteeBucket)
		testEmpty(current, str)
		str, _ = Zimgrc.Get(GiteeToken)
		testEmpty(current, str)
	} else {
		utils.Boom(errors.New("[500] config.config#Test"))
	}
}

func testEmpty(current, str string) {
	if str == "" {
		fmt.Println(fmt.Sprintf("当前正在使用的图床源[%s]未设置, 请使用[set]命令设置", current))
		os.Exit(0)
	}
}

func GetGroup(src string) []string {
	group := make([]string, 0)
	if src == Github {
		owner, _ := Get(GithubOwner)
		group = append(group, fmt.Sprintf("%s=%s", GithubOwner, owner))
		repo, _ := Get(GithubRepo)
		group = append(group, fmt.Sprintf("%s=%s", GithubRepo, repo))
		bucket, _ := Get(GithubBucket)
		group = append(group, fmt.Sprintf("%s=%s", GithubBucket, bucket))
		token, _ := Get(GithubToken)
		group = append(group, fmt.Sprintf("%s=%s", GithubToken, token))
	} else if src == Gitee {
		owner, _ := Get(GiteeOwner)
		group = append(group, fmt.Sprintf("%s=%s", GiteeOwner, owner))
		repo, _ := Get(GiteeRepo)
		group = append(group, fmt.Sprintf("%s=%s", GiteeRepo, repo))
		bucket, _ := Get(GiteeBucket)
		group = append(group, fmt.Sprintf("%s=%s", GiteeBucket, bucket))
		token, _ := Get(GiteeToken)
		group = append(group, fmt.Sprintf("%s=%s", GiteeToken, token))
	}
	return group
}
