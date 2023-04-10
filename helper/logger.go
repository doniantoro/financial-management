package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	LVL_WARN  = "warning"
	LVL_INFO  = "info"
	LVL_ERROR = "error"
)

type Log struct {
	Type         string
	Event        string
	StatusCode   interface{}
	ResponseTime time.Duration
	Method       string
	Request      interface{} `json:"request"`
	URL          string      `json:"url"`
	Message      string      `json:"message"`
	Response     interface{} `json:"response"`
	ClientIP     string      `json:"clien_ip"`
	UserAgent    string      `json:"user_agent"`
	SupportId    string      `json:"support_id"`
}

func CreateLog(data *Log, file_name string) error {

	logName := fmt.Sprintf(DynamicDir() + "logs/" + file_name + ".log") //generate file log name

	// file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	file, err := os.OpenFile(logName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	log.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      false,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableQuote:    true,
	})
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	message := logrus.Fields{
		"event":         data.Event,
		"status_code":   data.StatusCode,
		"response_time": data.ResponseTime,
		"method":        data.Method,
		"request":       data.Request,
		"message":       data.Message,
		"url":           data.URL,
		"response":      data.Response,
		"client-ip":     data.ClientIP,
		"user-agent":    data.UserAgent,
		"support_id":    data.SupportId,
	}
	if data.Type == LVL_WARN {
		log.WithFields(message).Warn(data.Message)
	}

	if data.Type == LVL_INFO {
		log.WithFields(message).Info(data.Message)
	}

	if data.Type == LVL_ERROR {
		log.WithFields(message).Error(data.Message)
	}

	log.Out = os.Stdout

	return nil
}
