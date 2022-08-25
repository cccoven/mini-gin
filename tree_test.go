package mini_gin

import (
	"fmt"
	"testing"
)

func TestTree(t *testing.T) {
	pattern := "/p/:lang/doc"

	parts := parsePattern(pattern)
	
	root := &node{}
	root.insert(pattern, parts, 0)
	
	path := "/p/go/doc"
	searchPath := parsePattern(path)
	result := root.search(searchPath, 0)
	fmt.Println(result)
}

func TestTree2(t *testing.T) {
	pattern := "/p/*/doc"
	parts := parsePattern(pattern)
	root := &node{}
	root.insert(pattern, parts, 0)

	path := "/p/go/doc"
	searchPath := parsePattern(path)
	result := root.search(searchPath, 0)
	fmt.Println(result)
}
