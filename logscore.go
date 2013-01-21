package main

import (
	"log"
)

// {
//   "id": 57221729,
//   "step": 1,
//   "server_id": 417,
//   "ts": 1356997998,
//   "monitor_id": 7,
//   "score": 20,
//   "offset": -0.00528514385223389
// }

// `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
// `monitor_id` int(10) unsigned DEFAULT NULL,
// `server_id` int(10) unsigned NOT NULL,
// `ts` datetime NOT NULL,
// `score` double NOT NULL DEFAULT '0',
// `step` double NOT NULL DEFAULT '0',
// `offset` double DEFAULT NULL,

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

func (ls *logScores) Sample(t int) *logScores {
	if t > len(*ls) {
		return ls
	}
	rate := len(*ls) / t
	i := 0
	r := make(logScores, 0)
	for _, l := range *ls {
		if i%rate == 0 {
			log.Printf("Adding number %v\n", i)
			r = append(r, l)
		}
		i++
	}
	return &r
}

func (ls *logScores) WorstOffset(t int) *logScores {
	if t > len(*ls) {
		return ls
	}
	rate := len(*ls) / t
	i := 0
	r := make(logScores, 0)

	var current *logScore
	var current_offset float64

	for _, l := range *ls {
		i++
		if i%rate == 0 {
			log.Printf("Adding number %v\n", i)
			r = append(r, current)
			current = &logScore{}
			current_offset = 0
		}

		switch l.Offset.(type) {
		case bool:
			// log.Println("bool...")
			continue

		case float64:
			// log.Println("float...")
			offset := l.Offset.(float64)
			if offset < 0 {
				offset = offset * -1
			}

			if offset > current_offset {
				// log.Println("Found worse offset", offset, " > ", current_offset)
				current_offset = offset
				current = l
			} else {
				// log.Println("Found better offset", offset, " < ", current_offset)

			}

		default:
			log.Printf("type %v %#v\n", l.Offset, l.Offset)
			panic("unknown type")
		}

	}
	return &r
}
