package av

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type RWBaser struct {
	lock               sync.Mutex
	timeout            time.Duration
	PreTime            time.Time
	BaseTimestamp      uint32
	LastVideoTimestamp uint32
	LastAudioTimestamp uint32
}

func NewRWBaser(duration time.Duration) RWBaser {
	return RWBaser{
		timeout: duration,
		PreTime: time.Now(),
	}
}

func (rw *RWBaser) BaseTimeStamp() uint32 {
	return rw.BaseTimestamp
}

func (rw *RWBaser) CalcBaseTimestamp() {
	if rw.LastAudioTimestamp > rw.LastVideoTimestamp {
		rw.BaseTimestamp = rw.LastAudioTimestamp
	} else {
		rw.BaseTimestamp = rw.LastVideoTimestamp
	}
}

func (rw *RWBaser) RecTimeStamp(timestamp, typeID uint32) {
	if typeID == TagVideo {
		rw.LastVideoTimestamp = timestamp
	} else if typeID == TagAudio {
		rw.LastAudioTimestamp = timestamp
	} else {
		log.Warnf("unexpected type id: %d", typeID)
	}
}

func (rw *RWBaser) SetPreTime() {
	rw.lock.Lock()
	defer rw.lock.Unlock()

	rw.PreTime = time.Now()
}

func (rw *RWBaser) Alive() bool {
	rw.lock.Lock()
	defer rw.lock.Unlock()

	alive := !(time.Now().Sub(rw.PreTime) >= rw.timeout)
	log.Debugf("pre_time: %s, timeout: %.0f", rw.PreTime.String(), rw.timeout.Seconds())

	return alive
}
