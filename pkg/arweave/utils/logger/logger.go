package logger

import (
	"context"
	"os"

	"github.com/Hubmakerlabs/hoover/pkg/arweave/utils/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ContextLoggerKey = "lg"
)

var logger logrus.FieldLogger

func init() {
	logger = logrus.New()
}

func Init(config *config.Config) (err error) {
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return
	}

	l := logrus.New()
	l.SetLevel(level)
	l.SetOutput(os.Stdout)

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999999999Z07:00",
	}
	l.SetFormatter(formatter)

	logger = l
	return nil
}

func InitWithLogger(l logrus.FieldLogger) {
	logger = l
}

func NewSublogger(tag string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{"module": "warp." + tag})
}

func L() *logrus.Entry {
	return NewSublogger("misc")
}

func LOG(ctx context.Context) *logrus.Entry {
	value := ctx.Value(ContextLoggerKey)
	if value == nil {
		// logrus.Panic("No logger in context")
		return NewSublogger("unclassified")
	}

	log, ok := value.(*logrus.Entry)
	if !ok {
		logrus.Panic("Logger bad type")
	}

	return log
}

func LOGE(c *gin.Context, err error, status int) *logrus.Entry {
	c.Status(status)
	_ = c.Error(err)
	c.Abort()

	entry := LOG(c).WithError(err)

	return entry
}

func SetupLogger(c *gin.Context, fields logrus.Fields) {
	c.Set(ContextLoggerKey,
		NewSublogger("req").WithFields(fields).WithContext(c))
}
