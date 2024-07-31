package config

type PrizeConfig struct {
	AllowedPrizes []int
}

type GamePrizeConfig struct {
	prizesByPlayersCount map[int][]PrizeConfig
}

var (
	GamePrizeConfigConst = GamePrizeConfig{
		prizesByPlayersCount: map[int][]PrizeConfig{
			6: {
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3, 4, 5},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3, 4, 5},
				},
			},
			10: {
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3, 5},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3},
				},
			},
			16: {
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2, 3},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2},
				},
			},
			20: {
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2},
				},
				PrizeConfig{
					AllowedPrizes: []int{0, 1, 2},
				},
			},
		},
	}
)

func (c *GamePrizeConfig) GetConfig(playersAmount int, place int) []int {
	if place >= len(c.prizesByPlayersCount[playersAmount]) {
		return make([]int, 0)
	}

	return c.prizesByPlayersCount[playersAmount][place].AllowedPrizes
}
