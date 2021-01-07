package week06

import (
	"sync"
	"time"
)

// RollingWindow 可以设定时间的滑动窗口计数器, 由hystrix-go修改而来
type RollingWindow struct {
	buckets    map[int64]*numberBucket
	windowSize int64
	mutex      *sync.RWMutex
}

type numberBucket struct {
	Value float64
}

// NewRollingWindow  初始化滑动窗口计数器，bucketSize单位为秒
func NewRollingWindow(bucketSize int64, bucketNumber int64) *RollingWindow {
	r := &RollingWindow{
		buckets:    make(map[int64]*numberBucket),
		windowSize: bucketSize * bucketNumber,
		mutex:      &sync.RWMutex{},
	}
	return r
}

// getCurrentBucket 获取当前计数桶，如果不存在则创建一个新的
func (r *RollingWindow) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool
	if bucket, ok = r.buckets[now]; !ok {
		bucket = &numberBucket{}
		r.buckets[now] = bucket
	}
	return bucket
}

// removeOldBuckets 删除过期的计数桶
func (r *RollingWindow) removeOldBuckets() {
	now := time.Now().Unix() - r.windowSize
	for timestamp := range r.buckets {
		if timestamp <= now {
			delete(r.buckets, timestamp)
		}
	}
}

// Increment 当前bucket计数1
func (r *RollingWindow) Increment(i float64) {
	if i == 0 {
		return
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

// UpdateMax 更新值最大的那个bucket(应用场景)
func (r *RollingWindow) UpdateMax(n float64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	b := r.getCurrentBucket()
	if n > b.Value {
		b.Value = n
	}
	r.removeOldBuckets()
}

// Sum 获取总和
func (r *RollingWindow) Sum(now time.Time) float64 {
	sum := float64(0)
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for timestamp, bucket := range r.buckets {
		if timestamp >= now.Unix()-r.windowSize {
			sum += bucket.Value
		}
	}
	return sum
}

// Max 获取最大值
func (r *RollingWindow) Max(now time.Time) float64 {
	var max float64
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for timestamp, bucket := range r.buckets {
		if timestamp >= now.Unix()-r.windowSize {
			if bucket.Value > max {
				max = bucket.Value
			}
		}
	}
	return max
}

// Avg 获取平均值
func (r *RollingWindow) Avg(now time.Time) float64 {
	return r.Sum(now) / float64(r.windowSize)
}
