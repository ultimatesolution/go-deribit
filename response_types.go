package deribit

import "encoding/json"

// TradeResponse is the data returned by a trade event notification
type TradeResponse struct {
	TradeID       int         `json:"tradeId"`
	Timestamp     int64       `json:"timeStamp"`
	Instrument    string      `json:"instrument"`
	Quantity      json.Number `json:"quantity"`
	Price         float64     `json:"price"`
	State         string      `json:"state"`
	Direction     string      `json:"direction"`
	OrderID       int         `json:"orderId"`
	MatchingID    int         `json:"matchingId"`
	MakerComm     float64     `json:"makerComm"`
	TakerComm     float64     `json:"takerComm"`
	IndexPrice    float64     `json:"indexPrice"`
	Label         string      `json:"label"`
	Me            string      `json:"me"`
	TickDirection int         `json:"tickDirection"`
	TradeSeq      int64       `json:"tradeSeq"`
	// T - if subscriber is taker, M - if subscriber is maker
	Liquidity string `json:"liquidity"`
}

type PositionResponse struct {
	Instrument string `json:"instrument"`
	// The type of instrument. "future" or "option"
	Kind              string  `json:"kind"`
	Size              int     `json:"size"`
	AveragePrice      float64 `json:"averagePrice"`
	Direction         string  `json:"direction"`
	SizeBTC           float64 `json:"sizeBtc"`
	FloatingPl        float64 `json:"floatingPl"`
	RealizedPl        float64 `json:"realizedPl"`
	EstLiqPrice       float64 `json:"estLiqPrice"`
	MarkPrice         float64 `json:"markPrice"`
	IndexPrice        float64 `json:"indexPrice"`
	MaintenanceMargin float64 `json:"maintenanceMargin"`
	InitialMargin     float64 `json:"initialMargin"`
	SettlementPrice   float64 `json:"settlementPrice"`
	Delta             float64 `json:"delta"`
	OpenOrderMargin   float64 `json:"openOrderMargin"`
	ProfitLoss        float64 `json:"profitLoss"`
}

type PortfolioResponse struct {
	Currency          string  `json:"currency"`
	Equity            float64 `json:"equity"`
	MaintenanceMargin float64 `json:"maintenanceMargin"`
	InitialMargin     float64 `json:"initialMargin"`
	AvailableFunds    float64 `json:"availableFunds"`
	unrealizedPl      float64 `json:"unrealizedPl"`
	realizedPl        float64 `json:"realizedPl"`
	totalPl           float64 `json:"totalPl"`
}

type PortfolioEvent struct {
	Portfolio []PortfolioResponse `json:"portfolio"`
	Positions []PositionResponse  `json:"positions"`
}

// OrderBookResponse is the data returned by an orderbook change
type OrderBookResponse struct {
	State           string            `json:"state"`
	SettlementPrice float64           `json:"settlementPrice"`
	Instrument      string            `json:"instrument"`
	Timestamp       int64             `json:"tstamp"`
	Last            float64           `json:"last"`
	Low             float64           `json:"low"`
	High            float64           `json:"high"`
	Mark            float64           `json:"mark"`
	Bids            []*OrderBookEntry `json:"bids"`
	Asks            []*OrderBookEntry `json:"asks"`
}

// OrderBookEntry is an entry in the orderbook
type OrderBookEntry struct {
	Quantity json.Number `json:"quantity"`
	Price    float64     `json:"price"`
	Cm       float64     `json:"cm"`
	CmAmount float64     `json:"cm_amount"`
}

// OrderResponse is a response to an OrderRequest
// It contains two fields: the created order in order and a list of the resulting trades in trades
type OrderResponse struct {
	Order  *OrderResponseDetail   `json:"order"`
	Trades []*OrderResponseTrades `json:"trades"`
}

// OrderResponseTrades trades is populated when a trade is executed immediately
type OrderResponseTrades struct {
	Label      string      `json:"label"`
	SelfTrade  bool        `json:"selfTrade"`
	Quantity   int         `json:"quantity"`
	Price      float64     `json:"price"`
	TradeSeq   int         `json:"tradeSeq"`
	MatchingID json.Number `json:matchingId`
}

// OrderResponseDetail is the full details of the order
type OrderResponseDetail struct {
	OrderID        json.Number `json:"orderId"`
	Direction      string      `json:"direction"`
	FilledQuantity int         `json:"filledQuantity"`
	Quantity       int         `json:"quantity"`
	AvgPrice       float64     `json:"avgPrice"`
	Price          float64     `json:"price"`
	Label          string      `json:"label"`
	Commission     float64     `json:"commission"`
	Created        int64       `json:"created"`
	LastUpdate     int64       `json:"lastUpdate"`
	State          string      `json:"state"`
	API            bool        `json:"api"`
	Triggered      bool        `json:"triggered"`
}
