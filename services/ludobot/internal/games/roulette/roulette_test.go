package roulette

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSetEvenAndOddTags(t *testing.T) {
	for i := 1; i < numSectors; i++ {
		// given / when
		tags := makeTags(i)

		// then
		if i%2 == 0 {
			assert.Containsf(t, tags, SectorTagEven, "No even tag for %d", i)
		} else {
			assert.Containsf(t, tags, SectorTagOdd, "No odd tag for %d", i)
		}
	}
}

func TestShouldSetZeroTagAndNoOtherTagsToZero(t *testing.T) {
	// given
	zero := 0

	// when
	tags := makeTags(zero)

	// then
	assert.Equal(t, 1, len(tags), "extra tags on zero")
	assert.Containsf(t, tags, SectorTagZero, "no zero tag on sector %d", zero)
}

func TestShouldSetBlackTagOnBlackSectors(t *testing.T) {
	for _, sec := range blackSectors() {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagBlack, "no black tag on black sector %d", sec)
	}
}

func TestShouldSetRedTagOnRedSectors(t *testing.T) {
	for _, sec := range redSectors() {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagRed, "no red tag on red sector %d", sec)
	}
}

func TestShouldSetFirstTwelveTagOnFirstTwelveSectors(t *testing.T) {
	for sec := 1; sec <= 12; sec++ {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagFirstTwelve, "no first twelve tag on sector %d", sec)
	}
}

func TestShouldSetSecondTwelveTagOnSecondTwelveSectors(t *testing.T) {
	for sec := 13; sec <= 24; sec++ {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagSecondTwelve, "no second twelve tag on sector %d", sec)
	}
}

func TestShouldSetThirdTwelveTagOnThirdTwelveSectors(t *testing.T) {
	for sec := 25; sec <= 36; sec++ {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagThirdTwelve, "no third twelve tag on sector %d", sec)
	}
}

func TestShouldSetFirstEighteenTagOnFirstEighteenSectors(t *testing.T) {
	for sec := 1; sec <= 18; sec++ {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagFirstEighteen, "no first eighteen tag on sector %d", sec)
	}
}

func TestShouldSetSecondEighteenTagOnSecondEighteenSectors(t *testing.T) {
	for sec := 19; sec <= 36; sec++ {
		// given / when
		tags := makeTags(sec)

		// then
		assert.Containsf(t, tags, SectorTagSecondEighteen, "no second eighteen tag on sector %d", sec)
	}
}

func TestShouldFilterWonAndLostBets(t *testing.T) {
	// given
	roulette := Roulette{}
	roulette.PlaceBet(123, 1, SectorTagRed)
	roulette.PlaceBet(321, 1, SectorTagBlack)

	// when
	roulette.updateBetStates(redSectors()[0])

	won := roulette.GetWon()
	lost := roulette.GetLost()

	// then
	assert.Equalf(t, 1, len(won), "Exactly 1 user was supposed to win, instead %d had won", len(won))
	assert.Equalf(t, 1, len(lost), "Exactly 1 user was supposed to lose, instead %d had lost", len(lost))

	assert.Equal(t, int64(123), won[0].UserID, "User 123 was supposed to win bet on red")
	assert.Equal(t, int64(321), lost[0].UserID, "User 321 was supposed to lose bet on black")
}

func TestShouldAddBetsConcurrently(t *testing.T) {
	// given
	roulette := Roulette{}
	num_bets := 10000000

	// when
	var wg sync.WaitGroup
	for i := range num_bets {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			roulette.PlaceBet(int64(i), 10, SectorTagRed)
		}(i)
	}
	wg.Wait()

	// then
	assert.Equal(t, num_bets, len(roulette.GetRegistered()), "Some bets weren't registered")
}
