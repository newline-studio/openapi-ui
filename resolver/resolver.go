package resolver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Resolver struct {
	file        string
	pathCleaner func(string) string
	result      map[string]subTree
}

func (r Resolver) traverseAndResolve(tree subTree) (subTree, error) {
	newTree := make(subTree)
	for key, value := range tree {
		switch v := value.(type) {
		case subTree:
			st, err := r.traverseAndResolve(v)
			if err != nil {
				return tree, err
			}
			newTree[key] = st
		case []any:
			newArr := make([]any, len(v))
			for i, val := range v {
				switch aV := val.(type) {
				case subTree:
					st, err := r.traverseAndResolve(aV)
					if err != nil {
						return tree, err
					}
					newArr[i] = st
				default:
					newArr[i] = val
				}
			}
			newTree[key] = newArr
		case string:
			if key != "$ref" || strings.IndexAny(v, "./") != 0 {
				newTree[key] = v
				break
			}
			pTree, err := r.resolveFile(v)
			if err != nil {
				return tree, err
			}
			for pK, pV := range pTree {
				newTree[pK] = pV
			}
		default:
			newTree[key] = value
		}
	}
	return newTree, nil
}

func (r *Resolver) resolveFile(filepath string) (subTree, error) {
	data := make(subTree)
	fileData, err := ioutil.ReadFile(r.pathCleaner(filepath))
	if err != nil {
		return data, errors.New(filepath + " > " + err.Error())
	}
	if strings.HasSuffix(filepath, ".yaml") || strings.HasSuffix(filepath, ".yml") {
		if err = yaml.Unmarshal(fileData, &data); err != nil {
			return data, errors.New(filepath + " > " + err.Error())
		}
	} else if strings.HasSuffix(filepath, ".json") {
		jsonData := make(map[string]any)
		if err = json.Unmarshal(fileData, &jsonData); err != nil {
			return data, errors.New(filepath + " > " + err.Error())
		}
		data = treeFromJsonData(jsonData)
	} else {
		return data, errors.New("unsupported file type " + filepath)
	}

	data, err = r.traverseAndResolve(data)
	if err != nil {
		return data, errors.New(filepath + " > " + err.Error())
	}
	return data, nil
}

func (r *Resolver) Resolve(filepath string) ([]byte, error) {
	tree, err := r.resolveFile(filepath)
	if err != nil {
		return []byte{}, err
	}
	if strings.HasSuffix(filepath, ".yaml") || strings.HasSuffix(filepath, ".yml") {
		return yaml.Marshal(tree)
	} else if strings.HasSuffix(filepath, ".json") {
		return json.Marshal(treeToJsonData(tree))
	}
	return []byte{}, errors.New("unsupported file type " + filepath)
}

func NewResolver(cleaner func(string) string) *Resolver {
	return &Resolver{
		pathCleaner: cleaner,
		result:      make(map[string]subTree),
	}
}
