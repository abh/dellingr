package scores

import (
	"log"
)

type LogScore struct {
	Id        int64       `json:"id"`
	MonitorId int         `json:"monitor_id"`
	ServerId  int         `json:"server_id"`
	Ts        uint64      `json:"ts"`
	Score     float64     `json:"score"`
	Step      float64     `json:"step"`
	Offset    interface{} `json:"offset"`
}

type LogScores []*LogScore

type filterState map[int]*LogScore

func (ls LogScores) First() *LogScore {
	return ls[0]
}

func (ls LogScores) Last() *LogScore {
	return ls[len(ls)-1]
}

func (ls LogScores) filter(wanted int, fn func(*LogScore, *filterState)) LogScores {
	if wanted > len(ls) {
		return ls
	}

	interval := float64(len(ls)) / float64(wanted)

	log.Printf("has %d, need %d, interval: %.2f\n", len(ls), wanted, interval)

	state := make(filterState)
	r := make(LogScores, 0)

	// the first data point comes after the first "full sample"
	next := 1

	for i, l := range ls {
		fn(l, &state)
		// log.Printf("at number %d, looking for %v >= %v (cur len %v)\n", i, float64(i+1)/interval, next, len(r))

		if float64(i+1)/interval >= float64(next) && l.Ts > 0 {
			next++
			log.Printf("=== Added number %v, len %d, next %v\n", i, len(r), float64(next)*interval)
			for _, l := range state {
				r = append(r, l)
			}
			state = make(filterState)
		}
		if len(r) > 200 {
			// break
		}

	}

	return r
}

func (ls LogScores) Sample(t int) LogScores {
	return ls.filter(t, func(l *LogScore, st *filterState) {
		(*st)[l.MonitorId] = l
	})
}

func (ls LogScores) WorstOffset(t int) LogScores {

	return ls.filter(t, func(l *LogScore, st *filterState) {

		var offset float64

		switch l.Offset.(type) {

		case nil:
			return

		case bool:
			// log.Println("bool...")
			return

		case float64:
			// log.Println("float...")
			offset = l.Offset.(float64)
			if offset < 0 {
				offset = offset * -1
			}

		default:
			log.Printf("type %#v %T %v %#v\n", l, l.Offset, l.Offset, l.Offset)
			panic("unknown type")
		}

		if current, exists := (*st)[l.MonitorId]; exists {
			currentOffset := current.Offset.(float64)
			if currentOffset < 0 {
				currentOffset *= -1
			}
			if offset > currentOffset {
				(*st)[l.MonitorId] = l
			}
		} else {
			switch l.Offset.(type) {
			case float64:
				(*st)[l.MonitorId] = l
			}
		}
	})

}
