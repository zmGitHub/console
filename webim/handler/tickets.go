package handler

import (
	"github.com/labstack/echo/v4"

	"bitbucket.org/forfd/custm-chat/webim/handler/adapter"
)

// GET /api/agent/tickets/v2/categories?enterprise_id=5869&browser_id=agent1548679440405
func (s *IMService) GetEntTickets(ctx echo.Context) (err error) {
	return jsonResponse(ctx, &adapter.Tickets{})
}

//GET /api/agent/tickets_v2/1EACfgGhNoogG9YFlTVr2OJt9lK
// {"tickets":[]}
func (s *IMService) GetTicketsV2(ctx echo.Context) (err error) {
	return jsonResponse(ctx, &adapter.Tickets{
		Categories: []*adapter.TicketCategory{},
	})
}
