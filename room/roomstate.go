package room

import (
  "time"

	"github.com/Senior-Design-Kappa/sync-server/models"
)

type RoomState struct {
  models.VideoState
  LastTime      time.Time
  VideoPlaying  bool

  Canvas models.CanvasState
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
    VideoPlaying: false,
    Canvas: models.NewCanvasState(),
  }
  return rs
}

func (rs * RoomState) GetVideoTime() float32 {
  videoTime := rs.LastVideoTime
  if rs.VideoPlaying {
    videoTime += float32(time.Now().Sub(rs.LastTime).Seconds())
  }
  return videoTime
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
    rs.Canvas.UpdateFromCanvasMessage(rm.Message)
  }
}
