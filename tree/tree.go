package tree

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jeffail/gabs"
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

func (t Tree) Insert(tag string, child interface{}, empty bool) error {
	gt, err := parseHummusTag(reflect.StructTag(tag))
	if err != nil && strings.Contains(err.Error(), "invalid struct tag") {
		return nil
	} else if err != nil {
		return err
	}

	if gt.omitEmpty && empty {
		return nil
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

			err = childTree.Insert(fmt.Sprintf("hummus:%q", at.childPath), child, empty)
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

func (t Tree) BuildJSON() *gabs.Container {
	jsonObj := gabs.New()

	for path, node := range t.NodeMap {
		if !node.IsArray {
			jsonObj.SetP(node.SingleChild, path)
		} else {
			jsonObj.ArrayP(path)

			for _, child := range node.ArrayChildren {
				if childTree, ok := child.(Tree); ok {
					childJsonObj := childTree.BuildJSON()
					jsonObj.ArrayAppendP(childJsonObj.Data(), path)
				} else {
					jsonObj.ArrayAppendP(child, path)
				}
			}
		}

		replaceHashTags(jsonObj, path)
	}

	return jsonObj
}

func replaceHashTags(obj *gabs.Container, path string) {
	keys := strings.Split(path, ".")
	curKey := keys[0]
	curSubTree := obj.Path(curKey)
	if strings.Contains(curKey, "#") {
		obj.DeleteP(curKey)
		curKey = strings.Replace(curKey, "#", ".", -1)
		existingSubTree := obj.S(curKey)
		subTreeData := existingSubTree.Data()
		subTreeData = mergeObjects(subTreeData, curSubTree.Data())
		obj.Set(subTreeData, curKey)
	}
	if len(keys) != 1 {
		replaceHashTags(curSubTree, strings.Join(keys[1:], "."))
	}
}

func mergeObjects(dst interface{}, src interface{}) interface{} {
	dstMap, dmok := dst.(map[string]interface{})
	srcMap, smok := src.(map[string]interface{})
	dstArr, daok := dst.([]interface{})
	srcArr, saok := src.([]interface{})
	if dmok && smok {
		for k, v := range srcMap {
			dstMap[k] = mergeObjects(dstMap[k], v)
		}
		dst = dstMap
	} else if daok && saok {
		for _, v := range srcArr {
			dstArr = append(dstArr, v)
		}
		dst = dstArr
	} else {
		dst = src
	}

	return dst
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
