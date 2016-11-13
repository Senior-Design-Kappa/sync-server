package room

import (
  "time"

  "github.com/Senior-Design-Kappa/sync-server/models"
)

type RoomState struct {
  models.VideoState
  LastTime      time.Time
  Actions  []interface{}
}

func NewRoomState() *RoomState {
  rs := &RoomState{
    VideoState: models.VideoState {
      CurrentTime: 0.0,
      Playing: false,
      Volume: 1.0,
      Muted: false,
    },
    LastTime: time.Now(),
    Actions: make([]interface{}, 0),
  }
  return rs
}

func (rs *RoomState) UpdateStateFromInboundMessage(m InboundMessage) {
  switch rm := parse(m.RawMessage); rm.MessageType {
  case "SYNC_VIDEO":
    rs.CurrentTime = rm.Video.CurrentTime
    rs.Playing = rm.Video.Playing
    rs.Volume = rm.Video.Volume
    rs.Muted = rm.Video.Muted
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
