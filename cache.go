package assetweb

import (
	"crypto"
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/obnahsgnaw/goutils/security/hsutil"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	cacheControlHeader = "Cache-Control"
	cacheControlValue  = "private, max-age=" // 缓存一个月
	eTagHeader         = "ETag"
	ifNoneMatchHeader  = "If-None-Match"
)

type EtagManager struct {
	fs      *embed.FS
	rootDir string
	etags   map[string]string
}

func newEtagManagerWithFs(fs *embed.FS, rootDir string) *EtagManager {
	s := &EtagManager{fs: fs, rootDir: rootDir, etags: make(map[string]string)}
	return s
}

func newEtagManagerWithDir(rootDir string) *EtagManager {
	s := &EtagManager{rootDir: rootDir, etags: make(map[string]string)}
	return s
}

func (s *EtagManager) Init() (err error) {
	var items []fs.DirEntry
	if s.fs == nil {
		items, err = os.ReadDir(s.rootDir)
	} else {
		items, err = s.fs.ReadDir(s.rootDir)
	}
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() {
			err = s.initDir(s.rootDir, item)
		} else {
			err = s.initFile(s.rootDir, item)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EtagManager) initDir(base string, entry fs.DirEntry) (err error) {
	var items []fs.DirEntry
	name := entry.Name()
	if base != "" {
		name = path.Join(base, name)
	}
	if s.fs == nil {
		items, err = os.ReadDir(name)
	} else {
		items, err = s.fs.ReadDir(name)
	}
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.IsDir() {
			err = s.initDir(name, item)
		} else {
			err = s.initFile(name, item)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EtagManager) initFile(base string, entry fs.DirEntry) (err error) {
	name := entry.Name()
	if strings.HasPrefix(name, ".") {
		return nil
	}
	if base != "" {
		name = path.Join(base, name)
	}
	var f fs.File
	if s.fs == nil {
		f, err = os.Open(name)
	} else {
		f, err = s.fs.Open(name)
	}
	if err != nil {
		return err
	}
	content, err1 := io.ReadAll(f)
	if err1 != nil {
		return err1
	}
	hash, err2 := hsutil.Hash(content, crypto.SHA1)
	if err2 != nil {
		return err2
	}
	s.etags[name] = string(hash)
	return nil
}

func (s *EtagManager) Etag(filename string) string {
	if etag, ok := s.etags[filename]; ok {
		return etag
	}
	return ""
}

// CacheMiddleware 是一个中间件，用于设置缓存控制头
func CacheMiddleware(s *Server, ttl int64) func(c *gin.Context) {
	return func(c *gin.Context) {
		if ttl > 0 {
			// 设置缓存控制头
			c.Header(cacheControlHeader, cacheControlValue+strconv.FormatInt(ttl, 10))

			// 生成并设置 ETag 头
			eTag := s.etagManager.Etag(c.Request.URL.Path)
			c.Header(eTagHeader, eTag)

			// 检查 If-None-Match 头与生成的 ETag 是否匹配，若匹配则返回 304 Not Modified
			if match := c.GetHeader(ifNoneMatchHeader); match != "" {
				if match == eTag {
					c.Status(http.StatusNotModified)
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}
