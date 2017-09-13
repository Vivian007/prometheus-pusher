package main

//import  导入包
import (
	"github.com/yunlzheng/prometheus-pusher/scrape"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/prometheus/retrieval"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/config"
	"fmt"
	"flag"
	"strings"
)

//var  定义变量名，，，cfg变量名，结构体
var cfg = struct {
	configFile        string
	customLabels      string
	customLabelValues string
}{}//	{}初始化

var (
	labels, values []string
)

//初始化函数，
func init() {
//将prometheus_pusher.yml打入到cfg.configFile变量
	flag.StringVar(
		&cfg.configFile, "config.file", "prometheus_pusher.yml",
		"Prometheus configuration file name.",	
	)
	flag.StringVar(
		&cfg.customLabels, "config.customLabels", "", "custom metrics labels",
	)
	flag.StringVar(
		&cfg.customLabelValues, "config.customLabelValues", "", "custom mertics label values",
	)
}

func main() {
	flag.Parse()		//解析命令行参数
	var (
		//引用 storage.	retrieval. scrape.   包里面的东西	
		sampleAppender = storage.Fanout{}		//数据采集器
		targetManager = retrieval.NewTargetManager(sampleAppender)	//管理target
		jobTargets = scrape.NewJobTargets(targetManager)	//管理job
	)
	//打印日志
	fmt.Println("Loading prometheus config file: " + cfg.configFile)
	fmt.Println("Custom labels: " + cfg.customLabels + "\t Custom label values: " + cfg.customLabelValues)

	if cfg.customLabels == "" {
		labels = []string{}
		values = []string{}
	} else {
		labels = strings.Split(cfg.customLabels, ",")
		values = strings.Split(cfg.customLabelValues, ",")
	}

	var (
		//采集job，labels, values
		scrapeManager = scrape.NewExporterScrape(jobTargets, labels, values)
	)
	
	//conf, err 是LoadFile方法的返回值，，先声明再赋值
	conf, err := config.LoadFile(cfg.configFile)
	//如果出错，就把日志返回给err，然后结束程序
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	
	//不出错的时候，把target拿到
	targetManager.ApplyConfig(conf)
	
	//运行，开始去抓取target
	go targetManager.Run()
	//在函数退出的时候，执行这个方法
	defer targetManager.Stop()

	scrapeManager.AppConfig(conf)

	go scrapeManager.Run()
	defer scrapeManager.Stop()
	
	//gin是golang的web框架
	r := gin.Default()
	//探活,自己
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/targets", func(c *gin.Context) {
		c.JSON(200, jobTargets.Targets())
	})
	r.Run()

}

