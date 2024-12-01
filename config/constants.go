package config

const (
	Current = "current"
	Github  = "github"
	Gitee   = "gitee"
)

var All = []string{Github, Gitee}

const (
	GithubOwner  = "github.owner"
	GithubRepo   = "github.repo"
	GithubBucket = "github.bucket"
	GithubToken  = "github.token"
)

const (
	GiteeOwner  = "gitee.owner"
	GiteeRepo   = "gitee.repo"
	GiteeBucket = "gitee.bucket"
	GiteeToken  = "gitee.token"
)
