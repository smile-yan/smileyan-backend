package services

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/models"
)

// 搜索服务
type SearchService struct {
	index bleve.Index
	mu    sync.RWMutex
}

var searchService *SearchService

// 初始化搜索服务
func InitSearchService() *SearchService {
	// 创建索引映射
	indexMapping := bleve.NewIndexMapping()

	// 设置文档映射
	documentMapping := bleve.NewDocumentMapping()
	indexMapping.DefaultMapping = documentMapping

	// 标题字段
	titleField := bleve.NewTextFieldMapping()
	titleField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("title", titleField)

	// 内容字段
	contentField := bleve.NewTextFieldMapping()
	contentField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("content", contentField)

	// 摘要字段
	excerptField := bleve.NewTextFieldMapping()
	excerptField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("excerpt", excerptField)

	// 创建内存索引
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		panic("Failed to create search index: " + err.Error())
	}

	searchService = &SearchService{
		index: index,
	}
	return searchService
}

func GetSearchService() *SearchService {
	return searchService
}

// 添加到搜索索引
func AddToSearchIndex(post models.Post) {
	if searchService == nil || searchService.index == nil {
		return
	}

	searchService.mu.Lock()
	defer searchService.mu.Unlock()

	doc := struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Content   string `json:"content"`
		Excerpt   string `json:"excerpt"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		Excerpt:   post.Excerpt,
		CreatedAt: post.CreatedAt.Format("2006-01-02"),
	}

	searchService.index.Index(post.Slug, doc)
}

// 从搜索索引中删除
func RemoveFromSearchIndex(id uint) {
	if searchService == nil || searchService.index == nil {
		return
	}

	searchService.mu.Lock()
	defer searchService.mu.Unlock()

	searchService.index.Delete(fmt.Sprint(id))
}

// 更新搜索索引
func UpdateSearchIndex(post models.Post) {
	RemoveFromSearchIndex(post.ID)
	AddToSearchIndex(post)
}

// 搜索文章
func SearchPosts(keyword string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	if searchService == nil || searchService.index == nil {
		return nil, 0, nil
	}

	searchService.mu.RLock()
	defer searchService.mu.RUnlock()

	query := bleve.NewQueryStringQuery(keyword)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.From = (page - 1) * pageSize
	searchRequest.Size = pageSize

	// 高亮
	highlight := bleve.NewHighlightWithStyle("html")
	searchRequest.Highlight = highlight

	result, err := searchService.index.Search(searchRequest)
	if err != nil {
		return nil, 0, err
	}

	// 从数据库获取完整的文章数据
	var posts []map[string]interface{}
	for _, hit := range result.Hits {
		slug := hit.ID

		// 从数据库查询完整文章
		var post models.Post
		if err := config.DB.Preload("Category").Preload("Tags").Preload("Author").
			Where("slug = ?", slug).First(&post).Error; err != nil {
			continue
		}

		// 处理摘要
		excerpt := post.Excerpt
		if excerpt == "" && post.Content != "" {
			if len(post.Content) > 200 {
				excerpt = post.Content[:200] + "..."
			} else {
				excerpt = post.Content
			}
		}

		postMap := map[string]interface{}{
			"id":          post.ID,
			"slug":        post.Slug,
			"title":       post.Title,
			"excerpt":     excerpt,
			"cover_image": post.CoverImage,
			"view_count":  post.ViewCount,
			"created_at":  post.CreatedAt,
			"status":      post.Status,
		}

		// 添加分类信息
		if post.CategoryID != nil {
			postMap["category"] = map[string]interface{}{
				"id":   post.Category.ID,
				"name": post.Category.Name,
			}
		}

		// 添加标签信息
		if len(post.Tags) > 0 {
			var tags []map[string]interface{}
			for _, tag := range post.Tags {
				tags = append(tags, map[string]interface{}{
					"id":   tag.ID,
					"name": tag.Name,
				})
			}
			postMap["tags"] = tags
		}

		posts = append(posts, postMap)
	}

	return posts, int64(result.Total), nil
}

// 重建搜索索引
func RebuildSearchIndex(posts []models.Post) {
	if searchService == nil {
		return
	}

	searchService.mu.Lock()
	defer searchService.mu.Unlock()

	// 清除旧索引
	searchService.index.Close()

	// 重新创建
	indexMapping := bleve.NewIndexMapping()
	documentMapping := bleve.NewDocumentMapping()
	indexMapping.DefaultMapping = documentMapping

	titleField := bleve.NewTextFieldMapping()
	titleField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("title", titleField)

	contentField := bleve.NewTextFieldMapping()
	contentField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("content", contentField)

	excerptField := bleve.NewTextFieldMapping()
	excerptField.Analyzer = "standard"
	documentMapping.AddFieldMappingsAt("excerpt", excerptField)

	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return
	}
	searchService.index = index

	// 重新索引所有文章
	for _, post := range posts {
		doc := struct {
			ID        uint   `json:"id"`
			Title     string `json:"title"`
			Content   string `json:"content"`
			Excerpt   string `json:"excerpt"`
			CreatedAt string `json:"created_at"`
		}{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			Excerpt:   post.Excerpt,
			CreatedAt: post.CreatedAt.Format("2006-01-02"),
		}
		index.Index(post.Slug, doc)
	}
}