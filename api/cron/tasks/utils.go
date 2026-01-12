package tasks

import "time"

var NowFunc = func() time.Time { return time.Now().UTC() }
