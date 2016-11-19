package room

import (
  "time"

	"github.com/Senior-Design-Kappa/sync-server/models"
)

type RoomState struct {
  LastVideoTime float32
  LastTime      time.Time
  VideoPlaying  bool

  Canvas models.CanvasState
}

func NewRoomState() *RoomState {
  rs := &RoomState{
    LastVideoTime: 0.0,
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
    rs.LastVideoTime = rm.Video.CurrentTime
    rs.VideoPlaying = rm.Video.Playing
    rs.LastTime = time.Now()
  case "SYNC_CANVAS":
    rs.Canvas.UpdateFromCanvasMessage(rm.Message)
  }
}
