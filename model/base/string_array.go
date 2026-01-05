// Package base @Author:冯铁城 [17615007230@163.com] 2026-01-05 15:46:45
package base

import (
	"encoding/json"
	"fmt"
)

// StringJSONArray 字符串数组字符串类型
type StringJSONArray []string

// UnmarshalJSON 自定义解析
func (s *StringJSONArray) UnmarshalJSON(data []byte) error {

	//1.如果失败，尝试解析为字符串（如 "[\"a\",\"b\"]"）
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("not a string or string array: %w", err)
	}

	//2.再把字符串内容当作 JSON 解析一次
	var arr []string
	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return fmt.Errorf("inner string is not valid JSON array: %w", err)
	}

	//3.字段赋值
	*s = arr

	//4.返回
	return nil
}
