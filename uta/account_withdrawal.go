package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// WithdrawService -- POST /api/v3/account/withdrawal (UTA withdrawal)
//
// Submits a withdrawal, either on-chain or as an internal transfer. For on-chain
// withdrawals chain is required; for internal transfers address holds the
// recipient UID/email/mobile per innerToType. The reply returns the created
// order's identifiers.
type WithdrawService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewWithdrawService(coin, transferType, address string, size decimal.Decimal) *WithdrawService {
	return &WithdrawService{c: c, body: map[string]any{
		"coin":         coin,
		"transferType": transferType,
		"address":      address,
		"size":         size.String(),
	}}
}

// SetChain sets the blockchain network (e.g. erc20, trc20); required when
// transferType is on_chain.
func (s *WithdrawService) SetChain(chain string) *WithdrawService {
	s.body["chain"] = chain
	return s
}

// SetInnerToType sets the internal withdrawal address type (uid, email, or
// mobile; defaults to uid).
func (s *WithdrawService) SetInnerToType(innerToType string) *WithdrawService {
	s.body["innerToType"] = innerToType
	return s
}

// SetAreaCode sets the area code; required when innerToType is mobile.
func (s *WithdrawService) SetAreaCode(areaCode string) *WithdrawService {
	s.body["areaCode"] = areaCode
	return s
}

// SetTag sets the address tag (required for certain coins such as EOS).
func (s *WithdrawService) SetTag(tag string) *WithdrawService {
	s.body["tag"] = tag
	return s
}

func (s *WithdrawService) SetRemark(remark string) *WithdrawService {
	s.body["remark"] = remark
	return s
}

func (s *WithdrawService) SetClientOid(clientOid string) *WithdrawService {
	s.body["clientOid"] = clientOid
	return s
}

// SetMemberCode sets the member code (bithumb, korbit, coinone).
func (s *WithdrawService) SetMemberCode(memberCode string) *WithdrawService {
	s.body["memberCode"] = memberCode
	return s
}

// SetIdentityType sets the identity type (company or user).
func (s *WithdrawService) SetIdentityType(identityType string) *WithdrawService {
	s.body["identityType"] = identityType
	return s
}

// SetCompanyName sets the company name (required when identityType is company).
func (s *WithdrawService) SetCompanyName(companyName string) *WithdrawService {
	s.body["companyName"] = companyName
	return s
}

// SetFirstName sets the first name (required when identityType is user).
func (s *WithdrawService) SetFirstName(firstName string) *WithdrawService {
	s.body["firstName"] = firstName
	return s
}

// SetLastName sets the last name (required when identityType is user).
func (s *WithdrawService) SetLastName(lastName string) *WithdrawService {
	s.body["lastName"] = lastName
	return s
}

// SetAccountType sets the account type(s) to deduct from (funding, uta, otc;
// comma-separated).
func (s *WithdrawService) SetAccountType(accountType string) *WithdrawService {
	s.body["accountType"] = accountType
	return s
}

func (s *WithdrawService) Do(ctx context.Context) (*WithdrawResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/withdrawal", s.body).WithSign()
	return request.Do[WithdrawResult](req)
}

type WithdrawResult struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// GetWithdrawalRecordsService -- GET /api/v3/account/withdrawal-records (UTA mgt. read)
//
// Returns the account's withdrawal records within the given time window,
// optionally filtered by coin or order identifier, paginated by cursor.
type GetWithdrawalRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetWithdrawalRecordsService(startTime, endTime time.Time) *GetWithdrawalRecordsService {
	return &GetWithdrawalRecordsService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetWithdrawalRecordsService) SetCoin(coin string) *GetWithdrawalRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetWithdrawalRecordsService) SetOrderID(orderId string) *GetWithdrawalRecordsService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetWithdrawalRecordsService) SetClientOid(clientOid string) *GetWithdrawalRecordsService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetWithdrawalRecordsService) SetLimit(limit string) *GetWithdrawalRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetWithdrawalRecordsService) SetCursor(cursor string) *GetWithdrawalRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetWithdrawalRecordsService) Do(ctx context.Context) ([]WithdrawalRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/withdrawal-records", s.params).WithSign()
	resp, err := request.Do[[]WithdrawalRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type WithdrawalRecord struct {
	OrderID     string          `json:"orderId"`
	ClientOid   string          `json:"clientOid"`
	RecordID    string          `json:"recordId"`
	Coin        string          `json:"coin"`
	Type        string          `json:"type"`
	Dest        string          `json:"dest"`
	Size        decimal.Decimal `json:"size"`
	Status      string          `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Chain       string          `json:"chain"`
	Fee         decimal.Decimal `json:"fee"`
	Confirm     string          `json:"confirm"`
	Tag         string          `json:"tag"`
	CreatedTime time.Time       `json:"createdTime"`
	UpdatedTime time.Time       `json:"updatedTime"`
}

// GetWithdrawAddressService -- GET /api/v3/account/withdraw-address (UTA mgt. read)
//
// Returns the account's saved withdrawal address book entries, optionally
// filtered by coin or address category, paginated by cursor.
type GetWithdrawAddressService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetWithdrawAddressService() *GetWithdrawAddressService {
	return &GetWithdrawAddressService{c: c, params: map[string]string{}}
}

func (s *GetWithdrawAddressService) SetCoin(coin string) *GetWithdrawAddressService {
	s.params["coin"] = coin
	return s
}

// SetType sets the address book category (EVM, regular, universal, internal).
func (s *GetWithdrawAddressService) SetType(addressType string) *GetWithdrawAddressService {
	s.params["type"] = addressType
	return s
}

func (s *GetWithdrawAddressService) SetLimit(limit string) *GetWithdrawAddressService {
	s.params["limit"] = limit
	return s
}

func (s *GetWithdrawAddressService) SetCursor(cursor string) *GetWithdrawAddressService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetWithdrawAddressService) Do(ctx context.Context) (*WithdrawAddressBook, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/withdraw-address", s.params).WithSign()
	return request.Do[WithdrawAddressBook](req)
}

type WithdrawAddressBook struct {
	AddressList []WithdrawAddress `json:"addressList"`
	Cursor      string            `json:"cursor"`
}

type WithdrawAddress struct {
	Coin         string    `json:"coin"`
	Chain        string    `json:"chain"`
	Address      string    `json:"address"`
	CountryCode  string    `json:"countryCode"`
	Label        string    `json:"label"`
	Memo         string    `json:"memo"`
	Type         string    `json:"type"`
	InternalType string    `json:"internalType"`
	CreatedTime  time.Time `json:"createdTime"`
}

// CancelWithdrawalService -- POST /api/v3/account/cancel-withdrawal (UTA withdrawal)
//
// Cancels a withdrawal still within its cooling-off period. Either orderId or
// clientOid must be set (orderId takes precedence); batch cancellation is not
// supported. The reply data is the literal string "success".
type CancelWithdrawalService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelWithdrawalService() *CancelWithdrawalService {
	return &CancelWithdrawalService{c: c, body: map[string]any{}}
}

func (s *CancelWithdrawalService) SetOrderID(orderId string) *CancelWithdrawalService {
	s.body["orderId"] = orderId
	return s
}

func (s *CancelWithdrawalService) SetClientOid(clientOid string) *CancelWithdrawalService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelWithdrawalService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/cancel-withdrawal", s.body).WithSign()
	return request.Do[string](req)
}
