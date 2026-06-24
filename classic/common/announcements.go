package common

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
)

// AnnouncementType enumerates the announcement categories accepted by the
// annType filter (and echoed back in each announcement's annType field).
type AnnouncementType string

const (
	AnnouncementTypeLatestNews               AnnouncementType = "latest_news"
	AnnouncementTypeCoinListings             AnnouncementType = "coin_listings"
	AnnouncementTypeProductUpdates           AnnouncementType = "product_updates"
	AnnouncementTypeSecurity                 AnnouncementType = "security"
	AnnouncementTypeAPITrading               AnnouncementType = "api_trading"
	AnnouncementTypeMaintenanceSystemUpdates AnnouncementType = "maintenance_system_updates"
	AnnouncementTypeSymbolDelisting          AnnouncementType = "symbol_delisting"
)

// AnnouncementLanguage enumerates the localization codes accepted by the
// required language parameter.
type AnnouncementLanguage string

const (
	AnnouncementLanguageZhCN AnnouncementLanguage = "zh_CN"
	AnnouncementLanguageEnUS AnnouncementLanguage = "en_US"
)

// GetAnnouncementsService -- GET /api/v2/public/annoucements (public)
//
// Returns Bitget platform notices/announcements for a given language,
// optionally filtered by category and creation-time window. Note the path uses
// Bitget's misspelling "annoucements".
type GetAnnouncementsService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetAnnouncementsService(language AnnouncementLanguage) *GetAnnouncementsService {
	return &GetAnnouncementsService{c: c, params: map[string]string{"language": string(language)}}
}

func (s *GetAnnouncementsService) SetAnnType(annType AnnouncementType) *GetAnnouncementsService {
	s.params["annType"] = string(annType)
	return s
}

func (s *GetAnnouncementsService) SetStartTime(t time.Time) *GetAnnouncementsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAnnouncementsService) SetEndTime(t time.Time) *GetAnnouncementsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAnnouncementsService) SetCursor(cursor string) *GetAnnouncementsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetAnnouncementsService) SetLimit(limit int) *GetAnnouncementsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAnnouncementsService) Do(ctx context.Context) ([]Announcement, error) {
	req := request.Get(ctx, s.c, "/api/v2/public/annoucements", s.params)
	resp, err := request.Do[[]Announcement](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Announcement is a single platform notice.
type Announcement struct {
	AnnID      string           `json:"annId"`
	AnnTitle   string           `json:"annTitle"`
	AnnDesc    string           `json:"annDesc"` // brief description (deprecated upstream)
	AnnType    AnnouncementType `json:"annType"`
	AnnSubType string           `json:"annSubType"`
	AnnURL     string           `json:"annUrl"`
	Language   string           `json:"language"`
	CTime      time.Time        `json:"cTime"`
}
