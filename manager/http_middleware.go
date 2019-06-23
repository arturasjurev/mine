package manager

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func (m *Manager) doMetrics(next http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Now().Sub(start)
		logrus.WithFields(logrus.Fields{
			"duration": duration,
		}).Info("metrics")
	}

	return http.HandlerFunc(mw)
}
