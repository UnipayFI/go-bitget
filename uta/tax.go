package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetTaxRecordsService -- GET /api/v3/tax/records (UTA mgt. read)
//
// Returns the unified account's tax records for a product (biz) type within a
// time window that must not exceed 7 days, paginated by cursor.
type GetTaxRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetTaxRecordsService(bizType string, startTime, endTime time.Time) *GetTaxRecordsService {
	return &GetTaxRecordsService{c: c, params: map[string]string{
		"bizType":   bizType,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetMarginType sets the margin type (effective when bizType=MARGIN; defaults to
// crossed).
func (s *GetTaxRecordsService) SetMarginType(marginType MarginMode) *GetTaxRecordsService {
	s.params["marginType"] = string(marginType)
	return s
}

// SetCoin filters to a single coin (returns all coins if unspecified).
func (s *GetTaxRecordsService) SetCoin(coin string) *GetTaxRecordsService {
	s.params["coin"] = coin
	return s
}

// SetLimit sets the page size (default 500, max 500).
func (s *GetTaxRecordsService) SetLimit(limit string) *GetTaxRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetTaxRecordsService) SetCursor(cursor string) *GetTaxRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetTaxRecordsService) Do(ctx context.Context) ([]TaxRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/tax/records", s.params).WithSign()
	resp, err := request.Do[[]TaxRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TaxRecord is a single tax (ledger) record. Amount, fee and balance are
// denominated in the record's coin.
type TaxRecord struct {
	ID      string          `json:"id"`
	Coin    string          `json:"coin"`
	Type    string          `json:"type"`
	Amount  decimal.Decimal `json:"amount"`
	Fee     decimal.Decimal `json:"fee"`
	Balance decimal.Decimal `json:"balance"`
	Ts      time.Time       `json:"ts"`
}
