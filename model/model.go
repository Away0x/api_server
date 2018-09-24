package model

import (
	"sync"
	"time"
)

// gorm update_at create_at 写入 0000-00-00 在 Mysql5.7 以上会报错 (需取消设置 mysql sql_mode 中的 NO_ZERO_DATE)
// mysql5.7版本中有了一个STRICT mode（严格模式），而在此模式下默认是不允许设置日期的值为全0值的
/*
mysql -uroot -p
- 查看
select @@GLOBAL.sql_mode
- 设置
set global sql_mode='...之前的 model,NO_ZERO_DATE';
*/
type BaseModel struct {
	// json:"-" 表示 json 时不解析这个字段
	Id uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	// MySQL的DATE/DATATIME类型可以对应Golang的time.Time
	CreatedAt time.Time `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time `gorm:"column:updatedAt" json:"-"`
	// 有 DeletedAt(类型需要是 *time.Time) 即支持 gorm 软删除
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

type UserInfo struct {
	Id        uint64 `json:"id"`
	Username  string `json:"username"`
	SayHello  string `json:"sayHello"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// userlist 查询时会使用 goroutine 所以需锁
type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserInfo // 由于用了协程，所以依赖这个 map(key 为 id) 来进行排序
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}
