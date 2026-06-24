package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetDepositAddressService -- GET /api/v3/account/deposit-address (UTA mgt. read)
//
// Returns the on-chain deposit address for a coin, optionally on a specific
// chain. When no chain is given the system auto-selects one.
type GetDepositAddressService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetDepositAddressService(coin string) *GetDepositAddressService {
	return &GetDepositAddressService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetDepositAddressService) SetChain(chain string) *GetDepositAddressService {
	s.params["chain"] = chain
	return s
}

// SetSize sets the deposit quantity (BTC Lightning Network only; range
// 0.000001 - 0.001).
func (s *GetDepositAddressService) SetSize(size string) *GetDepositAddressService {
	s.params["size"] = size
	return s
}

func (s *GetDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/deposit-address", s.params).WithSign()
	return request.Do[DepositAddress](req)
}

type DepositAddress struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag"`
	URL     string `json:"url"`
}

// GetSubDepositAddressService -- GET /api/v3/account/sub-deposit-address (UTA mgt. read)
//
// Returns the on-chain deposit address of a sub-account for a coin, optionally
// on a specific chain.
type GetSubDepositAddressService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubDepositAddressService(subUid, coin string) *GetSubDepositAddressService {
	return &GetSubDepositAddressService{c: c, params: map[string]string{
		"subUid": subUid,
		"coin":   coin,
	}}
}

func (s *GetSubDepositAddressService) SetChain(chain string) *GetSubDepositAddressService {
	s.params["chain"] = chain
	return s
}

// SetSize sets the deposit quantity (BTC Lightning Network only; range
// 0.000001 - 0.001).
func (s *GetSubDepositAddressService) SetSize(size string) *GetSubDepositAddressService {
	s.params["size"] = size
	return s
}

func (s *GetSubDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/sub-deposit-address", s.params).WithSign()
	return request.Do[DepositAddress](req)
}

// GetDepositRecordsService -- GET /api/v3/account/deposit-records (UTA mgt. read)
//
// Returns the unified account's deposit records within a start/end time window,
// paginated by cursor.
type GetDepositRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetDepositRecordsService(startTime, endTime time.Time) *GetDepositRecordsService {
	return &GetDepositRecordsService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetCoin filters records to a coin (all coins when omitted).
func (s *GetDepositRecordsService) SetCoin(coin string) *GetDepositRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetDepositRecordsService) SetOrderID(orderId string) *GetDepositRecordsService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetDepositRecordsService) SetLimit(limit string) *GetDepositRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetDepositRecordsService) SetCursor(cursor string) *GetDepositRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetDepositRecordsService) Do(ctx context.Context) ([]DepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/deposit-records", s.params).WithSign()
	resp, err := request.Do[[]DepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type DepositRecord struct {
	OrderID     string          `json:"orderId"`
	RecordID    string          `json:"recordId"`
	Coin        string          `json:"coin"`
	Type        string          `json:"type"`
	Dest        string          `json:"dest"`
	Size        decimal.Decimal `json:"size"`
	Status      string          `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Chain       string          `json:"chain"`
	CreatedTime time.Time       `json:"createdTime"`
	UpdatedTime time.Time       `json:"updatedTime"`
}

// GetSubDepositRecordsService -- GET /api/v3/account/sub-deposit-records (UTA mgt. read)
//
// Returns a sub-account's deposit records within a start/end time window,
// paginated by cursor.
type GetSubDepositRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubDepositRecordsService(subUid string, startTime, endTime time.Time) *GetSubDepositRecordsService {
	return &GetSubDepositRecordsService{c: c, params: map[string]string{
		"subUid":    subUid,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetCoin filters records to a coin (all coins when omitted).
func (s *GetSubDepositRecordsService) SetCoin(coin string) *GetSubDepositRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetSubDepositRecordsService) SetLimit(limit string) *GetSubDepositRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetSubDepositRecordsService) SetCursor(cursor string) *GetSubDepositRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSubDepositRecordsService) Do(ctx context.Context) ([]DepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/sub-deposit-records", s.params).WithSign()
	resp, err := request.Do[[]DepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
