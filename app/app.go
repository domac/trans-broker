package app

import (
	"errors"
	"github.com/domac/trans-broker/logger"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"net/http"
	"sync"
	"time"
)

//application context
var appCtx *App

type App struct {
	conf       *AppConfig
	lg         *logger.MyLogger
	stopCh     chan bool
	messageCh  chan *Packet
	dataWriter *KafkaOutput

	once sync.Once
}

func (app *App) initKafka() error {
	if len(app.conf.KafkaBrokers) == 0 {
		return errors.New("error occur because kafka brokers was empty")
	}
	var err error
	app.dataWriter, err = NewKafkaProducer(app.conf.KafkaBrokers)
	if err != nil {
		logger.Log().Error(err)
		return err
	}
	return nil

}

func newApp(conf *AppConfig) *App {

	app := new(App)
	app.conf = conf
	app.stopCh = make(chan bool, 1)
	app.messageCh = make(chan *Packet, 2<<10)

	//init kafka writer
	err := app.initKafka()
	if err != nil {
		panic(err)
	}

	if app.conf.BulkMessageFlushInterval > 0 {
		//开启批量写入模式
		go app.batchHandleReceiveData()
	} else {
		go app.handleReceivedData()
	}

	return app
}

func (app *App) handleReceivedData() {
	logger.Log().Info("open handleReceivedData")
	topic := app.conf.Topic
	for {
		select {
		case packet := <-app.messageCh:
			if app.dataWriter != nil {
				app.dataWriter.Write(topic, packet)
			}
		case <-app.stopCh:
			return
		}
	}
}

func (app *App) batchHandleReceiveData() {

	logger.Log().Info("open batchHandleReceiveData")

	topic := app.conf.Topic

	maxWirteBulkSize := app.conf.BulkMessagesCount
	interval := time.Duration(app.conf.BulkMessageFlushInterval) * time.Millisecond

	packets := make([]*Packet, 0, maxWirteBulkSize)

	for {
		select {
		case p := <-app.messageCh:

			packets = append(packets, p)
			chanlen := int(math.Min(float64(len(app.messageCh)), float64(maxWirteBulkSize)))

			for i := 0; i < chanlen; i++ {
				p := <-app.messageCh
				packets = append(packets, p)
			}

			if len(packets) > 0 {
				app.dataWriter.BulkWrite(topic, packets)
				packets = packets[:0]
			}

		case <-app.stopCh:
			return
		default:
			time.Sleep(interval)
		}
	}
}

func (app *App) Put(p *Packet) {
	app.messageCh <- p
}

func (app *App) setLogger(lg *logger.MyLogger) {
	app.lg = lg
}

func (app *App) Log() *zap.SugaredLogger {
	return app.lg.GetSugar()
}

func appRouterEngine() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	lg := logger.GetLogger()
	if lg != nil {
		r.Use(lg.GinzapHandler())
		defer lg.Flush()
	}

	r.Use(gin.Recovery())
	r.Use(VersionMiddleware())

	loadRouter(r)
	return r
}

//运行http服务
func RunHTTPServer(config *AppConfig) (err error) {
	lg := logger.GetLogger()
	//创建应用
	progrm := newApp(config)

	if progrm == nil {
		panic("app run fail")
	}

	addr := config.Bind
	if addr == "" {
		panic("server run without bind addr")
	}

	//设置上下文
	setApplicationContext(progrm)

	progrm.setLogger(lg)

	progrm.Log().Infof("server run : %s", config.Bind)

	err = gracehttp.Serve(&http.Server{
		Addr:    config.Bind,
		Handler: appRouterEngine(),
	})
	return
}

func ShutdownHTTPServer() {
	actx := getApplicationContext()
	if actx == nil {
		panic("app shutdown fail, app context is not exist")
	}
	actx.once.Do(func() {
		close(actx.stopCh)
		close(actx.messageCh)
		actx.Log().Info("app shutdown finish !")
	})
}

func getApplicationContext() *App {
	return appCtx
}

func setApplicationContext(ctx *App) {
	appCtx = ctx
}
