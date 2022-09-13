package main

import (
	"context"
	"github.com/jolestar/go-commons-pool"
	"fmt"
	"log"
	"net"
	"time"
)

var pCommonPool *pool.ObjectPool

type PoolTest struct {
	Conn *net.UDPConn
}

func (this *PoolTest) SendMsg(data []byte) {
	_, err := this.Conn.Write(data)
	if err != nil {
		log.Printf("write udp error : %+v", err)
		return
	}
	//读取回信
	result := make([]byte, 1024)
	n, remoteAddr, err := this.Conn.ReadFromUDP(result)
	if err != nil {
		log.Printf("read udp server msg error [data:%s] :%+v", string(data), err)
		return
	}
	log.Printf("Recived msg from %s, data:%s \n", remoteAddr, string(result[:n]))
}

func init() {
	// 初始化连接池配置项
	PoolConfig := pool.NewDefaultPoolConfig()
	// 连接池最大容量设置
	PoolConfig.MaxTotal = 3
	WithAbandonedConfig := pool.NewDefaultAbandonedConfig()
	// 注册连接池初始化链接方式
	pCommonPool = pool.NewObjectPoolWithAbandonedConfig(context.Background(), pool.NewPooledObjectFactorySimple(
		func(context.Context) (interface{}, error) {
			return connectUdp()
		}), PoolConfig, WithAbandonedConfig)
}

// 初始化链接类
func connectUdp() (*PoolTest, error) {
	// 创建一个 udp 句柄
	log.Println(">>>>> 创建一个 udp 句柄 ... ")
	// 连接服务器
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9998,
	})

	if err != nil {
		log.Println("Connect to udp server failed,err:", err)
		return nil, err
	}
	log.Printf("<<<<<< new udp connect %+v", conn)
	return &PoolTest{Conn: conn}, nil
}

func main() {

	for i := 0; i < 10; i++ {
		SendMsg(i)
	}

	time.Sleep(2 * time.Second)
}

func SendMsg(num int) {
	var client *PoolTest
	// 从连接池中获取一个实例
	obj, _ := pCommonPool.BorrowObject(context.Background())
	// 转换为对应实体
	if obj != nil {
		client = obj.(*PoolTest)
	}
	// 调用需要的方法
	//if err := client.Conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
	//	log.Printf("SetReadDeadlineerror: %+v", err)
	//	return
	//}
	client.SendMsg([]byte(fmt.Sprintf("send udp data is %d", num)))

	// 交还连接池
	pCommonPool.ReturnObject(context.Background(), client)
}
