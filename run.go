package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ArmanMaesumi/chess"
)

func main() {
	var depth int8
	depth = 5

	main_menu(depth)
}

func main_menu(depth int8) {
	fmt.Println("----------------")
	fmt.Println("CHESS MENU\n1- New Game\n2- Load Game\n3- Quit")
	fmt.Println("----------------")

	var fen string
	main_loop := true

	for main_loop {
		readerOptions := bufio.NewReader(os.Stdin)
		fmt.Print("Option -> ")
		option, _ := readerOptions.ReadString('\n')
		fmt.Println("Option: ", option)
		option = option[:len(option)-1]

		switch option {
		case "1":
			fmt.Println("New Game")
			fen = ""
			main_loop = false
		case "2":
			fmt.Println("Load Game")
			fen = fen_input()
			main_loop = false
		case "3":
			fmt.Println("Quit")
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}

	game := create_game(fen)
	game.AddTagPair("Event", "F/S Return Match")
	game_loop(depth, game)
}

func create_game(fen string) *chess.Game {
	if strings.Compare(fen, "") == 0 {
		fmt.Println("Default new game")
		game := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}))
		return game
	} else {
		fmt.Println("Loading game")
		fen, _ := chess.FEN(fen)
		game := chess.NewGame(fen)
		return game
	}
}

func fen_input() string {
	readerFen := bufio.NewReader(os.Stdin)
	fmt.Print("FEN -> ")
	fen, _ := readerFen.ReadString('\n')
	return fen
}

func game_loop(depth int8, game *chess.Game) {
	fmt.Println("Game Loop")

	//var movesHistory []string
	round := 1
	//lastMoveHistory := movesHistory
	lastGameState := game.Clone()
	lastRound := round

	for true {
		if game.Outcome() != "*" {
			fmt.Println("END")
			fmt.Println(game.Method())
			fmt.Println(game.Outcome())
			fmt.Println(game.Position().Board().Draw())
			fmt.Println(game.String())
			os.Exit(0)
		} else {
			fmt.Println("\n===================")
			fmt.Println("Current state of the board:")
			fmt.Println(game.Position().Board().Draw())
			fmt.Println(game.String())
			fmt.Println("----------------")

			fmt.Println("Round: ", round)
			fmt.Println("Turn: ", game.Position().Turn())
			fmt.Println("1- Add movement with prediction\n2- Add only movement\n3- Remove last movement\n4- Quit")
			fmt.Println("----------------")

			readerSubOptions := bufio.NewReader(os.Stdin)
			fmt.Print("Option -> ")
			subOption, _ := readerSubOptions.ReadString('\n')
			fmt.Println("Option: ", subOption)
			subOption = subOption[:len(subOption)-1]
			fmt.Println("----------------")

			switch subOption {
			case "1":
				fmt.Println("Predicting movement")
				gameClone := game.Clone()
				deepening(depth, gameClone, round)
				fallthrough
			case "2":
				lastGameState = game.Clone()
				lastRound = round
				fmt.Println("Add enemy movement")
				fmt.Println(game.ValidMoves())
				move := move_input()
				err := game.MoveStr(move)
				if err != nil {
					fmt.Println("Invalid move")
					fmt.Println(err)
					break
				}
				round++
				break
			case "3":
				fmt.Println("Remove last movement")
				if round > 1 {
					//movesHistory = lastMoveHistory
					game = lastGameState.Clone()
					round = lastRound
					fmt.Println(game.Position().Board().Draw())
					fmt.Println(game.String())
				}
				break
			case "4":
				fmt.Println("Quit")
				os.Exit(0)
			default:
				fmt.Println("Invalid option")
			}
			fmt.Println("===================\n")
		}
	}
}

func move_input() string {
	readerMove := bufio.NewReader(os.Stdin)
	fmt.Print("Move -> ")
	moveString, _ := readerMove.ReadString('\n')
	moveString = moveString[:len(moveString)-1]
	fmt.Println("Move: ", moveString)
	// moveString = "e2e4" in *chess.Move
	return moveString
}

// func quick_test1(depth int8) {
//	fen, _ := chess.FEN("2kr4/pp1n1p2/q1p1p3/2P5/1P1PPbp1/P1N3Pr/6K1/R2Q1R2 w - - 0 27")
//	game := chess.NewGame(fen)
//	iterative_deepening(depth, game)
//}

func iterative_deepening(max_depth int8, game *chess.Game) {
	maximize := true
	if game.Position().Turn() == chess.Black {
		maximize = false
	}

	var lastRecommendedMove *chess.Move

	for i := int8(3); i < max_depth+1; i++ {
		fmt.Println("(", "Depth: ", i, "MaxDepth: ", max_depth, ")")
		lastRecommendedMove = minimax_root(i, game, maximize, false, 1)
		fmt.Println("\n")
	}
	fmt.Println("\nLast recomended move: ", lastRecommendedMove)
	game.Move(lastRecommendedMove)
	fmt.Println("\nRecomended state of the board:")
	fmt.Println(game.Position().Board().Draw())
}

func deepening(depth int8, game *chess.Game, round int) {
	maximize := true
	if game.Position().Turn() == chess.Black {
		maximize = false
	}
	fmt.Println("Searching for best move at depth: ", depth)
	minimax_root(depth, game, maximize, true, round)
	fmt.Println("\nRecomended state of the board:")
	fmt.Println(game.Position().Board().Draw())
}
