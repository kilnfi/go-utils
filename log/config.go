package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Level  logrus.Level `json:"level"`
	Format string       `json:"format"`
	Out    io.Writer    `json:"-"`
}
