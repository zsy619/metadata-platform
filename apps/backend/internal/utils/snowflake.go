package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake 雪花算法生成器
type Snowflake struct {
	mutex        sync.Mutex
	lastTime     int64 // 上次生成ID的时间戳
	dataCenterID int64 // 数据中心ID
	machineID    int64 // 机器ID
	sequence     int64 // 序列号
}

// NewSnowflake 创建雪花算法生成器
func NewSnowflake(dataCenterID, machineID int64) *Snowflake {
	const (
		machineBits   = 5  // 机器ID位数
		dataCenterBits = 5 // 数据中心ID位数
		machineMax    = -1 ^ (-1 << machineBits)     // 机器ID最大值
		dataCenterMax = -1 ^ (-1 << dataCenterBits)  // 数据中心ID最大值
	)

	// 验证数据中心ID和机器ID
	if dataCenterID > dataCenterMax || dataCenterID < 0 {
		panic("dataCenterID out of range")
	}
	if machineID > machineMax || machineID < 0 {
		panic("machineID out of range")
	}

	return &Snowflake{
		dataCenterID: dataCenterID,
		machineID:    machineID,
		sequence:     0,
		lastTime:     0,
	}
}

// GenerateID 生成唯一ID
func (sf *Snowflake) GenerateID() int64 {
	const (
		sequenceBits  = 12 // 序列号位数
		machineBits   = 5  // 机器ID位数
		dataCenterBits = 5 // 数据中心ID位数
		sequenceMax   = -1 ^ (-1 << sequenceBits)    // 序列号最大值
		machineShift  = sequenceBits                 // 机器ID左移位数
		dataCenterShift = sequenceBits + machineBits // 数据中心ID左移位数
		timestampShift = sequenceBits + machineBits + dataCenterBits // 时间戳左移位数
		epoch         = 1609459200000 // 起始时间戳 (2021-01-01 00:00:00)
	)

	currentTime := time.Now().UnixNano() / 1000000 // 毫秒级时间戳

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	if currentTime < sf.lastTime {
		panic("clock moved backwards")
	}

	if currentTime == sf.lastTime {
		// 同一毫秒内，序列号自增
		sf.sequence = (sf.sequence + 1) & sequenceMax
		if sf.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for currentTime <= sf.lastTime {
				currentTime = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		// 不同毫秒，序列号重置
		sf.sequence = 0
	}

	sf.lastTime = currentTime

	// 生成ID：时间戳 << 时间戳左移位数 | 数据中心ID << 数据中心ID左移位数 | 机器ID << 机器ID左移位数 | 序列号
	return ((currentTime - epoch) << timestampShift) |
		(sf.dataCenterID << dataCenterShift) |
		(sf.machineID << machineShift) |
		sf.sequence
}

// GenerateIDString 生成唯一ID字符串
func (sf *Snowflake) GenerateIDString() string {
	return fmt.Sprintf("%d", sf.GenerateID())
}
