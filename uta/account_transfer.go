package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetTransferableCoinsService -- GET /api/v3/account/transferable-coins (UTA mgt. read)
//
// Returns the coins that can be transferred between the given source and target
// account types. The reply data is an array of coin name strings.
//
// Account types: spot, p2p, coin_futures, usdt_futures, usdc_futures,
// crossed_margin, isolated_margin, uta.
type GetTransferableCoinsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetTransferableCoinsService(fromType, toType string) *GetTransferableCoinsService {
	return &GetTransferableCoinsService{c: c, params: map[string]string{
		"fromType": fromType,
		"toType":   toType,
	}}
}

func (s *GetTransferableCoinsService) Do(ctx context.Context) ([]string, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/transferable-coins", s.params).WithSign()
	resp, err := request.Do[[]string](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TransferService -- POST /api/v3/account/transfer (UTA mgt. read & write)
//
// Transfers a coin between two account types within the same account. Account
// types: spot, p2p, coin_futures, usdt_futures, usdc_futures, crossed_margin,
// isolated_margin, uta.
type TransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewTransferService(fromType, toType, coin string, amount decimal.Decimal) *TransferService {
	return &TransferService{c: c, body: map[string]any{
		"fromType": fromType,
		"toType":   toType,
		"coin":     coin,
		"amount":   amount.String(),
	}}
}

// SetSymbol sets the isolated spot margin trading pair (e.g. BTCUSDT).
func (s *TransferService) SetSymbol(symbol string) *TransferService {
	s.body["symbol"] = symbol
	return s
}

// SetAllowBorrow enables ("yes") or disables ("no") automatic margin borrowing
// when the balance is insufficient.
func (s *TransferService) SetAllowBorrow(allowBorrow string) *TransferService {
	s.body["allowBorrow"] = allowBorrow
	return s
}

// SetClientOid sets the client-defined transaction identifier.
func (s *TransferService) SetClientOid(clientOid string) *TransferService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *TransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/transfer", s.body).WithSign()
	return request.Do[TransferResult](req)
}

type TransferResult struct {
	TransferId string `json:"transferId"`
	ClientOid  string `json:"clientOid"`
}

// SubTransferService -- POST /api/v3/account/sub-transfer (UTA mgt. read & write)
//
// Transfers a coin between a main account and a sub-account. Requires the main
// account API key. Account types: spot, p2p, usdt_futures, coin_futures,
// usdc_futures, crossed_margin, uta.
type SubTransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSubTransferService(fromType, toType, coin string, amount decimal.Decimal, fromUserId, toUserId, clientOid string) *SubTransferService {
	return &SubTransferService{c: c, body: map[string]any{
		"fromType":   fromType,
		"toType":     toType,
		"coin":       coin,
		"amount":     amount.String(),
		"fromUserId": fromUserId,
		"toUserId":   toUserId,
		"clientOid":  clientOid,
	}}
}

// SetAllowBorrow enables ("yes") or disables ("no") automatic margin borrowing
// when the balance is insufficient.
func (s *SubTransferService) SetAllowBorrow(allowBorrow string) *SubTransferService {
	s.body["allowBorrow"] = allowBorrow
	return s
}

func (s *SubTransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/sub-transfer", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// GetSubTransferRecordsService -- GET /api/v3/account/sub-transfer-record (UTA mgt. read)
//
// Returns the main/sub-account transfer records, paginated by cursor and bounded
// to a 90-day query span.
type GetSubTransferRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubTransferRecordsService() *GetSubTransferRecordsService {
	return &GetSubTransferRecordsService{c: c, params: map[string]string{}}
}

// SetSubUid filters to a sub-account UID; omit to retrieve main account records.
func (s *GetSubTransferRecordsService) SetSubUid(subUid string) *GetSubTransferRecordsService {
	s.params["subUid"] = subUid
	return s
}

// SetRole filters by account role ("initiator" or "receiver"; default
// "initiator").
func (s *GetSubTransferRecordsService) SetRole(role string) *GetSubTransferRecordsService {
	s.params["role"] = role
	return s
}

func (s *GetSubTransferRecordsService) SetCoin(coin string) *GetSubTransferRecordsService {
	s.params["coin"] = coin
	return s
}

// SetStartTime filters records at or after t (max 90-day span).
func (s *GetSubTransferRecordsService) SetStartTime(t time.Time) *GetSubTransferRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 90-day span).
func (s *GetSubTransferRecordsService) SetEndTime(t time.Time) *GetSubTransferRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSubTransferRecordsService) SetClientOid(clientOid string) *GetSubTransferRecordsService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetSubTransferRecordsService) SetLimit(limit string) *GetSubTransferRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetSubTransferRecordsService) SetCursor(cursor string) *GetSubTransferRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSubTransferRecordsService) Do(ctx context.Context) (*SubTransferRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/sub-transfer-record", s.params).WithSign()
	return request.Do[SubTransferRecords](req)
}

type SubTransferRecords struct {
	List   []SubTransferRecord `json:"list"`
	Cursor string              `json:"cursor"`
}

type SubTransferRecord struct {
	TransferId    string          `json:"transferId"`
	FromType      string          `json:"fromType"`
	ToType        string          `json:"toType"`
	Amount        decimal.Decimal `json:"amount"`
	Coin          string          `json:"coin"`
	FromUserId    string          `json:"fromUserId"`
	ToUserId      string          `json:"toUserId"`
	Status        string          `json:"status"`
	ClientOid     string          `json:"clientOid"`
	CreatedTime   time.Time       `json:"createdTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
	OldTransferId string          `json:"oldTransferId"`
}

// SubMasterTransferService -- POST /api/v3/account/sub-master-transfer (UTA mgt. read & write)
//
// Transfers a coin from a sub-account back to the main account. The sub-account
// initiating the transfer must be the API key holder. Account types: from spot
// or uta; to spot, p2p, or uta.
type SubMasterTransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSubMasterTransferService(fromType, toType, coin string, amount decimal.Decimal) *SubMasterTransferService {
	return &SubMasterTransferService{c: c, body: map[string]any{
		"fromType": fromType,
		"toType":   toType,
		"coin":     coin,
		"amount":   amount.String(),
	}}
}

// SetClientOid sets the client order ID (max 64 characters).
func (s *SubMasterTransferService) SetClientOid(clientOid string) *SubMasterTransferService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *SubMasterTransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/sub-master-transfer", s.body).WithSign()
	return request.Do[TransferResult](req)
}
