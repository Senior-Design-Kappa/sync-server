package room

import (
  "time"

  "github.com/Senior-Design-Kappa/sync-server/models"
)

type RoomState struct {
  LastVideoTime float32
  LastTime      time.Time
  VideoPlaying  bool

  Actions  []interface{}
}

func NewRoomState() *RoomState {
  rs := &RoomState{
    LastVideoTime: 0.0,
    LastTime: time.Now(),
    VideoPlaying: false,
    Actions: make([]interface{}, 0),
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
      rs.Actions = append(rs.Actions,
        struct {
          models.LineSegment
          Type string `json:"t"`
        } {
          LineSegment: models.LineSegment {
            PrevX: rm.PrevX,
            PrevY: rm.PrevY,
            CurrX: rm.CurrX,
            CurrY: rm.CurrY,
          },
          Type: "DRAW_LINE",
        })
    } else if rm.Message == "ERASE" {
      rs.Actions = append(rs.Actions,
        struct {
          models.ErasePoint
          Type string `json:"t"`
        } {
          ErasePoint: models.ErasePoint {
            X: rm.X,
            Y: rm.Y,
          },
          Type: "ERASE",
        })
    }
  }
}
