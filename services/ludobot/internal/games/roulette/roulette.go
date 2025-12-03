package roulette

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram/bot/models"
)

type SectorTag string
type BetState int

const (
	SectorTagZero           SectorTag = "zero"
	SectorTagRed            SectorTag = "red"
	SectorTagBlack          SectorTag = "black"
	SectorTagOdd            SectorTag = "odd"
	SectorTagEven           SectorTag = "even"
	SectorTagFirstTwelve    SectorTag = "1-12"
	SectorTagSecondTwelve   SectorTag = "13-24"
	SectorTagThirdTwelve    SectorTag = "25-36"
	SectorTagFirstEighteen  SectorTag = "1-18"
	SectorTagSecondEighteen SectorTag = "19-36"
)

const (
	BetStateRegistered BetState = iota
	BetStateWon
	BetStateLost
)

const numSectors = 37

type UserBet struct {
	UserID      int64
	TokenAmount int
	Sector      SectorTag
	State       BetState
}

type RouletteRules struct {
	blackSectors map[int]bool
	redSectors   map[int]bool
}

type Roulette struct {
	bets []UserBet
	mu   sync.RWMutex
	seed int64
}

var validSectorTags = map[SectorTag]struct{}{
	SectorTagZero:           {},
	SectorTagRed:            {},
	SectorTagBlack:          {},
	SectorTagOdd:            {},
	SectorTagEven:           {},
	SectorTagFirstTwelve:    {},
	SectorTagSecondTwelve:   {},
	SectorTagThirdTwelve:    {},
	SectorTagFirstEighteen:  {},
	SectorTagSecondEighteen: {},
}

func (r *Roulette) PlaceBet(userId int64, amount int, sector SectorTag) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bets = append(r.bets, UserBet{
		userId,
		amount,
		sector,
		BetStateRegistered,
	})
}

func (r *Roulette) Spin() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := rand.New(rand.NewSource(r.seed)).Intn(numSectors)
	r.updateBetStates(result)

	return result
}

func (r *Roulette) PlaceBetFromMessage(message *models.Message) error {
	amount, sector, err := r.extractBetAndSector(message.Text)

	if err != nil {
		return err
	} else if message.From == nil {
		return errors.New("can't get message sender")
	}

	r.PlaceBet(
		message.From.ID,
		amount,
		sector,
	)

	return nil
}

func (r *Roulette) GetRegistered() []UserBet {
	return r.getBetsByState(BetStateRegistered)
}

func (r *Roulette) GetWon() []UserBet {
	return r.getBetsByState(BetStateWon)
}

func (r *Roulette) GetLost() []UserBet {
	return r.getBetsByState(BetStateLost)
}

func (r *Roulette) ClearBets() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bets = r.bets[:0]
}

func NewRoulette() *Roulette {
	seed := time.Now().UnixNano()
	slog.Info(fmt.Sprintf("Registered new roulette with seed: %d", seed))
	return &Roulette{
		bets: make([]UserBet, 0),
		mu:   sync.RWMutex{},
		seed: seed,
	}
}

func blackSectors() []int {
	return []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}
}

func redSectors() []int {
	return []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
}

func makeTags(result int) []SectorTag {
	tags := []SectorTag{}

	if result < 0 {
		return tags
	}

	if result == 0 {
		tags = append(tags, SectorTagZero)
		return tags
	}

	if result%2 == 0 {
		tags = append(tags, SectorTagEven)
	} else {
		tags = append(tags, SectorTagOdd)
	}

	rules := newRouletteRules()

	if rules.blackSectors[result] {
		tags = append(tags, SectorTagBlack)
	} else if rules.redSectors[result] {
		tags = append(tags, SectorTagRed)
	}

	if 1 <= result && result <= 12 {
		tags = append(tags, SectorTagFirstTwelve)
	} else if 13 <= result && result <= 24 {
		tags = append(tags, SectorTagSecondTwelve)
	} else if 25 <= result && result <= 36 {
		tags = append(tags, SectorTagThirdTwelve)
	}

	if 1 <= result && result <= 18 {
		tags = append(tags, SectorTagFirstEighteen)
	} else if 19 <= result && result <= 36 {
		tags = append(tags, SectorTagSecondEighteen)
	}

	return tags
}

func newRouletteRules() RouletteRules {
	result := RouletteRules{
		blackSectors: make(map[int]bool),
		redSectors:   make(map[int]bool),
	}

	for _, num := range redSectors() {
		result.redSectors[num] = true
	}

	for _, num := range blackSectors() {
		result.blackSectors[num] = true
	}

	return result
}

func calculateBetState(betSector SectorTag, winningTags []SectorTag) BetState {
	if slices.Contains(winningTags, betSector) {
		return BetStateWon
	}

	return BetStateLost
}

func updateBetState(bet UserBet, winningTags []SectorTag) UserBet {
	result := bet
	result.State = calculateBetState(bet.Sector, winningTags)
	return result
}

func parseSectorTag(s string) (SectorTag, error) {
	tag := SectorTag(s)
	_, tagExists := validSectorTags[tag]

	if !tagExists {
		return "", errors.New("invalid bet sector")
	}

	return tag, nil
}

func (r *Roulette) updateBetStates(result int) {
	winningTags := makeTags(result)
	for i := range r.bets {
		r.bets[i] = updateBetState(r.bets[i], winningTags)
	}
}

func (r *Roulette) getBetsByState(state BetState) []UserBet {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := []UserBet{}
	for _, bet := range r.bets {
		if bet.State == state {
			result = append(result, bet)
		}
	}

	return result
}

func (r *Roulette) extractBetAndSector(message string) (int, SectorTag, error) {
	makeErrorRv := func(message string) (int, SectorTag, error) {
		slog.Warn(fmt.Sprintf("Error occured: %s", message))
		return 0, "", errors.New(message)
	}

	messageParts := strings.Split(message, " ")

	betRaw := messageParts[1]
	sectorRaw := messageParts[2]

	bet, err := strconv.Atoi(betRaw)
	if err != nil {
		return makeErrorRv("Invalid bet")
	} else if bet <= 0 {
		// todo: transfer this to validateBet func
		return makeErrorRv("Bet must be greater than zero")
	}

	sector, err := parseSectorTag(sectorRaw)
	if err != nil {
		return makeErrorRv(err.Error())
	}

	return bet, sector, nil
}
