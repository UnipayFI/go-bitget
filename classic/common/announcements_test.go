package common

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestAnnouncements(t *testing.T) {
	c := NewCommonClient(apitest.PublicOptions()...)
	ctx := apitest.Ctx(t)

	// GET /api/v2/public/annoucements (public) — language is required.
	resp, err := c.NewGetAnnouncementsService(AnnouncementLanguageEnUS).Do(ctx)
	if err != nil {
		t.Fatalf("GetAnnouncements: %v", err)
	}
	raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/public/annoucements",
		map[string]string{"language": string(AnnouncementLanguageEnUS)}, false)
	apitest.AssertCovers(t, "Announcements", raw, resp)
}
