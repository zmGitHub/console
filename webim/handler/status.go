package handler

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/common"
	"bitbucket.org/forfd/custm-chat/webim/common/logging"
	"bitbucket.org/forfd/custm-chat/webim/conf"
	"bitbucket.org/forfd/custm-chat/webim/db"
	"bitbucket.org/forfd/custm-chat/webim/events"
	"bitbucket.org/forfd/custm-chat/webim/imclient"
	"bitbucket.org/forfd/custm-chat/webim/models"
)

const (
	statusOffline       = "offline"
	statusOnline        = "online"
	clientVisitor       = "visitor"
	clientAgent         = "agent"
	defaultSyncInterval = 25 * time.Second
)

// SyncStatusReq
// ClientType == "visitor" --> Values: {"trace_id": "xxx", "ent_id": "xxx"}
// ClientType == "agent"   --> Values: {"agent_id": "xxx"}
type SyncStatusReq struct {
	StatusType string            `json:"status_type"` // offline / online
	ClientType string            `json:"client_type"` // visitor / agent
	Ping       bool              `json:"ping"`
	Values     map[string]string `json:"values"`
}

// PUT /api/v1/status/sync
func (s *IMService) SyncStatus(ctx echo.Context) error {
	req := &SyncStatusReq{}
	if err := ctx.Bind(req); err != nil {
		return invalidParameterResp(ctx, err.Error())
	}

	if req.StatusType == statusOnline {
		if req.ClientType == clientAgent {
			agentID := req.Values["agent_id"]
			if agentID != "" {
				agent, err := models.AgentByID(db.Mysql, agentID)
				if err != nil {
					return dbErrResp(ctx, err.Error())
				}

				if err := models.IncrOnlineAgentCount(db.Mysql, agent.EntID, agentID, req.Ping); err != nil {
					return dbErrResp(ctx, err.Error())
				}

				if !req.Ping {
					go s.sendAgentStatusUpdateEvent(agentID, models.AgentAvailableStatus, true)
				}
			}
		}

		if req.ClientType == clientVisitor {
			entID, traceID := req.Values["ent_id"], req.Values["trace_id"]
			if entID != "" && traceID != "" {
				if err := models.UpdateOnlineVisitor(db.Mysql, traceID, entID); err != nil {
					return dbErrResp(ctx, err.Error())
				}

				if req.Ping {
					if err := models.IncrVisitResidenceTimeSec(db.Mysql, traceID, 25); err != nil {
						log.Logger.Warnf("[SyncStatus] IncrVisitResidenceTimeSec error: %v", err)
					}
					if err := models.IncrVisitorResidenceTimeSec(db.Mysql, traceID, 25); err != nil {
						log.Logger.Warnf("[SyncStatus] IncrVisitorResidenceTimeSec error: %v", err)
					}
				}

				if !req.Ping {
					go s.sendVisitorStatusUpdateEvent(VisitorOnline, entID, traceID, 0)
				}
			}
		}
	}

	if req.StatusType == statusOffline {
		if req.ClientType == clientAgent {
			agentID := req.Values["agent_id"]
			if agentID != "" {
				agent, err := models.AgentByID(db.Mysql, agentID)
				if err != nil {
					return dbErrResp(ctx, err.Error())
				}

				if err := models.DecrOnlineAgentCount(db.Mysql, agent.EntID, agentID); err != nil {
					return dbErrResp(ctx, err.Error())
				}

				// go s.sendAgentStatusUpdateEvent(agentID, models.AgentUnavailableStatus)
			}
		}

		if req.ClientType == clientVisitor {
			entID, traceID := req.Values["ent_id"], req.Values["trace_id"]
			if entID != "" && traceID != "" {
				v, err := models.OnlineVisitorByEntIDTraceID(db.Mysql, entID, traceID)
				if err != nil {
					return dbErrResp(ctx, err.Error())
				}

				//if err := models.DeleteOnlineVisitor(db.Mysql, traceID, entID); err != nil {
				//	return dbErrResp(ctx, err.Error())
				//}

				queue := &models.VisitorQueue{TrackID: traceID}
				if err := queue.Delete(db.Mysql); err == nil {
					go s.sendVisitorQueueingRemove(entID, traceID)
				}

				d := time.Now().UTC().Second() - v.CreatedAt.UTC().Second()
				go s.sendVisitorStatusUpdateEvent(VisitorOffline, entID, traceID, d)
			}
		}
	}

	return jsonResponse(ctx, &SuccessResp{Success: true})
}

func (s *IMService) sendVisitorStatusUpdateEvent(action, entID, traceID string, d int) {
	visits, err := models.VisitsByTraceIDs(db.Mysql, []string{traceID})
	if err != nil {
		return
	}
	if len(visits) > 0 {
		convs, err := models.ConversationsByTraceID(db.Mysql, traceID, 0, 1, false)
		if err != nil {
			return
		}

		if len(convs) > 0 {
			body := &visitorStatusUpdateEventBody{ResidenceTimeSec: d}
			event := VisitorStatusUpdateEvent(action, entID, traceID, visits[0].ID, body)
			eventContent, err := common.Marshal(event)
			if err == nil {
				sendMessageToAgent(s.imCli, convs[0].AgentID, eventContent)
			}
		}
	}
}

func (s *IMService) sendVisitorQueueingRemove(entID, trackID string) {
	event := events.NewQueueingRemove(trackID, entID)
	eventContent, err := common.Marshal(event)
	if err == nil {
		s.sendEventToAllAgents(entID, eventContent)
	}
}

func SyncStatus() {
	defer func() {
		if e := recover(); e != nil {
			log.Logger.Errorf("recover error: %v", e)
		}
	}()

	interval := defaultSyncInterval
	d := conf.IMConf.CentrifugoConf.PingInterval.Duration
	if d >= 25*time.Second {
		interval = d
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)

	syncFn := func() {
		ctx := context.Background()
		channels, err := imclient.CentriClient.ActiveChannels(ctx)
		if err != nil {
			log.Logger.Warnf("active channels error: %v", err)
			return
		}

		var tracks []*models.VisitorTrack
		for _, ch := range channels {
			if strings.HasSuffix(ch, "_message") {
				continue
			}

			ids := strings.Split(ch, "_")
			if len(ids) != 2 {
				continue
			}

			trackID, entID := ids[0], ids[1]
			tracks = append(tracks, &models.VisitorTrack{
				TrackID: trackID,
				EntID:   entID,
			})
		}

		if err := models.BulkUpdateOnlineVisitor(db.Mysql, tracks); err != nil {
			log.Logger.Warnf("BulkUpdateOnlineVisitor error: %v", err)
		}
	}

	log.Logger.Infof("Sync Status Interval: %v", interval)
	ticker := time.NewTicker(interval)
Loop:
	for {
		select {
		case <-ticker.C:
			syncFn()
		case <-quit:
			ticker.Stop()
			log.Logger.Info("Sync Status stopped")
			break Loop
		default:
		}
	}
}
