package main

import (
	"database/sql"
)

// 1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
func FilterError(err error) error {
	// 不应该wrap这个error抛给上层，因为无数据是正常的。并且抛给上层，上层还要多处理一遍判断是否数据为空的业务逻辑。
	if err != nil && err == sql.ErrNoRows {
		return nil
	}
	return err
}
