package room

import (
  "time"

  "github.com/Senior-Design-Kappa/sync-server/models"
)

type RoomState struct {
  LastVideoTime float32
  LastTime      time.Time
  VideoPlaying  bool

  LineSegments  []models.LineSegment
}

func NewRoomState() *RoomState {
  rs := &RoomState{
    LastVideoTime: 0.0,
    LastTime: time.Now(),
    VideoPlaying: false,
    LineSegments: make([]models.LineSegment, 0),
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
    if rm.Message == "DRAW_LINE" {
      rs.LineSegments = append(rs.LineSegments, models.LineSegment{
        PrevX: rm.PrevX,
        PrevY: rm.PrevY,
        CurrX: rm.CurrX,
        CurrY: rm.CurrY,
      })
    }
  }
}
