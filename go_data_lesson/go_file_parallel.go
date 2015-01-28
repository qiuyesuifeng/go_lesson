package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

var Workers = runtime.NumCPU()

// Job
type Job struct {
	data   string
	result chan<- Result
}

func (this *Job) String() string {
	if this == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[Job](%+v)", *this)
}

// Result
type Result struct {
	data string
}

func (this *Result) String() string {
	if this == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[Result](%+v)", *this)
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func AddJobs(jobs chan<- Job, result chan<- Result, input *os.File) {
	br := bufio.NewReader(input)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else {
			realLine := strings.TrimRight(line, "\n")
			// 向job channel中添加任务
			// Notice: 数据不处理，直接发送给下游
			jobs <- Job{realLine, result}
		}
	}

	// 关闭job channel
	close(jobs)
}

func DoJob(job Job) {
	// v=XXXX
	fields := strings.Split(job.data, "=")
	if len(fields) != 2 {
		return
	}

	value := fields[1]
	data := value + "\n"
	job.result <- Result{data}
}

func DoJobs(done chan<- struct{}, jobs <-chan Job) {
	// 在channel中取出任务并计算
	for job := range jobs {
		DoJob(job)
	}

	// 所有工作任务完成后的结束标志
	done <- struct{}{}
}

func AwaitJobDone(done <-chan struct{}, result chan<- Result) {
	for i := 0; i < Workers; i++ {
		<-done
	}

	// 关闭result channel
	close(result)
}

func DoResult(result <-chan Result, output *os.File) {
	for item := range result {
		_, err := output.WriteString(item.data)
		if err != nil {
			fmt.Println("[DoResult][WriteString]Fail\t" + item.data)
		}
	}
}

func main() {
	flag.Usage = Usage
	input := flag.String("i", "", "Input File Name")
	output := flag.String("o", "", "Output File Name")
	flag.Parse()

	if *input == "" {
		Usage()
	}

	if *output == "" {
		Usage()
	}

	fmt.Printf("[Workers]%d\n", Workers)

	var err error
	inputFile, err := os.OpenFile(*input, os.O_RDONLY, 0)
	CheckError(err)
	defer inputFile.Close()

	outputFile, err := os.OpenFile(*output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	CheckError(err)
	defer outputFile.Close()

	jobs := make(chan Job, Workers)
	done := make(chan struct{}, Workers)
	result := make(chan Result, Workers)

	// 将需要并发处理的任务添加到jobs的channel中
	go AddJobs(jobs, result, inputFile)

	// 根据cpu的数量启动对应个数的goroutines从jobs争夺任务进行处理
	for i := 0; i < Workers; i++ {
		go DoJobs(done, jobs)
	}

	// 等待所有worker routiines的完成结果, 并将结果通知主routine
	go AwaitJobDone(done, result)

	DoResult(result, outputFile)
}
