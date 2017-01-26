package models

import (
  "encoding/json"
  "log"
)

type location struct {
  X int
  Y int
}

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
  Points map[location][]Point
}

func NewCanvasState() CanvasState {
  return CanvasState {
    Points: make(map[location][]Point),
  }
}

func (c CanvasState) MarshalJSON() ([]byte, error) {
  p := make([]Point, 0)
  for _, v := range c.Points {
    p = append(p, v...)
  }
  return json.Marshal(struct {
    Points []Point `json:"points"`
  } {
    Points: p,
  })
}

func (c * CanvasState) UpdateFromCanvasMessage(m string) {
  var cm CanvasMessage
  if err := json.Unmarshal([]byte(m), &cm); err != nil {
    log.Printf("error unmarshaling canvas message: %+v\n", err)
    return
  }
  if cm.MessageType == "POINTS" {
    for _, e := range cm.Points {
      loc := location {
        X: e.X,
        Y: e.Y,
      }
      c.Points[loc] = append(c.Points[loc], e)
    }
  } else if cm.MessageType == "ERASE" {
    for _, e := range cm.Points {
      loc := location {
        X: e.X,
        Y: e.Y,
      }
      for i, p := range c.Points[loc] {
        if p.T1 <= e.T1 && p.T1 <= e.T2 {
          c.Points[loc][i].T2 = e.T1 - 0.0000001
        }
      }
    }
  }
}

