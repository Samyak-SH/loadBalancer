package hashring

import (
	"errors"
	"hash/fnv"
	"making-loadbalancer/server"
	"sort"
	"strconv"
	"sync"
)

type HashRing struct {
	nodes            []int
	hashMap          map[int]*server.Server
	virtualNodeCount int
	mu               sync.RWMutex
}

func hashKey(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}

func (hr *HashRing) AddServer(server *server.Server) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	for i := 0; i < hr.virtualNodeCount; i++ {
		virtualNodeKey := server.GetServerURL() + strconv.Itoa(i)
		hash := hashKey(virtualNodeKey)
		hr.nodes = append(hr.nodes, hash)
		hr.hashMap[hash] = server
	}
	sort.Ints(hr.nodes)
}

func (hr *HashRing) RemoveServer(server *server.Server) {
	hr.mu.Lock()
	defer hr.mu.Unlock()
	filtered := hr.nodes[:0]
	for _, hash := range hr.nodes {
		if hr.hashMap[hash] == server {
			delete(hr.hashMap, hash)
		} else {
			filtered = append(filtered, hash)
		}
	}
	hr.nodes = filtered
}

func (hr *HashRing) GetServer(ip string) (*server.Server, error) {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	if len(hr.nodes) == 0 {
		return nil, errors.New("No healthy server")
	}
	ipHash := hashKey(ip)

	for _, nodeHash := range hr.nodes {
		if nodeHash >= ipHash {
			return hr.hashMap[nodeHash], nil
		}
	}
	return hr.hashMap[hr.nodes[0]], nil

}

func InitializeHashRing(virtualNodeCount int, servers []*server.Server) *HashRing {
	hr := &HashRing{
		nodes:            []int{},
		hashMap:          make(map[int]*server.Server),
		virtualNodeCount: virtualNodeCount,
	}
	for _, server := range servers {
		hr.AddServer(server)
	}
	return hr
}
