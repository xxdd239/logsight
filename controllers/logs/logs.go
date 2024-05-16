package logs

func InitRouter() {
	initOplogRouter()
	initLokiRouter()
	initSyslogRouter()
}
