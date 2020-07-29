> 官方文档：[https://prometheus.io/docs/prometheus/latest/](https://prometheus.io/docs/prometheus/latest/)
>
> Godoc：[https://godoc.org/github.com/prometheus/client_golang/prometheus](https://godoc.org/github.com/prometheus/client_golang/prometheus)
>
> Process Collector 和Go Collector  介绍[https://povilasv.me/prometheus-go-metrics/](https://povilasv.me/prometheus-go-metrics/)
>
> Prometheus 原理和源码分析 [https://www.infoq.cn/article/Prometheus-theory-source-code](https://www.infoq.cn/article/Prometheus-theory-source-code)  
> 四种数据类型和代码实现讲的较清晰
>
> 聊聊时序数据存储系统的容量管理 [https://mp.weixin.qq.com/s?__biz=MzUyMzA3MTY1NA==&mid=2247485003&idx=1&sn=8981b9a676448c09b21959b6c6d47f86&chksm=f9c37f82ceb4f6948ef1a21e44457ea972d5c14df4562d3070d3cfa3f32a5a0a6bc48185ea8e&scene=21#wechat_redirect](https://mp.weixin.qq.com/s?__biz=MzUyMzA3MTY1NA==&mid=2247485003&idx=1&sn=8981b9a676448c09b21959b6c6d47f86&chksm=f9c37f82ceb4f6948ef1a21e44457ea972d5c14df4562d3070d3cfa3f32a5a0a6bc48185ea8e&scene=21#wechat_redirect)
>
> 百度智能监控系统的过载保护实践 [https://www.infoq.cn/article/zZqZUtJwSjNNtKMkS_6g?utm_source=related_read_bottom&utm_medium=article](https://www.infoq.cn/article/zZqZUtJwSjNNtKMkS_6g?utm_source=related_read_bottom&utm_medium=article)
>


[TOC]

# prometheus架构与运作逻辑

大致使用逻辑是这样： 
1. Prometheus server 定期从静态配置的 targets 或者服务发现的 targets 拉取数据。 
2.  当新拉取的数据大于配置内存缓存区的时候，Prometheus 会将数据持久化到磁盘（如果使用 remote storage 将持久化到云端）。
3. Prometheus 可以配置 rules，然后定时查询数据，当条件触发的时候，会将 alert 推送到配置的 Alertmanager。 
4. Alertmanager 收到警告的时候，可以根据配置，聚合，去重，降噪，最后发送警告。

## Micrometer Api
- Counter:是表示单个单调递增计数器的累积度量，其值只能增加,或在重启时重置为零。
- Gauge:是一个瞬时度量，表示可以任意上下的单个数值。
- Timer/Summary:记录持续时间和响应大小等，是Counter和Gauge的混合示例，有三个指标，总次数和总时长，和时间范围内的max值


## 表达式函数

官网API文档：https://prometheus.io/docs/prometheus/latest/querying/functions/

**increase()**

increase():calculates the increase in the time series in the range vector.指定时间范围内的增量

过去5分钟内新增请求总数：  
increase(loan_total{appName="grafana-design-competition-server"}[5m])


**sum() 、count()**

sum():分组求和。类似select sum() group by。

例如以下result是领取红包编码，表示每个红包的领取成功记录counter求和。  
sum(action_draw_total{appName="red-packet-service",group="lc-yxx",instance="192.168.13.42:8024",type="winning"})by(result)

**注意：** sum()与count()的区别在于，count()方法是select count() group by。

例如以下，表示为每个红包，结果为result = A，数量为1；result = b，数量为1 。。。  
count(action_draw_total{appName="red-packet-service",group="lc-yxx",instance="192.168.13.42:8024",type="winning"})by(result)

组合使用可得，所有存在红包的总数，注意，不是所有红包的领取成功记录数之和：  
sum(count(action_draw_total{appName="red-packet-service",group="lc-yxx",instance="192.168.13.42:8024",type="winning"})by(result))


**rate():**

rate():calculates the per-second average rate of increase of the time series in the range vector. 过去时间范围内，每秒的增量相对该时间范围内总量的占比。

1.QPS统计，根据过去5分钟内每秒增加请求数的占比：  
rate(lz_http_requests_total{job="02_lzmh_microservice_base_service_docker"}[5m])>0

2.QPS统计，根据过去5分钟内每秒增加请求数的占比，并根据handler分组：  
sum(rate(lz_http_requests_total{job="lzmh_microservice_weixin_applet_api"}[5m]))by(handler)>0

3.平均响应时间占比：每秒请求总时长占比除以每秒请求总数占比：  
(rate(lz_http_response_time_milliseconds_sum{job="02_lzmh_microservice_base_service_docker"}[5m])/
rate(lz_http_response_time_milliseconds_count{job="02_lzmh_microservice_base_service_docker"}[5m])
)>0

**offset**

offset：往过去偏移指定时长

获取一天的近4分钟内loan请求总数的增量：  
increase(loan_total{appName="grafana-design-competition-server",type="success"}[5m] offset 1d)

环比计算方法：例如增量请求环比，以30分钟为一个周期，每个周期相对上一个周期请求增量之比
sum(increase(http_requests_total{appName="grafana-design-competition-server"}[10m] offset 10m))by(appName)/
sum(increase(http_requests_total{appName="grafana-design-competition-server"}[10m]))by(appName)


# Retarget Dsp metrics 实践代码：

```

//*************** Metrics const ******************
type Metrics int

const (
	Panic             Metrics = iota //服务Panic次数统计
	GitTag                           //服务发版情况
	ConcurrencyFilter                //服务并发过滤流量统计

	RetargetDspTime //Rt-Dsp 耗时统计
	JunoTime        //Juno 耗时统计
	PolarisTime     //Polaris 耗时统计
	RankTime        //Rank 耗时统计

	CampaignFilter //请求、juno单子过滤、polaris模版召回、rank单子排序 各个模块处理结果的单子数量，数量变化巨大说明模块有隐含的业务问题或提升点
	CampaignBid    //rt-dsp bid 的单子监控
	TemplateBid    //rt-dsp bid 的模版监控

	RankPrice //rt-dsp 算法出价监控
)

//************** SetMetrics *******************
//监控指标维度
type Labels struct {
	GitTag       string
	FunctionName string //触发写metrics函数名字

	Adx         string
	AdType      string
	CountryCode string

	ModelName   string //模块名字
	ModelStatus string //模块处理状态，分为【"can not connected", "timeout", "ok", "fail"】四种基本状态

	RetargetType     string //rt type
	CampaignId       string
	CampaignCostType string
	CampaignPkgName  string
	TemplateId       string
}

var countryTop10 = map[string]bool{}

func (l *Labels) GetCountryTop10() string {
	if countryTop10[l.CountryCode] {
		return l.CountryCode
	}
	return "Others"
}

var metricsAddFunction = map[Metrics]func(labels Labels, value float64){}

func SetMetrics(item Metrics, labels Labels, value float64) {
	addMetrics, ok := metricsAddFunction[item]
	if !ok {
		return
	}
	addMetrics(labels, value)
}

//************** InitMetrics *******************

//InitMetrics: 初始化监控指标
func InitMetrics() {
	metricsConfig := config.RetargetCfg.MetricsConfig
	for _, countryCode := range metricsConfig.TopCountry {
		countryTop10[countryCode] = true
	}

	var (
		region = metricsConfig.Region
		ip     = ""
	)

	//系统健康状态监控
	initSystemMetrics(region, ip)

	//各模块耗时统计
	initModelTimeMetrics(region, ip)

	//单子投放监控
	initCampaignMetrics(region, ip)

	//算法监控
	initRankMetrics(region, ip)
}

func initSystemMetrics(region, ip string) {

	//服务Panic次数统计
	PanicCounter := newCounterMetrics("panic", "Serve Panic Count", []string{"type", "region", "ip"})
	metricsAddFunction[Panic] = func(labels Labels, value float64) {
		PanicCounter.WithLabelValues(labels.FunctionName, region, ip).Inc()
	}

	//服务并发过滤流量统计
	ConcurrencyFilterCounter := newCounterMetrics("concurrency_filter", "Concurrency Filter Count", []string{"region", "ip"})
	metricsAddFunction[ConcurrencyFilter] = func(labels Labels, value float64) {
		ConcurrencyFilterCounter.WithLabelValues(region, ip).Inc()
	}

	//服务发版情况
	TagGauge := newGaugeMetrics("tags", "Retarget Dsp Serve Git Tag Online", []string{"tag", "region", "ip"})
	metricsAddFunction[GitTag] = func(labels Labels, value float64) {
		TagGauge.WithLabelValues(labels.GitTag, region, ip).Set(1)
	}

}

func initModelTimeMetrics(region, ip string) {

	//Rt-Dsp 耗时统计, ModelStatus 是整体的 describe
	RtDspTimeHistogram := newHistogramMetrics("rtdsp_time", "Retarget Dsp Model Time and Status Metrics", []string{"adx", "ad_type", "status", "region"}, []float64{20, 40, 120})
	metricsAddFunction[RetargetDspTime] = func(labels Labels, value float64) {
		RtDspTimeHistogram.WithLabelValues(labels.Adx, labels.AdType, labels.ModelStatus, region).Observe(value)
	}

	//Juno 耗时统计
	JunoTimeHistogram := newHistogramMetrics("juno_time", "Juno Model Time and Status Metrics", []string{"adx", "ad_type", "status", "region"}, []float64{5, 10, 30})
	metricsAddFunction[JunoTime] = func(labels Labels, value float64) {
		JunoTimeHistogram.WithLabelValues(labels.Adx, labels.AdType, labels.ModelStatus, region).Observe(value)
	}

	//Polaris 耗时统计
	PolarisTimeHistogram := newHistogramMetrics("polaris_time", "Polaris Model Time and Status Metrics", []string{"adx", "ad_type", "status", "region"}, []float64{10, 20, 60})
	metricsAddFunction[PolarisTime] = func(labels Labels, value float64) {
		PolarisTimeHistogram.WithLabelValues(labels.Adx, labels.AdType, labels.ModelStatus, region).Observe(value)
	}

	//Rank 耗时统计
	RankTimeHistogram := newHistogramMetrics("rank_time", "Rank Model Time and Status Metrics", []string{"adx", "ad_type", "status", "region"}, []float64{5, 10, 30})
	metricsAddFunction[RankTime] = func(labels Labels, value float64) {
		RankTimeHistogram.WithLabelValues(labels.Adx, labels.AdType, labels.ModelStatus, region).Observe(value)
	}
}

func initCampaignMetrics(region, ip string) {

	//Campaign 过滤情况
	CampaignNumberHistogram := newHistogramMetrics("campaign_filter", "Campaign Filter Number Metrics", []string{"adx", "ad_type", "country", "model_name"}, []float64{100})
	metricsAddFunction[CampaignFilter] = func(labels Labels, value float64) {
		CampaignNumberHistogram.WithLabelValues(labels.Adx, labels.AdType, labels.GetCountryTop10(), labels.ModelName).Observe(value)
	}

	//Campaign 投放情况监控
	CampaignBidCounter := newCounterMetrics("campaign_bid", "Campaign Bid Count", []string{"rt_type", "adx", "campaign_id", "campaign_pkg"})
	metricsAddFunction[CampaignBid] = func(labels Labels, value float64) {
		CampaignBidCounter.WithLabelValues(labels.RetargetType, labels.Adx, labels.CampaignId, labels.CampaignPkgName).Inc()
	}

	//Template 投放情况监控
	TemplateBidCounter := newCounterMetrics("template_bid", "Template Bid Count", []string{"rt_type", "adx", "template_id", "ad_type"})
	metricsAddFunction[TemplateBid] = func(labels Labels, value float64) {
		TemplateBidCounter.WithLabelValues(labels.RetargetType, labels.Adx, labels.TemplateId, labels.AdType).Inc()
	}
}

func initRankMetrics(region, ip string) {

	//算法出价监控
	RankPriceHistogram := newHistogramMetrics("rank_price", "Rank Price Metrics", []string{"adx", "cost_type"}, []float64{0.01, 1, 10})
	metricsAddFunction[RankPrice] = func(labels Labels, value float64) {
		RankPriceHistogram.WithLabelValues(labels.Adx, labels.CampaignCostType).Observe(value)
	}

}

//************* base function ******************

var (
	NameSpace = "Retarget"
	SubSystem = "Dsp"
)

func newCounterMetrics(name, help string, labels []string) *prometheus.CounterVec {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: NameSpace,
		Subsystem: SubSystem,
		Name:      name,
		Help:      help + "< labels=" + strings.Join(labels, ",") + " >",
	}, labels)
	prometheus.MustRegister(counter)
	return counter
}

func newHistogramMetrics(name, help string, labels []string, buckets []float64) *prometheus.HistogramVec {
	histogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: NameSpace,
		Subsystem: SubSystem,
		Name:      name,
		Help:      help + "< labels=" + strings.Join(labels, ",") + " >",
		Buckets:   buckets,
	}, labels)
	prometheus.MustRegister(histogram)
	return histogram
}

func newGaugeMetrics(name, help string, labels []string) *prometheus.GaugeVec {
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: NameSpace,
		Subsystem: SubSystem,
		Name:      name,
		Help:      help + "< labels=" + strings.Join(labels, ",") + " >",
	}, labels)
	prometheus.MustRegister(gauge)
	return gauge
}

```