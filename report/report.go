package report

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)



type ErrorLogReportConfig struct {
	//每秒钟:n次上报一次
	SecondMaxNumToReport	int64
	MinuteMaxNumToReport	int64
	HourMaxNumToReport	int64
	//每次上报间隔描述
	MinTimeToReport	int64
}

type ErrorLogReport struct {
	second   sync.Map
	minute	 sync.Map
	hour	 sync.Map
	reportTime	sync.Map
	config   *ErrorLogReportConfig
}

func NewErrorLogReport(config *ErrorLogReportConfig) (*ErrorLogReport)  {
	return &ErrorLogReport{
		second: sync.Map{},
		minute:sync.Map{},
		hour:sync.Map{},
		reportTime:sync.Map{},
		config:config,
	}
}

func (e *ErrorLogReport) WriteSecond(msg string) {
	var num int64
	k := fmt.Sprintf("%s_%d",msg,time.Now().Second())
	val,ok := e.second.Load(k)
	if ok {
		num = val.(int64)
		//每秒钟，写入n次，上报一次
		if num % e.config.SecondMaxNumToReport == 0 {
			go e.Report(msg)
		}
	}
	e.second.Store(k, num+1)
}

func (e *ErrorLogReport) WriteMinute(msg string) {
	var num int64
	k := fmt.Sprintf("%s_%d",msg,time.Now().Minute())
	val,ok := e.minute.Load(k)
	if ok {
		num = val.(int64)
		//每分钟，大于n次，上报一次
		if num % e.config.MinuteMaxNumToReport == 0 {
			go e.Report(msg)
		}
	}
	e.minute.Store(k,num + 1)
}

func (e *ErrorLogReport) WriteHour(msg string) {
	var num int64
	k := fmt.Sprintf("%s_%d",msg,time.Now().Hour())
	val,ok := e.hour.Load(k)
	if ok {
		num = val.(int64)
		//每小时，大于n次，上报一次
		if num % e.config.HourMaxNumToReport == 0 {
			go e.Report(msg)
		}
	}
	e.hour.Store(k,num + 1)
}

func (e *ErrorLogReport) Report(msg string) {
	var lastTime int64
	val,ok := e.reportTime.Load(msg)
	if ok {
		lastTime = val.(int64)
		//n秒钟上报一次
		if (time.Now().Unix() - lastTime) < e.config.MinTimeToReport {
			return
		}
	}
	lastTime = time.Now().Unix()
	e.reportTime.Store(msg,lastTime)
	fmt.Println("report",msg,lastTime)
}

func (e *ErrorLogReport) Reset() {
	for {
		second := time.Now().Second()
		e.second.Range(func(key, value interface{}) bool {
			keyStr := key.(string)
			ks,_ := strconv.Atoi( strings.Split(keyStr,"_")[1])
			if ks < second {
				e.second.Delete(key)
				fmt.Println("del second", key,value)
			}
			return true
		})
		minute := time.Now().Minute()
		e.minute.Range(func(key, value interface{}) bool {
			keyStr := key.(string)
			ks,_ := strconv.Atoi( strings.Split(keyStr,"_")[1])
			if ks < minute {
				e.minute.Delete(key)
				fmt.Println("del minute", key,value)
			}
			return true
		})
		hour := time.Now().Hour()
		e.hour.Range(func(key, value interface{}) bool {
			keyStr := key.(string)
			ks,_ := strconv.Atoi( strings.Split(keyStr,"_")[1])
			if ks < hour {
				e.hour.Delete(key)
				fmt.Println("del hour", key,value)
			}
			return true
		})
		time.Sleep(10*time.Second)
	}
}
