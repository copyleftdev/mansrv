package main

import (
	"bytes"
	"container/list"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

type cacheEntry struct {
	key         string
	content     []byte
	listElement *list.Element
}

type LRUCache struct {
	mutex     sync.Mutex
	capacity  int
	cacheMap  map[string]*cacheEntry
	cacheList *list.List
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity:  capacity,
		cacheMap:  make(map[string]*cacheEntry),
		cacheList: list.New(),
	}
}

func (l *LRUCache) Get(key string) ([]byte, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	log.Println("LRU Cache Get: Acquired mutex")

	if entry, found := l.cacheMap[key]; found {
		l.cacheList.MoveToFront(entry.listElement)
		log.Println("LRU Cache Get: Cache hit for key", key)
		return entry.content, true
	}
	log.Println("LRU Cache Get: Cache miss for key", key)
	return nil, false
}

func (l *LRUCache) Set(key string, content []byte) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	log.Println("LRU Cache Set: Acquired mutex")

	if entry, found := l.cacheMap[key]; found {
		l.cacheList.MoveToFront(entry.listElement)
		entry.content = content
		log.Println("LRU Cache Set: Updated key", key)
		return
	}

	if l.cacheList.Len() == l.capacity {
		l.evict()
	}
	entry := &cacheEntry{
		key:     key,
		content: content,
	}
	element := l.cacheList.PushFront(entry)
	entry.listElement = element
	l.cacheMap[key] = entry
	log.Println("LRU Cache Set: Added new key", key)
}

func (l *LRUCache) evict() {
	back := l.cacheList.Back()
	if back != nil {
		entry := back.Value.(*cacheEntry)
		delete(l.cacheMap, entry.key)
		l.cacheList.Remove(back)
		log.Println("LRU Cache Evict: Evicted key", entry.key)
	}
}

var cache = NewLRUCache(50) // Set the cache capacity based on your requirements

func main() {
	router := gin.Default()
	router.Use(gin.Logger()) // Enable logging for HTTP requests

	router.GET("/man/:command", func(c *gin.Context) {
		command := c.Param("command")
		if content, found := cache.Get(command); found {
			c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
			return
		}

		content, err := getManPage(command)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Man page not found"})
			log.Printf("Warning: Man page not found for command %s", command)
			return
		}
		cache.Set(command, content)
		c.Data(http.StatusOK, "text/plain; charset=utf-8", content)
	})

	router.Run(":8887")
}

func getManPage(command string) ([]byte, error) {
	for _, dir := range []string{"/usr/share/man/man1", "/usr/local/man/man1"} {
		filePath := filepath.Join(dir, command+".1.gz")
		if _, err := os.Stat(filePath); err == nil {
			return convertToPlainText(filePath)
		}
	}
	log.Printf("Error: File not found for command %s", command)
	return nil, os.ErrNotExist
}

func convertToPlainText(filePath string) ([]byte, error) {
	cmd := exec.Command("zcat", filePath)             // Decompress the gzipped man page
	groff := exec.Command("groff", "-Tascii", "-man") // Convert troff to plain text
	var out bytes.Buffer
	var stderr bytes.Buffer
	groff.Stdin, _ = cmd.StdoutPipe()
	groff.Stdout = &out
	groff.Stderr = &stderr

	if err := groff.Start(); err != nil {
		log.Printf("Error: Groff start failed for file %s, %v", filePath, err)
		return nil, err
	}

	if err := cmd.Run(); err != nil {
		log.Printf("Error: Zcat failed for file %s, %v, %s", filePath, err, stderr.String())
		return nil, err
	}

	if err := groff.Wait(); err != nil {
		log.Printf("Error: Groff processing failed for file %s, %v, %s", filePath, err, stderr.String())
		return nil, err
	}

	cleanContent := removeFormatting(out.Bytes())
	return cleanContent, nil
}

func removeFormatting(text []byte) []byte {
	re := regexp.MustCompile(".")
	return re.ReplaceAllLiteral(text, nil)
}
