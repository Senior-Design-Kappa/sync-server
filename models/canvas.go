package models

import (
  "encoding/json"
  "log"
)

type Point struct {
  X int `json:"x"`
  Y int `json:"y"`
  T1 float32 `json:"t1"`
  T2 float32 `json:"t2"`
  R uint8 `json:"r"`
  G uint8 `json:"g"`
  B uint8 `json:"b"`
}

type CanvasMessage struct {
  MessageType string `json:"type"`

  Points []Point `json:"points"`
}


type CanvasState struct {
  Points []Point `json:"points"`
}

func NewCanvasState() CanvasState {
  return CanvasState {
    Points: []Point{},
  }
}

func (c * CanvasState) UpdateFromCanvasMessage(m string) {
  var cm CanvasMessage
  if err := json.Unmarshal([]byte(m), &cm); err != nil {
    log.Printf("error unmarshaling canvas message: %+v\n", err)
    return
  }
  if cm.MessageType == "POINTS" {
    c.Points = append(c.Points, cm.Points...)
  } else if cm.MessageType == "ERASE" {
    for _, e := range cm.Points {
      for i, p := range c.Points {
        if p.X == e.X && p.Y == e.Y && p.T1 <= e.T1 && p.T1 <= e.T2 {
          c.Points[i].T2 = e.T1 - 0.0000001
        }
      }
    }
  }
}

