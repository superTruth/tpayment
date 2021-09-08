package tlog

import (
	"container/list"
	"sync"
	"time"
	"tpayment/pkg/gls"
)

type loggerNode struct {
	tagID         uint64
	log           *Logger
	lifeTimeStamp uint64
}

var (
	loggerMap       = make(map[uint64]*loggerNode)
	loggerList      = list.New()      // 链表，用来处理超时释放
	defaultLifeTime = 5 * time.Minute // 最大超时时间
	lock            sync.RWMutex
)

func SetDefaultLifeTime(d time.Duration) {
	defaultLifeTime = d
}

func SetGoroutineLogger(logger *Logger) {
	lifeTime := uint64(time.Now().Unix()) + uint64(defaultLifeTime/time.Second)

	node := &loggerNode{
		tagID:         gls.GetGoroutineID(),
		log:           logger,
		lifeTimeStamp: lifeTime,
	}

	insertCache(node)
}

func SetLoggerCacheWithTimeout(logger *Logger, d time.Duration) {
	lifeTime := uint64(time.Now().Unix()) + uint64(d/time.Second)

	node := &loggerNode{
		tagID:         gls.GetGoroutineID(),
		log:           logger,
		lifeTimeStamp: lifeTime,
	}

	insertCache(node)
}

func insertCache(node *loggerNode) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()

	nowTime := uint64(time.Now().Unix())

	// 如果已经存在缓存，则删除掉现有的，换成最新的，极低概率会进入这个流程，性能可以忽略不计
	if _, ok := loggerMap[node.tagID]; ok {
		for elm := loggerList.Front(); elm != nil; elm = elm.Next() {
			nodeTmp, _ := elm.Value.(*loggerNode)
			if nodeTmp.tagID == node.tagID { // 查找到目标数据，删除掉
				loggerList.Remove(elm)
				delete(loggerMap, nodeTmp.tagID)
				break
			}
		}
	}

	// 添加到淘汰列表
	elm := loggerList.Back()
	if elm == nil {
		loggerList.PushBack(node)
	} else {
		// 从后往前遍历，遇到比自己生命周期更小的则直接插入到后面
		inserted := false
		for ; elm != nil; elm = elm.Prev() {
			nodeTmp, _ := elm.Value.(*loggerNode)
			if node.lifeTimeStamp >= nodeTmp.lifeTimeStamp {
				loggerList.InsertAfter(node, elm)
				inserted = true
				break
			}
		}
		if !inserted { // 如果还没有插入，则说明需要插入最开头
			loggerList.PushFront(node)
		}
	}

	// 插入到map表
	loggerMap[node.tagID] = node

	// 清理过期列表
	elm = loggerList.Front()
	for elm != nil {
		nodeTmp, _ := elm.Value.(*loggerNode)
		if nodeTmp.lifeTimeStamp >= nowTime {
			// 遇到非过期数据，则不需要继续遍历
			break
		}
		// 过期数据
		loggerList.Remove(elm)
		delete(loggerMap, nodeTmp.tagID)

		elm = loggerList.Front()
	}
}

func GetGoroutineLogger() *Logger {
	lock.RLock()
	defer func() {
		lock.RUnlock()
	}()

	id := gls.GetGoroutineID()
	logNode, ok := loggerMap[id]
	if !ok {
		logger := NewLog(nil)
		return logger
	}

	return logNode.log
}
