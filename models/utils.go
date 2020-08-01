package models

import (
	"strings"
)

func CombQueryCondition(equalData, likeData map[string]string) string {
	// 拼接where条件
	sb := strings.Builder{}

	for k, v := range equalData {
		if sb.Len() != 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString(k)
		sb.WriteString(" = \"")
		sb.WriteString(v)
		sb.WriteString("\"")
	}

	// 添加机构筛选
	for k, v := range likeData {
		if sb.Len() != 0 {
			sb.WriteString(" AND ")
		}
		sb.WriteString(k)
		sb.WriteString(" LIKE \"")
		sb.WriteString(v)
		sb.WriteString("%\"")
	}

	return sb.String()
}
