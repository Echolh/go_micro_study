package model

// 用户
type UserInfo struct {
	ID   int64  `json:"id"`   // 用户ID
	Name string `json:"name"` // 姓名
	Age  int    `json:"age"`  // 年龄
}
