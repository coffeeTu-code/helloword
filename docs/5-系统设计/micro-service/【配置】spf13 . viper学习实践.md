# viper 适用于Go应用程序的完整配置解决方案


> https://github.com/spf13/viper



[TOC]



# What is Viper?


Viper是Go应用程序的完整配置解决方案，包括12-Factor应用程序。它旨在在应用程序中工作，并可以处理所有类型的配置需求和格式。它支持：

- 设置默认值
- 从JSON，TOML，YAML，HCL和Java属性配置文件中读取
- 实时观看和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取，并观察变化
- 从命令行标志读取
- 从缓冲区读取
- 设置显式值

Viper可以被认为是所有应用程序配置需求的注册表。


# Why Viper?


在构建现代应用程序时，您无需担心配置文件格式。您想专注于构建出色的软件。Viper就是为此提供帮助的。

Viper为您做了以下事情：

1. 以JSON，TOML，YAML，HCL或Java属性格式查找，加载和解组配置文件。
2. 提供一种机制来为不同的配置选项设置默认值。
3. 提供一种机制来为通过命令行标志指定的选项设置覆盖值。
4. 提供别名系统以轻松重命名参数而不会破坏现有代码。
5. 可以很容易地区分用户提供命令行或配置文件与默认值相同的时间。

Viper使用以下优先顺序。每个项目优先于其下面的项目：

- explicit call to Set
- flag
- env
- config
- key/value store
- default

毒蛇配置键不区分大小写。

# Viper 命令行与 cobra 的结合

通过 Flags() 和 BindPFlags 将viper命令行功能和 cobra命令绑定，
```
serverCmd.Flags().String("conf", "config/service/dsp_retarget.toml", "config file (default is $HOME/configs/service/dsp_retarget.toml)")
_ = viper.BindPFlags(serverCmd.Flags())
```

通过 Get() 将命令行参数提取出，应用到应用程序中。
```
viper.GetString(config.RetargetCfgFile)
```

源代码：
```
func main() {
	var rootCmd = &cobra.Command{Use: "dsp_retarget_server"}

	rootCmd.AddCommand(commands.NewServerCmd())

	_ = rootCmd.Execute()
}


func NewServerCmd() *cobra.Command {
	var serverCmd = &cobra.Command{
		Use:   "serve",
		Short: "server start",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Serve Init PersistentPreRun error:", err)
					rtmetrics.SetMetrics(rtmetrics.Panic, rtmetrics.Labels{FunctionName: "PersistentPreRun"}, 1)
					debug.PrintStack()
				}
			}()

			timeBegin := time.Now()
			log.Println("retarget server start @" + timeBegin.String())
			fmt.Println("config file=", viper.GetString(config.RetargetCfgFile))

			var err error

			//init config
			err = config.RetargetCfg.LoadConfig(viper.GetString(config.RetargetCfgFile))
			if err != nil {
				log.Printf("ParseConfig error: %s\n", err.Error())
				os.Exit(1)
			}
			log.Println("init retarget config success")

			//init logger
			err = rtlogger.Logger.Init(config.RetargetCfg.LogConfig)
			if err != nil {
				log.Println("Load Logger err: ", err)
				os.Exit(1)
			}
			log.Println("init logger success")

			//init metrics
			rtmetrics.InitMetrics()
			rtmetrics.SetMetrics(rtmetrics.GitTag, rtmetrics.Labels{FunctionName: "main", GitTag: reference.VersionMsg.Tag}, 1)
			fmt.Println("git tag=", reference.VersionMsg.Tag)
			log.Println("init metrics success")

			//init juno client
			err = juno_rpc.InitJunoClient()
			if err != nil {
				log.Println("Load Juno err: ", err)
				os.Exit(1)
			}
			log.Println("init juno client success")

			//init polaris lib
			err = polaris_lib.InitPolaris()
			if err != nil {
				log.Println("Load Polaris err: ", err)
				os.Exit(1)
			}
			log.Println("init polaris lib success")

			//init rank lib
			err = rank_lib.InitRank()
			if err != nil {
				log.Println("Load Rank err: ", err)
				os.Exit(1)
			}
			log.Println("init rank lib success")
		},
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("Serve Init Run error:", err)
					rtmetrics.SetMetrics(rtmetrics.Panic, rtmetrics.Labels{FunctionName: "PersistentRun"}, 1)
					debug.PrintStack()
				}
			}()

			signChan := make(chan os.Signal, 1)
			signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
			log.Println(time.Now().Format("2019-06-17 21:50:05"), " Http connector is started.")

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// rpc server
			rpcServer, err := retarget_service.NewRpcServer()
			if err != nil {
				log.Printf("new rpc server error! err[%s]", err.Error())
				os.Exit(-1)
			}
			go func() {
				log.Printf("rpc server start port[%d]", config.RetargetCfg.ServerConfig.RpcPort)
				if err := rpcServer.Start(ctx); err != nil {
					log.Printf("server start error! err[%s]", err.Error())
					return
				}
			}()
			log.Println("rpc server success", time.Now().Format("2019-06-17 21:50:05"))

			//http server
			httpServer := retarget_service.NewHttpServer()
			httpServer.InitRouter()
			go func() {
				log.Printf("http server start port[%d]", config.RetargetCfg.ServerConfig.HttpPort)
				if err := httpServer.Start(ctx); err != nil {
					log.Printf("server start error! err[%s]", err.Error())
				}
			}()
			log.Println("http server success", time.Now().Format("2019-06-17 21:50:05"))

			c := <-signChan
			log.Printf("server stop by signal[%s]", c.String())
			_ = httpServer.Stop(context.Background())
		},
	}
	serverCmd.Flags().String(config.RetargetCfgFile, "config/service/dsp_retarget.toml", "config file (default is $HOME/configs/service/dsp_retarget.toml)")
	_ = viper.BindPFlags(serverCmd.Flags())
	return serverCmd
}

```