package main

import (
	"log"
)

type logScore struct {
	Id        uint64      `json:"id"`
	MonitorId uint32      `json:"monitor_id"`
	ServerId  uint32      `json:"server_id"`
	Ts        uint64      `json:"ts"`
	Score     float64     `json:"score"`
	Step      float64     `json:"step"`
	Offset    interface{} `json:"offset"`
}

type logScores []*logScore

type filterState map[uint32]*logScore

func (ls *logScores) First() *logScore {
	scores := *ls
	return scores[0]
}

func (ls *logScores) Last() *logScore {
	scores := *ls
	return scores[len(scores)-1]
}

func (ls *logScores) filter(wanted int, fn func(*logScore, *filterState)) *logScores {
	if wanted > len(*ls) {
		return ls
	}
	rate := len(*ls) / wanted

	state := make(filterState)
	r := make(logScores, 0)
	for i, l := range *ls {
		fn(l, &state)
		if (i+1)%rate == 0 && l.Ts > 0 {
			// log.Printf("Adding number %v\n", i)

			for _, l := range state {
				r = append(r, l)
			}
			state = make(filterState)

		}
	}
	return &r
}

func (ls *logScores) Sample(t int) *logScores {
	return ls.filter(t, func(l *logScore, st *filterState) {
		(*st)[l.MonitorId] = l
	})
}

func (ls *logScores) WorstOffset(t int) *logScores {

	return ls.filter(t, func(l *logScore, st *filterState) {

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
