package tree

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	Path          string
	IsArray       bool
	SingleChild   interface{}
	ArrayChildren []interface{}
}

type Tree struct {
	NodeMap map[string]Node
}

type hummusTag struct {
	tagName   string
	omitEmpty bool
}

type arrayTag struct {
	arrayPath  string
	arrayIndex int
	childPath  string
}

func NewTree() Tree {
	return Tree{
		NodeMap: make(map[string]Node),
	}
}

func (t Tree) Insert(tag string, child interface{}) error {
	gt, err := parseHummusTag(reflect.StructTag(tag))
	if err != nil {
		return err
	}

	childToInsert := child

	if at, yes := parseArrayTag(gt.tagName); yes {
		if _, exists := t.NodeMap[at.arrayPath]; !exists {
			t.NodeMap[at.arrayPath] = Node{
				Path:          at.arrayPath,
				IsArray:       true,
				ArrayChildren: []interface{}{},
			}
		}

		node := t.NodeMap[at.arrayPath]
		if len(node.ArrayChildren) <= at.arrayIndex {
			arrayChildrenCopy := make([]interface{}, at.arrayIndex+1)
			for i, c := range node.ArrayChildren {
				arrayChildrenCopy[i] = c
			}
			node.ArrayChildren = arrayChildrenCopy
		}

		if at.childPath != "" {
			childToInsert = node.ArrayChildren[at.arrayIndex]
			if childToInsert == nil {
				childToInsert = NewTree()
			}

			childTree, ok := childToInsert.(Tree)
			if !ok {
				return errors.New("fatal error: existing subchild is not a tree")
			}

			err = childTree.Insert(fmt.Sprintf("hummus:%q", at.childPath), child)
			if err != nil {
				return err
			}
		}

		node.ArrayChildren[at.arrayIndex] = childToInsert
		t.NodeMap[at.arrayPath] = node
	} else {
		t.NodeMap[gt.tagName] = Node{
			Path:        gt.tagName,
			IsArray:     false,
			SingleChild: child,
		}
	}

	return nil
}

func parseHummusTag(tag reflect.StructTag) (hummusTag, error) {
	hummusTagString := tag.Get("hummus")
	if hummusTagString == "" {
		return hummusTag{}, fmt.Errorf("error: invalid struct tag %s", tag)
	}

	tagFields := strings.Split(hummusTagString, ",")
	var omitEmpty bool

	if len(tagFields) > 2 {
		return hummusTag{}, errors.New("error: invalid number of struct tag fields")
	}

	if len(tagFields) == 2 && tagFields[1] == "omitempty" {
		omitEmpty = true
	}

	return hummusTag{
		tagName:   tagFields[0],
		omitEmpty: omitEmpty,
	}, nil
}

func parseArrayTag(tag string) (arrayTag, bool) {
	if strings.Contains(tag, "[") {
		//regex needs to be non-greedy in order to catch the parent array path first
		//(in case of arrays inside arrays)
		arrayRegex := regexp.MustCompile("(.*?)\\[(\\d+)\\]\\.*(.*)")

		matches := arrayRegex.FindStringSubmatch(tag)
		if len(matches) != 4 {
			return arrayTag{}, false
		}

		arrayIndex, err := strconv.Atoi(matches[2])
		if err != nil {
			return arrayTag{}, false
		}

		return arrayTag{
			arrayPath:  matches[1],
			arrayIndex: arrayIndex,
			childPath:  matches[3],
		}, true
	}
	return arrayTag{}, false
}
