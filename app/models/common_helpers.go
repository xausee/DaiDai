package models

func RemoveDuplicates(slice []ArticleBasicInfo) []ArticleBasicInfo {
	result := []ArticleBasicInfo{}
	seen := map[ArticleBasicInfo]ArticleBasicInfo{}
	for _, val := range slice {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}
