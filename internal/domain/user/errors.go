// Package user 用户领域
// 领域错误定义
package user

import "errors"

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")

	// ErrInsufficientCoins 金币不足
	ErrInsufficientCoins = errors.New("insufficient coins")

	// ErrInsufficientDiamonds 钻石不足
	ErrInsufficientDiamonds = errors.New("insufficient diamonds")

	// ErrInvalidAmount 无效金额
	ErrInvalidAmount = errors.New("invalid amount")
)

