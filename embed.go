package emb

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/dist/*
var distFS embed.FS

// DistDirFS 返回构建好的前端静态资源文件系统
func DistDirFS() http.FileSystem {
	sub, err := fs.Sub(distFS, "web/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}
