// HousingListings.go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// ListingInfo structure
//------------------------------------------------------
type ListingInfo struct {
	StreetAddress string
	Price         int
	Sold          bool
}

// Condo listing info
type Condo struct {
	ListingInfo
	CondoFees int // Condos require fees
}

// Helper structure for lotsize
type Rectangle struct {
	Width float32
	Depth float32
}

// House listing info
type House struct {
	ListingInfo
	Lotsize Rectangle // Houses are on a lot
}

// Townhouse
type TownHouse struct {
	ListingInfo
	FloorLevels int // Townhouses have multiple levels
}

// Interface for interacting with ListingInfo
//-----------------------------------------------------------------------------------------------------------------
type RealEstate interface {
	getStreetAddress() string
	getPrice() int
	getSold() bool
	setSold(s bool)
}

// Implementation for interfacing with Condo
//Condo
func (L *Condo) getStreetAddress() string {
	return L.StreetAddress
}
func (L *Condo) getPrice() int {
	return L.Price
}
func (L *Condo) getSold() bool {
	return L.Sold
}
func (L *Condo) setSold(s bool) {
	L.Sold = s
}

//House
func (L *House) getStreetAddress() string {
	return L.StreetAddress
}
func (L *House) getPrice() int {
	return L.Price
}
func (L *House) getSold() bool {
	return L.Sold
}
func (L *House) setSold(s bool) {
	L.Sold = s
}

//TownHouse
func (L *TownHouse) getStreetAddress() string {
	return L.StreetAddress
}
func (L *TownHouse) getPrice() int {
	return L.Price
}
func (L *TownHouse) getSold() bool {
	return L.Sold
}
func (L *TownHouse) setSold(s bool) {
	L.Sold = s
}

// Buyer structure
type Buyer struct {
	Name   string
	Active bool // If active Buyer will participate in bidding
	// Internal helper variables to control bidding process
	bidMinimum int
	bidMaximum int
	bidStep    int
	bidCurrent int
	bidDelay   time.Duration
}

// Buyer Factory
func NewBuyer(name string) *Buyer {
	// randomly initialize the Buyer
	var b Buyer
	b.Name = name
	n := rand.Intn(10) // maximum 10 second delay
	b.bidDelay = time.Duration((n+5)*10) * time.Millisecond
	b.bidMinimum = rand.Intn(10)*25000 + 600000
	b.bidMaximum = int((1.1 + float64(rand.Intn(50))/100.0) * float64(b.bidMinimum))
	b.bidStep = (b.bidMaximum - b.bidMinimum) / (10 + rand.Intn(10))
	b.Active = true
	fmt.Println(b)
	return &b
}

// Call to receive bid from buyer
// second return parameter will be false if bid is invalid
func (b *Buyer) nextBid() (int, bool) {
	if !b.Active {
		return 0, false
	}
	time.Sleep(b.bidDelay)
	if b.bidCurrent < b.bidMinimum {
		b.bidCurrent = b.bidMinimum
		return b.bidCurrent, true
	} else {
		if b.bidCurrent < b.bidMaximum {
			b.bidCurrent = b.bidCurrent + b.bidStep
			return b.bidCurrent, true
		}
	}
	b.bidCurrent = b.bidMinimum
	return b.bidCurrent, false
}

// Seller of real estate objects
type Seller struct {
	Name         string
	Object       RealEstate
	OfferChan    chan int  // Channel for concurrent reception of bids
	ResponseChan chan bool // Channel for response to bid
	bidAccept    int       // threshold for acceptance
}

// Direct bid submission not using channels
func (s *Seller) acceptBid(offer int) bool {
	if !s.Object.getSold() && offer >= s.bidAccept {
		s.Object.setSold(true)
		return true
	} else {
		return false
	}
}

// Factory function for seller
func NewSeller(name string, obj RealEstate) *Seller {
	var s Seller
	s.Name = name
	s.Object = obj
	s.bidAccept = int(0.95 * float64(obj.getPrice()))
	s.OfferChan = make(chan int)
	s.ResponseChan = make(chan bool)
	go func() {
		for {
			select {
			case offer := <-s.OfferChan:
				s.ResponseChan <- s.acceptBid(offer)
			}
		}
	}()
	return &s
}

// Helper function to test if any Buyer is active
func buyerActive(allBuyers []*Buyer) bool {
	for _, b := range allBuyers {
		if b.Active {
			return true
		}
	}
	return false
}

// Helper function if any real estate is not sold yet
func objectForSale(allSellers []*Seller) bool {
	for _, s := range allSellers {
		if !s.Object.getSold() {
			return true
		}
	}
	return false
}

func main() {
	// Seeding the pseudo randon number generator; resulting in different bidding processes
	rand.Seed(time.Now().UnixNano())

	//Real estate listing Scenario
	listings := []RealEstate{
		&Condo{ListingInfo{"Goulburn Ave 1120", 750000, false}, 900},
		&Condo{ListingInfo{"Summerset Street 10", 950000, false}, 850},
		&Condo{ListingInfo{"Wilbord Avenue 999", 1250000, false}, 1250}}
	listings = append(listings, &TownHouse{ListingInfo{"Elgin 123", 2100000, false}, 2})
	listings = append(listings, &House{ListingInfo{"Maplewood 889", 850000, false}, Rectangle{50, 110}})

	// Seller for every listing
	sellers := []*Seller{NewSeller("Eve", listings[0]), NewSeller("Monica", listings[1]), NewSeller("Ramon", listings[2])}
	sellers = append(sellers, NewSeller("Paul", listings[3]))
	sellers = append(sellers, NewSeller("Mary", listings[4]))

	// Buyers
	buyers := []*Buyer{NewBuyer("Zara"), NewBuyer("Jim"), NewBuyer("Claude"), NewBuyer("Emilie"),
		NewBuyer("Amelie"), NewBuyer("Ali")}

	// Bidding process
	for _, buy := range buyers { 
		// This is the bidding process for one buyer
		for amount, valid := buy.nextBid(); valid; amount, valid = buy.nextBid() {
			// Buyers try to buy any property for current bidding amount
			go func() {
				for _, s := range sellers {
					if s.Object.getSold() {
						continue
					}
					fmt.Printf("%s bids %d on %s.\n", buy.Name, amount, s.Object.getStreetAddress())
					// Is bid accepted? - close deal

					if _, ok := <-s.OfferChan; ok {
						fmt.Printf("Buyer %s buys from Seller %s the Object %s for $ %d\n", buy.Name,
							s.Name, s.Object.getStreetAddress(), amount)
						buy.Active = false
						break
					} /**/
					/*
						if s.acceptBid(amount) {
							fmt.Printf("Buyer %s buys from Seller %s the Object %s for $ %d\n", buy.Name,
								s.Name, s.Object.getStreetAddress(), amount)
							buy.Active = false
							break
						} /**/
				}
			}()
		}
		buy.Active = false // Buyer went to upper limit or was successful
	}
	// Synchronization
	// Ensure that we only exit if bidding process is complete
	for buyerActive(buyers) && objectForSale(sellers) {
		time.Sleep(50 * time.Millisecond)
	}
}
