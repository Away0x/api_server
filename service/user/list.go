package service

import (
	"fmt"
	"sync"

	"api_server/model"
	"api_server/utils"
)

func ListUser(username string, offset, limit int) ([]*model.UserInfo, uint64, error) {
	infos := make([]*model.UserInfo, 0)
	users, count, err := model.ListUser(username, offset, limit)
	if err != nil {
		return nil, count, err
	}

	// 获得 id 列表，记录顺序
	ids := []uint64{}
	for _, user := range users {
		ids = append(ids, user.Id)
	}

	wg := sync.WaitGroup{}
	userList := model.UserList{
		Lock:  new(sync.Mutex),
		IdMap: make(map[uint64]*model.UserInfo, len(users)),
	}

	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	// Improve query efficiency in parallel
	for _, u := range users {
		wg.Add(1)

		// 对列表的每一项都做操作，如果操作复杂或条数太多，会造成 api 响应延迟
		//    该例中的复杂操作只是
		// 所以这里使用并行查询
		go func(u *model.UserModel) {
			defer wg.Done()

			shortId, err := utils.GenShortId()
			if err != nil {
				errChan <- err // 报错了，发送消息给 errChan
				return
			}

			// 在并发处理中，更新同一个变量为了保证数据一致性，通常需要做锁处理
			userList.Lock.Lock()
			defer userList.Lock.Unlock()
			// 并发时数据被打乱了顺序，所以这里使用 map，id 为 key 以便后续排序
			userList.IdMap[u.Id] = &model.UserInfo{
				Id:        u.Id,
				Username:  u.Username,
				SayHello:  fmt.Sprintf("Hello %s", shortId),
				Password:  u.Password,
				CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: u.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
		}(u)
	}

	go func() {
		wg.Wait() // 上面多个 goroutine 的并行处理完会发送消息给 finished
		close(finished)
	}()

	// 等待消息 (无可用 case 也无 default 会堵塞)
	select {
	case <-finished:
	case err := <-errChan:
		return nil, count, err
	}

	// 将 goroutine 中处理过的乱序数据排序
	for _, id := range ids {
		infos = append(infos, userList.IdMap[id])
	}

	return infos, count, nil
}
