package fetch

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type ModelListFile struct {
	Versions []string      `json:"versions"`
	Models   []interface{} `json:"models"`
}

type FetchService struct {
	CachedModelList map[string][]interface{}
}

func NewFetchService() *FetchService {
	// 从 assets/llm-models/ 读取各个文件
	files, err := os.ReadDir("assets/llm-models")
	if err != nil {
		panic(err)
	}

	// 按照文件的名称排序
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// 读取文件
	cachedModelList := make(map[string][]interface{})

	for _, file := range files {
		modelListFile, err := os.ReadFile("assets/llm-models/" + file.Name())
		if err != nil {
			panic(err)
		}

		var modelList ModelListFile

		if err := json.Unmarshal(modelListFile, &modelList); err != nil {
			panic(err)
		}

		for _, version := range modelList.Versions {
			cachedModelList[version] = modelList.Models
		}
	}

	return &FetchService{
		CachedModelList: cachedModelList,
	}
}

func (s *FetchService) FetchModelList(c *gin.Context, version string) ([]interface{}, error) {
	// if _, ok := s.CachedModelList[version]; !ok {
	// 	return s.CachedModelList["default"], nil
	// } else {
	// 	return s.CachedModelList[version], nil
	// }

	for k, v := range s.CachedModelList {
		// 如果version是以key开头的
		if k == version || strings.HasPrefix(version, k+".") {
			return v, nil
		}
	}

	// 如果没有找到对应的版本，就返回默认的模型列表
	return s.CachedModelList["default"], nil
}
