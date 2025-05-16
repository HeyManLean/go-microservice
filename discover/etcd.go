package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func NewEtcdClient(endpoints []string) *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Panicf("init etcd client failed, %v", err)
	}
	return client
}

func getNodes(ctx context.Context, client *clientv3.Client, prefix string) []string {
	/* 获取相关服务的节点列表 */
	var nodes []string
	rsp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return nodes
	}
	for _, kv := range rsp.Kvs {
		nodes = append(nodes, string(kv.Value))
	}
	fmt.Println("get nodes", nodes)
	return nodes
}
func Watch(ctx context.Context, client *clientv3.Client, prefix string) {
	// 监听 prefix 为 order_service 的节点变动事件，并更新当前节点列表
	fmt.Println("get nodes", getNodes(ctx, client, prefix))
	w := client.Watch(ctx, prefix, clientv3.WithPrefix())
	for rsp := range w {
		for _, ev := range rsp.Events {
			fmt.Printf("Event: %s key:%s value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
		fmt.Println("get nodes", getNodes(ctx, client, prefix))
	}
}

func Update(ctx context.Context, client *clientv3.Client) {
	// 注册 order_service 节点, id=1
	_, err := client.Put(ctx, "order_service:1", `{"host":"host1","port":8080}`)
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// 注册 order_service 节点, id=2
	_, err = client.Put(ctx, "order_service:2", `{"host":"host2","port":8080}`)
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}

	// 发现所有 order_service 节点
	resp, err := client.Get(ctx, "order_service", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("get kv, key=%s, value=%s\n", kv.Key, kv.Value)
	}

	// 下线 order_service 节点，id=1
	_, err = client.Delete(ctx, "order_service:1")
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return
	}

	// 重新获取 order_service 节点列表
	resp, err = client.Get(ctx, "order_service", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("get kv, key=%s, value=%s\n", kv.Key, kv.Value)
	}

	// 申请租户，ttl 为 2
	rsp, err := client.Grant(ctx, 2)
	if err != nil {
		fmt.Printf("grant etcd failed, err:%v\n", err)
		return
	}
	leaseId := rsp.ID

	// 注册 order_service 节点，id为1，ttl为2，2秒后自动删除
	_, err = client.Put(ctx, "order_service:1", `{"host":"host1","port":8080}`, clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	krsp, err := client.KeepAliveOnce(ctx, leaseId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ttl:", krsp.TTL)

	time.Sleep(time.Second * 3)
	resp, err = client.Get(ctx, "order_service", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("get etcd failed, err:%v\n", err)
		return
	}
	for _, kv := range resp.Kvs {
		fmt.Printf("get kv, key=%s, value=%s\n", kv.Key, kv.Value)
	}
	time.Sleep(time.Second)
	fmt.Println("OK")
}

func RegisterAndHeartbeat(ctx context.Context, client *clientv3.Client) {
	/* 通过 etcd 实现分布式锁

	相对于 redis，能保证全局一致性，但性能相对较差，涉及磁盘操作
	redis可能存在单点故障，性能较好，内存操作
	*/
	var (
		ttl        int64 = 30 // 超过 30 秒没有续期则认为节点不存在
		retryTimes int64 = 3  // 注册重试次数
		interval   int64 = 5  // 每 5 秒续期
	)

	// 申请租约，可用时间为 30
	lease := clientv3.NewLease(client)
	leaseRsp, err := lease.Grant(ctx, ttl)
	if err != nil {
		log.Panicf("grant etcd error: %v\n", err)
		return
	}

	leaseId := leaseRsp.ID
	fmt.Printf("grant lease succes, ID: %d\n", leaseId)

	servicePath := fmt.Sprintf("/services/%s/%d", "order", leaseId)
	serviceInfo := `{"host":"host1","port":80}`

	kv := clientv3.NewKV(client)

	// 尝试注册到 etcd 中，最大重试3次
	var registerOk bool
	for i := 0; i < int(retryTimes); i++ {
		_, err = kv.Put(ctx, servicePath, serviceInfo, clientv3.WithLease(leaseId))
		if err != nil {
			fmt.Printf("grant etcd error: %v\n", err)
			continue
		}
		registerOk = true
		break
	}
	if !registerOk {
		log.Panicf("register failed")
		return
	}

	// 保持keepalive
	leaveCh, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		log.Panicf("register failed")
		return
	}
	for rsp := range leaveCh {
		if rsp != nil {
			fmt.Printf("heartbeat, ID=%d, TTL=%d\n", rsp.ID, rsp.TTL)
		}
		time.Sleep(time.Duration(interval))
	}
}

func main() {
	client := NewEtcdClient([]string{"127.0.0.1:2379"})
	defer client.Close()
	ctx := context.Background()
	go Watch(ctx, client, "order_service")
	Update(ctx, client)

	go Watch(ctx, client, "/services/order")
	go RegisterAndHeartbeat(ctx, client)
	go RegisterAndHeartbeat(ctx, client)

	for {
		time.Sleep(time.Second * 4)
	}
}
