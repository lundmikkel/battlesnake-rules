package maps

import (
	"sort"

	"github.com/BattlesnakeOfficial/rules"
)

type ClassicSnakeMap struct{}

func init() {
	globalRegistry.RegisterMap("classic_snake", ClassicSnakeMap{})
}

func (m ClassicSnakeMap) ID() string {
	return "classic_snake"
}

func (m ClassicSnakeMap) Meta() Metadata {
	return Metadata{
		Name:        "Classic Snake",
		Description: "Classic Snake food behavior with exactly one food on the board",
		Author:      "Battlesnake",
		Version:     1,
		MinPlayers:  1,
		MaxPlayers:  16,
		BoardSizes:  OddSizes(rules.BoardSizeSmall, rules.BoardSizeXXLarge),
		Tags:        []string{TAG_FOOD_PLACEMENT},
	}
}

func (m ClassicSnakeMap) SetupBoard(initialBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(0)

	if len(initialBoardState.Snakes) > int(m.Meta().MaxPlayers) {
		return rules.ErrorTooManySnakes
	}

	snakeIDs := make([]string, 0, len(initialBoardState.Snakes))
	for _, snake := range initialBoardState.Snakes {
		snakeIDs = append(snakeIDs, snake.ID)
	}

	tempBoardState, err := rules.CreateDefaultBoardState(rand, initialBoardState.Width, initialBoardState.Height, snakeIDs)
	if err != nil {
		return err
	}

	for _, snake := range tempBoardState.Snakes {
		editor.PlaceSnake(snake.ID, snake.Body, snake.Health)
	}

	return applyClassicFoodRules(rand, tempBoardState, editor)
}

func (m ClassicSnakeMap) PreUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	return nil
}

func (m ClassicSnakeMap) PostUpdateBoard(lastBoardState *rules.BoardState, settings rules.Settings, editor Editor) error {
	rand := settings.GetRand(lastBoardState.Turn)
	foodCount := len(lastBoardState.Food)

	if foodCount == 1 {
		return nil
	}

	if foodCount > 1 {
		selected, ok := selectDeterministicFood(lastBoardState.Food)
		if !ok {
			return nil
		}
		editor.ClearFood()
		editor.AddFood(selected)
		return nil
	}

	return applyClassicFoodRules(rand, lastBoardState, editor)
}

func applyClassicFoodRules(rand rules.Rand, boardState *rules.BoardState, editor Editor) error {
	if selected, ok := selectDeterministicFood(boardState.Food); ok {
		editor.ClearFood()
		editor.AddFood(selected)
		return nil
	}

	unoccupiedPoints := rules.GetUnoccupiedPoints(boardState, false, true)
	if len(unoccupiedPoints) == 0 {
		return nil
	}

	placeFoodRandomlyAtPositions(rand, boardState, editor, 1, unoccupiedPoints)
	return nil
}

func selectDeterministicFood(food []rules.Point) (rules.Point, bool) {
	if len(food) == 0 {
		return rules.Point{}, false
	}

	foodCopy := append([]rules.Point(nil), food...)
	sort.Slice(foodCopy, func(i, j int) bool {
		if foodCopy[i].X == foodCopy[j].X {
			return foodCopy[i].Y < foodCopy[j].Y
		}
		return foodCopy[i].X < foodCopy[j].X
	})

	return foodCopy[0], true
}
