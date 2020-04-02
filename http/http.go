package http

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/urfave/cli"
	"github.com/valyala/fasthttp"
)

// TestOption 测试参数
type TestOption struct {
	URL      string // 地址
	Co       int    // 协程数量
	Number   int    // 总请求次数
	progress int8   // 当前进度，0-10
}

// RequestResult 请求结果
type RequestResult struct {
	BeginTime  time.Time // 请求开始时间
	EndTime    time.Time // 请求结束时间
	Success    bool      // 是否成功
	ByteLength int       // 数据长度
}

// TestData 测试数据
type TestData struct {
	requestCount    int32 // 请求计数
	TotalRequest    int32 // 总请求次数
	TotalTime       int64 // 请求总耗时
	TotalByteLength int64 // 总的数据长度
	Success         int   // 成功数量
	Failed          int   // 失败数量
}

// Test 测试
func Test(c *cli.Context) error {
	co := c.Int("co")
	testOption := &TestOption{
		URL:    c.String("url"),
		Co:     co,
		Number: c.Int("number"),
	}
	testData := &TestData{}
	// 请求结果通道
	channel := make(chan RequestResult, testOption.Co)
	tasksGroup := sync.WaitGroup{}
	parseGroup := sync.WaitGroup{}
	tasksGroup.Add(testOption.Co)
	parseGroup.Add(1)

	// 接收请求响应结果协程
	go func() {
		defer parseGroup.Done()
		ParseResult(testOption, testData, channel)
	}()

	// 请求工作协程
	for i := 0; i < testOption.Co; i++ {
		go func() {
			defer tasksGroup.Done()
			TestTask(testOption, testData, channel)
		}()
	}

	tasksGroup.Wait()
	close(channel)
	parseGroup.Wait()
	return nil
}

// TestTask 测试任务
func TestTask(testOption *TestOption, testData *TestData, channel chan RequestResult) {
	maxNumber := int32(testOption.Number)
	for {
		if atomic.AddInt32(&testData.requestCount, 1) > maxNumber {
			return
		}

		request := fasthttp.AcquireRequest()
		request.SetRequestURI(testOption.URL)
		response := fasthttp.AcquireResponse()

		result := RequestResult{
			BeginTime: time.Now(),
		}
		err := fasthttp.Do(request, response)
		response.Body()
		result.EndTime = time.Now()

		if err != nil {
			result.Success = false
		} else {
			result.Success = true
			result.ByteLength += len(response.Body())
		}
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)

		channel <- result
	}
}

// ParseResult 处理结果
func ParseResult(testOption *TestOption, testData *TestData, channel chan RequestResult) {
	fmt.Printf("Testing %s\n", testOption.URL)
	beginTime := time.Now()
	for item := range channel {
		if item.Success {
			testData.Success++
		} else {
			testData.Failed++
		}
		testData.TotalRequest++
		testData.TotalTime += item.EndTime.UnixNano() - item.BeginTime.UnixNano()
		testData.TotalByteLength += int64(item.ByteLength)

		if testData.TotalRequest >= int32(float32(testOption.Number)*((float32(testOption.progress)+1)/10)) {
			testOption.progress++
			fmt.Printf("Completed %d requests\n", testData.TotalRequest)
		}
	}
	since := time.Since(beginTime)
	fmt.Println("\nTest result:")
	fmt.Printf("Total requests: %d\n", testData.TotalRequest)
	fmt.Printf("Total time: %s\n", since)
	fmt.Printf("Success requests: %d\n", testData.Success)
	fmt.Printf("Failed requests: %d\n", testData.Failed)
	fmt.Printf("Transfer bytes: %d bytes\n", testData.TotalByteLength)
	fmt.Printf("Time per request: %s\n", time.Duration(float32(testData.TotalTime)/float32(testData.TotalRequest)))
	fmt.Printf("Transfer rate: %f Kb/s\n", float64(testData.TotalByteLength)/since.Seconds()/1024)
	fmt.Printf("Requests per second: %f/s\n", 1.0/since.Seconds()*float64(testData.TotalRequest))
}
