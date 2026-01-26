package utils

import (
	"testing"
)

// TestNewSnowflake_ValidParameters 测试使用有效参数创建Snowflake实例
func TestNewSnowflake_ValidParameters(t *testing.T) {
	// 使用有效参数创建实例
	sf := NewSnowflake(1, 1)
	if sf == nil {
		t.Error("NewSnowflake should return a valid instance")
	}
}

// TestNewSnowflake_InvalidMachineID 测试使用无效机器ID创建Snowflake实例
func TestNewSnowflake_InvalidMachineID(t *testing.T) {
	// 机器ID超出范围（0-1023）
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewSnowflake should panic with invalid machineID")
		}
	}()
	sf := NewSnowflake(1024, 1)
	if sf != nil {
		t.Error("NewSnowflake should not return instance with invalid machineID")
	}
}

// TestNewSnowflake_InvalidDatacenterID 测试使用无效数据中心ID创建Snowflake实例
func TestNewSnowflake_InvalidDatacenterID(t *testing.T) {
	// 数据中心ID超出范围（0-31）
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewSnowflake should panic with invalid datacenterID")
		}
	}()
	sf := NewSnowflake(1, 32)
	if sf != nil {
		t.Error("NewSnowflake should not return instance with invalid datacenterID")
	}
}

// TestSnowflake_GenerateID_Unique 测试生成的ID是否唯一
func TestSnowflake_GenerateID_Unique(t *testing.T) {
	sf := NewSnowflake(1, 1)

	// 生成1000个ID，检查是否唯一
	idMap := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := sf.GenerateID()
		if idMap[id] {
			t.Errorf("Generated duplicate ID: %d", id)
		}
		idMap[id] = true
	}
}

// TestSnowflake_GenerateID_Sequence 测试生成的ID是否按时间递增
func TestSnowflake_GenerateID_Sequence(t *testing.T) {
	sf := NewSnowflake(1, 1)

	// 生成多个ID，检查是否递增
	prevID := int64(0)
	for i := 0; i < 100; i++ {
		currID := sf.GenerateID()
		if currID <= prevID {
			t.Errorf("ID should be greater than previous ID, got: %d <= %d", currID, prevID)
		}
		prevID = currID
	}
}

// TestSnowflake_GenerateID_MultipleInstances 测试多个实例生成的ID是否唯一
func TestSnowflake_GenerateID_MultipleInstances(t *testing.T) {
	// 创建两个不同的实例
	sf1 := NewSnowflake(1, 1)
	sf2 := NewSnowflake(2, 1)

	// 生成ID并检查唯一性
	idMap := make(map[int64]bool)
	for i := 0; i < 500; i++ {
		id1 := sf1.GenerateID()
		id2 := sf2.GenerateID()

		if idMap[id1] {
			t.Errorf("Generated duplicate ID from instance 1: %d", id1)
		}
		idMap[id1] = true

		if idMap[id2] {
			t.Errorf("Generated duplicate ID from instance 2: %d", id2)
		}
		idMap[id2] = true
	}
}

// TestSnowflake_GenerateID_Performance 测试生成ID的性能
func TestSnowflake_GenerateID_Performance(t *testing.T) {
	sf := NewSnowflake(1, 1)

	// 生成10000个ID，检查性能
	for i := 0; i < 10000; i++ {
		sf.GenerateID()
	}
}

// BenchmarkSnowflake_GenerateID 性能测试生成的ID是否达到要求 (>100k/s)
func BenchmarkSnowflake_GenerateID(b *testing.B) {
	sf := NewSnowflake(1, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sf.GenerateID()
	}
}
