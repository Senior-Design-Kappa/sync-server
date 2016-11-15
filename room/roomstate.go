package room

import (
  "time"
)

type point struct {
  x int
  y int
}

type color struct {
  r uint8
  g uint8
  b uint8
  a uint8
}

type RoomState struct {
  LastVideoTime float32
  LastTime      time.Time
  VideoPlaying  bool

  Canvas map[point]color
}

func NewRoomState() *RoomState {
  rs := &RoomState{
    LastVideoTime: 0.0,
    LastTime: time.Now(),
    VideoPlaying: false,
    Canvas: make(map[point]color),
  }
  return rs
}

func (rs *RoomState) UpdateStateFromInboundMessage(m InboundMessage) {
  switch rm := parse(m.RawMessage); rm.MessageType {
  case "SYNC_VIDEO":
    rs.LastVideoTime = rm.Video.CurrentTime
    rs.VideoPlaying = rm.Video.Playing
    rs.LastTime = time.Now()
  case "SYNC_CANVAS":
    if rm.Message == "DRAW_POINTS" {
      for _, pt := range rm.Points {
        coords := point {
          x: pt[0],
          y: pt[1],
        };
        rs.Canvas[coords] = color {
          r: 0,
          g: 0,
          b: 0,
          a: 255,
        };
      }
    } else if rm.Message == "ERASE" {
      for _, pt := range rm.Points {
        coords := point {
          x: pt[0],
          y: pt[1],
        };
        delete(rs.Canvas, coords)
      }
    }
  }
}
