package process

import (
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"path/filepath"
	"sip-parser/pkg/utils/csv_utils"
	"sip-parser/pkg/utils/telnet"
	"sync"
)

// BatchProcessor 用于管理批次任务的并发处理
type BatchProcessor struct {
	// 最大并发数
	concurrencyLimit int
	// 每批次处理的数量
	batchSize int
	// 等待所有任务完成的同步机制
	wg sync.WaitGroup
	// 用于获取处理结果的通道
	resultChan chan int
	// 控制并发的信号量
	sema chan struct{}
}

// NewBatchProcessor 返回一个 BatchProcessor 实例
func NewBatchProcessor(concurrencyLimit, batchSize int) *BatchProcessor {
	return &BatchProcessor{
		concurrencyLimit: concurrencyLimit,
		batchSize:        batchSize,
		resultChan:       make(chan int),
		sema:             make(chan struct{}, concurrencyLimit),
	}
}

// ProcessBatch 处理一批任务
func (bp *BatchProcessor) ProcessBatch(batch []*csv_utils.PcapCsv, startIdx int) {
	// 启动并发任务
	// 启动并发任务
	for idx, pcap := range batch {
		// 计算实际的元素索引
		actualIdx := startIdx + idx
		bp.wg.Add(1)
		bp.sema <- struct{}{} // 限制并发
		go bp.processRow(pcap, actualIdx)
	}
}

// processPcapCsv 处理单个任务
func (bp *BatchProcessor) processRow(pcap *csv_utils.PcapCsv, idx int) {
	defer bp.wg.Done()

	// 模拟处理：5秒后输出 CallId
	//time.Sleep(5 * time.Second)

	err := handleRow(pcap)
	if err != nil {
		log.Println("Skip row:", err)
		idx = -1 //设置为失败
	}

	bp.resultChan <- idx

	// 释放信号量
	<-bp.sema
}

// Wait 等待所有任务完成并关闭结果通道
func (bp *BatchProcessor) Wait() {
	// 等待所有的 goroutine 完成
	bp.wg.Wait()

	// 关闭结果通道，避免其他 goroutine 写入
	close(bp.resultChan)
}

// 输出所有处理的结果
func (bp *BatchProcessor) OutputResults() []int {
	results := []int{}
	for result := range bp.resultChan {
		results = append(results, result)
	}
	return results
}

func NewCalculateSipCost(path string) {
	// 创建客户端实例
	client = telnet.NewTelnetClient("127.0.0.1", "4320")

	csvFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	rows := []*csv_utils.PcapCsv{}

	if err := gocsv.UnmarshalFile(csvFile, &rows); err != nil { // Load clients from file
		panic(err)
	}

	all_count := len(rows)

	// 创建一个 BatchProcessor 实例，设置并发限制为 3，批次大小为 10
	bp := NewBatchProcessor(3, 10)

	// 按批次处理数据
	for i := 0; i < len(rows); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(rows) {
			end = len(rows)
		}
		batch := rows[i:end]

		// 处理当前批次 传入起始索引
		bp.ProcessBatch(batch, i)

		// 等待当前批次的 goroutine 完成并输出结果
		bp.Wait()

		results := bp.OutputResults()
		for _, index := range results {
			if index == -1 {
				log.Println("Skip row:", index)
			} else {
				log.Println("handle->", index, rows[index])
			}
		}

		log.Printf("processing->%d/%d", i, all_count)

		fileName := filepath.Base(path)
		fileName = "res_" + fileName

		csvWriteFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}

		//每操作一次写入一次
		err = gocsv.MarshalFile(&rows, csvWriteFile) // Use this to save the CSV back to the file
		if err != nil {
			panic(err)
		}

		csvWriteFile.Close()

		// 重新初始化 resultChan 为下一个批次清空
		bp.resultChan = make(chan int)
	}

}
