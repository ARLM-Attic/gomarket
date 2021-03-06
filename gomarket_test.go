
package gomarket

import (
	"testing"
)

type TestActor struct {
	*StandardTrader
	BuySums Resources
	SellSums Resources
	BuyPrices Resources
	SellPrices Resources
}
func NewTestActor() *TestActor {
	t := &TestActor{
		nil,
		make(Resources),
		make(Resources),
		make(Resources), 
		make(Resources)}
	t.StandardTrader = NewStandardTrader(t)
	return t
}
func (a *TestActor) Buy(bid, ask *Order, price float64) {
	a.BuySums[bid.Resource] = a.BuySums[bid.Resource] + bid.Units
	a.BuyPrices[bid.Resource] = price
	ask.Carrier.Deliver(bid, ask, price)
}
func (a *TestActor) Deliver(bid, ask *Order, price float64) {
	a.SellSums[ask.Resource] = a.SellSums[ask.Resource] + ask.Units
	a.SellPrices[ask.Resource] = price
}

func Check(t *testing.T,
	ask_units, ask_prices, bid_units, bid_prices []float64, 
	expected_price float64, 
	expected_sells, expected_buys []float64) {
	m := NewMarket()
	sellers := make([]*TestActor, 0)
	buyers := make([]*TestActor, 0)
	product := "product"
	for i := 0; i < len(ask_units); i++ {
		seller := NewTestActor()
		seller.Ask(ask_units[i], product, ask_prices[i])
		m.Add(seller)
		sellers = append(sellers, seller)
	}
	for i := 0; i < len(bid_units); i++ {
		buyer := NewTestActor()
		buyer.Bid(bid_units[i], product, bid_prices[i])
		m.Add(buyer)
		buyers = append(buyers, buyer)
	}
	m.Trade()
	if m.prices[product] != expected_price {
		t.Error("When selling",ask_units,"for",ask_prices,"and buying",bid_units,"for",bid_prices,"expected price to be",expected_price,"but was",m.prices[product])
	}
	for i := 0; i < len(expected_sells); i++ {
		if sellers[i].SellSums[product] != expected_sells[i] {
			t.Error("When selling",ask_units,"for",ask_prices,"and buying",bid_units,"for",bid_prices,"expected seller",i,"to sell",expected_sells[i],"units, but sold",sellers[i].SellSums[product],"units.")
		}
		if sellers[i].SellSums[product] > 0 && sellers[i].SellPrices[product] != expected_price {
			t.Error("When selling",ask_units,"for",ask_prices,"and buying",bid_units,"for",bid_prices,"expected seller",i,"to sell for",expected_price,"but sold for",sellers[i].SellPrices[product])
		}
	}
	for i := 0; i < len(expected_buys); i++ {
		if buyers[i].BuySums[product] != expected_buys[i] {
			t.Error("When selling",ask_units,"for",ask_prices,"and buying",bid_units,"for",bid_prices,"expected buyer",i,"to buy",expected_buys[i],"units, but bought",buyers[i].BuySums[product],"units.")
		}
		if buyers[i].BuySums[product] > 0 && buyers[i].BuyPrices[product] != expected_price {
			t.Error("When selling",ask_units,"for",ask_prices,"and buying",bid_units,"for",bid_prices,"expected buyer",i,"to buy for",expected_price,"but bought for",buyers[i].BuyPrices[product])
		}
	}
}

func TestOneSellerNoBuyers(t *testing.T) {
	Check(t,
		[]float64{10.0}, []float64{5.0}, []float64{}, []float64{},
		5.0,
		[]float64{0.0}, []float64{})
}

func TestNoSellersOneBuyer(t *testing.T) {
	Check(t,
		[]float64{}, []float64{}, []float64{10.0}, []float64{10.0},
		10.0,
		[]float64{}, []float64{0.0})
}

func TestOneSellerOneBuyerNoDeal(t *testing.T) {
	Check(t,
		[]float64{10.0}, []float64{5.0}, []float64{10.0}, []float64{2.0},
		3.5,
		[]float64{0.0}, []float64{0.0})
}

func TestOneSellerOneBuyerDeal(t *testing.T) {
	Check(t,
		[]float64{10.0}, []float64{5.0}, []float64{10.0}, []float64{10.0},
		7.5,
		[]float64{10.0}, []float64{10.0})
}

func TestOneSellerManyBuyersDeal(t *testing.T) {
	Check(t,
		[]float64{10.0}, []float64{5.0}, []float64{5.0, 5.0, 5.0}, []float64{10.0, 15.0, 12.0},
		12.0,
		[]float64{10.0}, []float64{0.0, 5.0, 5.0})
}

func TestOneSellerManyBuyersNoDeal(t *testing.T) {
	Check(t,
		[]float64{10.0}, []float64{25.0}, []float64{5.0, 5.0, 5.0}, []float64{10.0, 15.0, 12.0},
		20.0,
		[]float64{0.0}, []float64{0.0, 0.0, 0.0})
}

func TestManySellersOneBuyerDeal(t *testing.T) {
	Check(t,
		[]float64{7.0, 5.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0}, []float64{10.0},
		5.0,
		[]float64{5.0, 5.0, 0.0}, []float64{10.0})
}

func TestManySellersOneBuyerNoDeal(t *testing.T) {
	Check(t,
		[]float64{7.0, 5.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0}, []float64{1.0},
		2.5,
		[]float64{0.0, 0.0, 0.0}, []float64{0.0})
}

func TestManySellersManyBuyersDealScarcity(t *testing.T) {
	Check(t,
		[]float64{7.0, 5.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0, 20.0}, []float64{20.0, 10.0},
		10.0,
		[]float64{7.0, 5.0, 15.0}, []float64{10.0, 17.0})
}

func TestManySellersManyBuyersDealAbundance(t *testing.T) {
	Check(t,
		[]float64{7.0, 25.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0, 20.0}, []float64{20.0, 10.0},
		5.0,
		[]float64{5.0, 25.0, 0.0}, []float64{10.0, 20.0})
}

func TestManySellersManyBuyersNoDealScarcity(t *testing.T) {
	Check(t,
		[]float64{7.0, 5.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0, 20.0}, []float64{2.0, 1.0},
		3.0,
		[]float64{0.0, 0.0, 0.0}, []float64{0.0, 0.0})
}

func TestManySellersManyBuyersNoDealAbundance(t *testing.T) {
	Check(t,
		[]float64{7.0, 25.0, 15.0}, []float64{5.0, 4.0, 6.0}, []float64{10.0, 20.0}, []float64{2.0, 1.0},
		3.0,
		[]float64{0.0, 0.0, 0.0}, []float64{0.0, 0.0})
}

func TestValue(t *testing.T) {
	m := NewMarket()
	rice := "rice"
	shoes := "shoes"
	m.prices[rice] = 2.0
	m.prices[shoes] = 3.0
	if m.Value(Resources{rice: 4.0, shoes: 5.0}) != 23.0 {
		t.Error("wrong price!")
	}
}

func TestEq(t *testing.T) {
	rice := "rice"
	shoes := "shoes"
	r1 := Resources(Resources{rice: 4.0, shoes: 3.2})
	r2 := Resources(Resources{rice: 4.0, shoes: 3.2})
	r3 := Resources(Resources{rice: 4.0})
	r4 := Resources(Resources{rice: 4.0, shoes: 3.1})
	if !r1.Eq(r2) {
		t.Error("should be equal!")
	}
	if r1.Eq(r3) {
		t.Error("should not be equal!")
	}
	if r1.Eq(r4) {
		t.Error("should not be equal!")
	}
}

func TestMergeIn(t *testing.T) {
	rice := "rice"
	shoes := "shoes"
	r1 := Resources(Resources{rice: 4.0, shoes: 3.2})
	r2 := Resources(Resources{rice: 1.0, shoes: 0.2})
	r1.MergeIn(r2)
	if r1.Eq(Resources{rice: 5.0, shoes: 3.4}) {
		t.Error("wrong merge!")
	}
}

func TestClone(t *testing.T) {
	rice := "rice"
	shoes := "shoes"
	r1 := Resources(Resources{rice: 4.0, shoes: 3.2})
	if !r1.Eq(r1.Clone()) {
		t.Error("not equal!")
	}
}
