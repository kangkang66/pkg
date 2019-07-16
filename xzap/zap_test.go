package xzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestTimeEncoder(t *testing.T) {
	log,err := NewZLog([]string{"stderr","/ssd/kang.log"}, zapcore.DebugLevel)
	if err != nil {
		panic(err)
	}

	log2,err := NewZLog([]string{"stderr","/ssd/kang.log"}, zapcore.DebugLevel)
	if err != nil {
		panic(err)
	}

	for {
		log.Info("hello",zap.String("name","kangkang"))
		log2.Info("hello",zap.String("name","kangkang"))
		time.Sleep(1*time.Second)
	}
}
