package resolver

type subTree = map[any]any

func treeFromJsonData(data map[string]any) subTree {
	tree := make(subTree)
	for k, v := range data {
		tree[k] = v
	}
	return tree
}

func treeToJsonData(tree subTree) map[string]any {
	data := make(map[string]any)
	for k, v := range tree {
		data[k.(string)] = v
	}
	return data
}
