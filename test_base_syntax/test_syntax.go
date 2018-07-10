package main

/*cgdb :

1. 编译: go build -gcflags "-N -l"  test_syntax.go:关闭优化
2. l main.go:8  l main.main查看源码
3. 断点 b main.main     b test_syntax.go:94

*/
/*
1. 测试基本语法，比如for, range, 数组,等的用法
2. 测试管道,利用读取signal使main阻塞
3. 为何单独的select却不行，会报错 fatal error: all goroutines are asleep - deadlock!
*/
import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Int int

//和c#一样的ToString函数
func (i *Int) ToString() string {
	return strconv.Itoa(int(*i))
}
func print(ch chan int) {
	fmt.Println("Hello world")
	ch <- 1
}

func test_select() {
	in := make(chan string)
	quit := make(chan string)

	//stat := make(chan int)
	var wg sync.WaitGroup

	//子线程 接收阻塞任务 但主线程有绝对控制权,可以随时让子线程退出 需求来源于light中异步接收码值，
	//主线程多次异步发送消息后会导致多个子线程接收数据混乱，所以每次开始发送时需要清除掉之前的子线
	//程 除了使用条件变量外使用for+select+WaitGroup也是不错的选择
	go func() {
		defer func() {
			//stat <- 1
			fmt.Println("thread 1 exit")
			//wg.Add(-1)
			wg.Done()
		}()
		var str string
		var q string
		wg.Add(1)

		for {
			select {
			//会阻塞
			case str = <-in:
				fmt.Println("recv thread 1:", str)
			case q = <-quit:
				//成功退出给主线程发信号
				fmt.Println("quit thread 1:", q)
				return
			}
		}
	}()

	/*
		go func(in chan string, quit chan string) {
			defer func() {
				fmt.Println("thread 2 exit")
				wg.Done()
				//stat <- 1
			}()
			wg.Add(1)
			var str string
			var q string
			for {
				select {
				//会阻塞
				case str = <-in:
					fmt.Println("recv thread 2:", str)
				case q = <-quit:
					//成功退出给主线程发信号
					fmt.Println("quit thread 2:", q)
					return
				}
			}
		}(in, quit)
	*/
	var i Int
	for i = 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		in <- i.ToString()
	}
	quit <- "quit"
	//<-stat
	wg.Wait()
	fmt.Println("father over")
	//fmt.Println("father over")
}

/*
1. 空的interface相当于C语言里面的void
2. 常用来作为函数参数 如:json.Marshal()
*/
func test_interface() {
	var a int
	a = 123
	var v interface{}
	v = a
	fmt.Println("interface :", v)
}

//管道的写入和读取都是阻塞的
func test_chan() {

	chs := make([]chan int, 10) //Panic occurs if no initial length: missing len argument to make([]chan int)
	// 一般的由变量控制的for循环
	for i := 0; i < len(chs); i++ {
		chs[i] = make(chan int)
		go print(chs[i])
	}
	//key 此处为数组的索引
	for key, ch := range chs {

		//<-ch
		fmt.Println("chan key:", key, <-ch)
	}
}
func test_other() {
	//数组元素为string
	str := []string{"welcome", "for", "Chengdu!"}
	//迭代str数组中的各个元素
	for _, v := range str {
		//v为str的备份
		//_为坐标i，这里不需要用到所以用下划线代替
		fmt.Println(v)
	}
	//len cap均指元素的个数
	fmt.Println("长度:", len(str), "容量:", cap(str))

	//字符串转化为[]byte所以用()
	bytes := []byte("11223344")
	fmt.Println(bytes) //打印十进制数[49 49 50 50 51 51 52 52]
	//动态申请空间 此处make申请的为切片数组元素个数为2048,此处默认切片的容量为2048
	//bytes1 := make([]byte, 2048)
	//第二个参数为切片长度,切片容量通过第三个参数指定,此处容量为2060
	bytes2 := make([]byte, 2048, 2060)
	fmt.Println("长度:", len(bytes2), "容量:", cap(bytes2))
	//匿名函数用法
	a := 1
	// 返回参数要么都不写名字，要么全部写上否则会报错：cannot use 0 (type int) as type error in return argument:
	// 返回参数声明用括号括起来
	func() (ret int, err error) {
		//for 3种用法 该示例为死循环用法
		for {
			switch {
			case a == 1:
				fmt.Println("a==1")
				return 0, nil
			}
		}
	}() //加括号为调用

	//map用法
	var student map[int]string     //nil的map
	student = make(map[int]string) //非nil的map

	//student := make(map[int]string) //直接创建
	//初始化+赋值
	studenta := map[int]string{
		2: "aa",
		3: "bb",
	}
	go func() {
		student[1] = "aa" //直接赋值
		student[1] = "bb" //直接赋值
		student[2] = "cc"
		student[3] = "dd"
		fmt.Println("print map1:")
		for key, value := range student {
			fmt.Println(key, ":", value)
		}
		fmt.Println("print map2:")
		for key, value := range studenta {
			fmt.Println(key, ":", value)
		}

		//查找键值 ok为bool类型
		if v, ok := student[2]; ok {
			fmt.Println("found value:", v)
		}
		v, ok := student[22]
		if ok { //正常if用法 if后可跟语句如上所示
			fmt.Println("found key:22 value:", v)
		} else {
			fmt.Println("can not found key:22 from map")
		}
	}()

}

func test_udp() {
	conn, err := net.Dial("udp", "192.168.20.175:7000")
	defer conn.Close()
	if err != nil {
		fmt.Println("udp err")
		os.Exit(1)
	}
	for {
		conn.Write([]byte("00000"))

		time.Sleep(1 * time.Second)
	}
}
func main() {
	fmt.Println("test begin!!")

	//测试udp
	//test_udp()

	//基本语法的测试  ok
	test_other()
	//管道用法 ok
	test_chan()
	test_interface()
	test_select()
	fmt.Println("test over!!")
	//select {} //为何不能使用该句?
	//使用下面的函数阻塞main函数
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	fmt.Println("signal", <-sigChan)

}
