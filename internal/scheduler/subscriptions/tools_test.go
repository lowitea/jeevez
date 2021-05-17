package subscriptions

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestGetNowTimeInterval проверяет функцию получения временного интервала
func TestGetNowTimeInterval(t *testing.T) {
	cases := [...]struct {
		Time       string
		expMinTime int
		expMaxTime int
	}{
		{"15:04:05", 54000, 54590},
		{"11:00:00", 39600, 40190},
		{"20:10:00", 72600, 73190},
		{"0:00:00", 0, 590},
		{"23:59:59", 85800, 86390},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("time=%s", c.Time), func(t *testing.T) {
			testTime, _ := time.Parse(time.RFC3339, fmt.Sprintf("2021-06-06T%sZ", c.Time))
			actualMinTime, actualMaxTime := getNowTimeInterval(testTime)
			assert.Equal(t, c.expMinTime, actualMinTime)
			assert.Equal(t, c.expMaxTime, actualMaxTime)
		})
	}
}

// TestSend проверяет функцию отправки сообщений в соответствии с подписками
func TestSend(t *testing.T) {

}
