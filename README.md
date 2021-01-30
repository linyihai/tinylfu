# tinylfu

## go语言实现的LFU缓存，使用哈希表和双链表实现。
主要接口
```
func (lfu *LFUCache) Get(key interface{}) interface{}
func (lfu *LFUCache) Put(key interface{}, value interface{}) 
func (lfu *LFUCache) GetIterator() func() *dbNode
func (lfu *LFUCache) GetAll() []interface{} 
```
提供一种思路实现LFU算法
