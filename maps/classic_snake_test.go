package maps_test

import (
	"testing"

	"github.com/BattlesnakeOfficial/rules"
	"github.com/BattlesnakeOfficial/rules/maps"
	"github.com/stretchr/testify/require"
)

func TestClassicSnakeMapInterface(t *testing.T) {
	var _ maps.GameMap = maps.ClassicSnakeMap{}
}

func TestClassicSnakeMapSetupBoardSingleFood(t *testing.T) {
	m := maps.ClassicSnakeMap{}
	initialBoardState := rules.NewBoardState(7, 7)
	settings := rules.Settings{}.WithRand(rules.MinRand)

	nextBoardState := rules.NewBoardState(7, 7)
	editor := maps.NewBoardStateEditor(nextBoardState)

	err := m.SetupBoard(initialBoardState, settings, editor)

	require.NoError(t, err)
	require.Len(t, nextBoardState.Food, 1)
}

func TestClassicSnakeMapPostUpdateKeepsSingleFood(t *testing.T) {
	m := maps.ClassicSnakeMap{}
	initialBoardState := rules.NewBoardState(7, 7).WithFood([]rules.Point{{X: 3, Y: 3}})
	settings := rules.Settings{}.WithRand(rules.MinRand)

	nextBoardState := initialBoardState.Clone()
	editor := maps.NewBoardStateEditor(nextBoardState)

	err := m.PostUpdateBoard(initialBoardState, settings, editor)

	require.NoError(t, err)
	require.Len(t, nextBoardState.Food, 1)
	require.Equal(t, initialBoardState.Food, nextBoardState.Food)
}

func TestClassicSnakeMapPostUpdateAddsFoodAfterEaten(t *testing.T) {
	m := maps.ClassicSnakeMap{}
	initialBoardState := rules.NewBoardState(3, 3)
	settings := rules.Settings{}.WithRand(rules.MinRand)

	nextBoardState := initialBoardState.Clone()
	editor := maps.NewBoardStateEditor(nextBoardState)

	err := m.PostUpdateBoard(initialBoardState, settings, editor)

	require.NoError(t, err)
	require.Len(t, nextBoardState.Food, 1)
	require.Equal(t, rules.Point{X: 0, Y: 0}, nextBoardState.Food[0])
}

func TestClassicSnakeMapPostUpdateTrimsExtraFood(t *testing.T) {
	m := maps.ClassicSnakeMap{}
	initialBoardState := rules.NewBoardState(7, 7).WithFood([]rules.Point{{X: 5, Y: 2}, {X: 1, Y: 4}, {X: 1, Y: 3}})
	settings := rules.Settings{}.WithRand(rules.MinRand)

	nextBoardState := initialBoardState.Clone()
	editor := maps.NewBoardStateEditor(nextBoardState)

	err := m.PostUpdateBoard(initialBoardState, settings, editor)

	require.NoError(t, err)
	require.Len(t, nextBoardState.Food, 1)
	require.Equal(t, []rules.Point{{X: 1, Y: 3}}, nextBoardState.Food)
}

func TestClassicSnakeMapPostUpdateNoRoomForFood(t *testing.T) {
	m := maps.ClassicSnakeMap{}
	initialBoardState := rules.NewBoardState(2, 2).WithSnakes([]rules.Snake{{
		ID:     "1",
		Body:   []rules.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 0}, {X: 1, Y: 1}},
		Health: 100,
	}})
	settings := rules.Settings{}.WithRand(rules.MinRand)

	nextBoardState := initialBoardState.Clone()
	editor := maps.NewBoardStateEditor(nextBoardState)

	err := m.PostUpdateBoard(initialBoardState, settings, editor)

	require.NoError(t, err)
	require.Len(t, nextBoardState.Food, 0)
}
