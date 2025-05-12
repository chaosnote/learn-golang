package main

import (
	"context"
	"fmt"
	"log"

	"github.com/looplab/fsm"
)

// 定義遊戲狀態
const (
	StateIdle         = "idle"      // 遊戲尚未開始
	StateDealingCards = "dealing"   // 發牌中
	StatePreflop      = "preflop"   // 前翻牌圈
	StateFlop         = "flop"      // 翻牌圈
	StateTurn         = "turn"      // 轉牌圈
	StateRiver        = "river"     // 河牌圈
	StateShowdown     = "showdown"  // 攤牌
	StateGameOver     = "game_over" // 遊戲結束
)

// 定義遊戲事件
const (
	EventStartGame                = "start_game"
	EventDealCards                = "deal_cards"
	EventPreflopBettingRoundStart = "preflop_betting_start"
	EventPlayerBet                = "player_bet"
	EventPlayerCall               = "player_call"
	EventPlayerRaise              = "player_raise"
	EventPlayerFold               = "player_fold"
	EventPlayerCheck              = "player_check"
	EventPreflopBettingRoundEnd   = "preflop_betting_end"
	EventFlop                     = "flop"
	EventFlopBettingRoundStart    = "flop_betting_start"
	EventFlopBettingRoundEnd      = "flop_betting_end"
	EventTurn                     = "turn"
	EventTurnBettingRoundStart    = "turn_betting_start"
	EventTurnBettingRoundEnd      = "turn_betting_end"
	EventRiver                    = "river"
	EventRiverBettingRoundStart   = "river_betting_start"
	EventRiverBettingRoundEnd     = "river_betting_end"
	EventShowdown                 = "showdown"
	EventEndGame                  = "end_game"
)

// Game 結構體包含狀態機和遊戲數據
type Game struct {
	fsm *fsm.FSM
	// 遊戲相關數據，例如玩家列表、彩池金額、當前下注額等
	potAmount     int
	currentBet    int
	playersInHand []string
	// ... 其他遊戲數據
}

// NewGame 創建一個新的德州撲克遊戲實例
func NewGame() *Game {
	g := &Game{
		fsm: fsm.NewFSM(
			StateIdle,
			fsm.Events{
				{Name: EventStartGame, Src: []string{StateIdle}, Dst: StateDealingCards},
				{Name: EventDealCards, Src: []string{StateDealingCards}, Dst: StatePreflop},
				{Name: EventPreflopBettingRoundStart, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPlayerBet, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPlayerCall, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPlayerRaise, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPlayerFold, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPlayerCheck, Src: []string{StatePreflop}, Dst: StatePreflop},
				{Name: EventPreflopBettingRoundEnd, Src: []string{StatePreflop}, Dst: StateFlop},
				{Name: EventFlop, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventFlopBettingRoundStart, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventPlayerBet, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventPlayerCall, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventPlayerRaise, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventPlayerFold, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventPlayerCheck, Src: []string{StateFlop}, Dst: StateFlop},
				{Name: EventFlopBettingRoundEnd, Src: []string{StateFlop}, Dst: StateTurn},
				{Name: EventTurn, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventTurnBettingRoundStart, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventPlayerBet, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventPlayerCall, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventPlayerRaise, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventPlayerFold, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventPlayerCheck, Src: []string{StateTurn}, Dst: StateTurn},
				{Name: EventTurnBettingRoundEnd, Src: []string{StateTurn}, Dst: StateRiver},
				{Name: EventRiver, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventRiverBettingRoundStart, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventPlayerBet, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventPlayerCall, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventPlayerRaise, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventPlayerFold, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventPlayerCheck, Src: []string{StateRiver}, Dst: StateRiver},
				{Name: EventRiverBettingRoundEnd, Src: []string{StateRiver}, Dst: StateShowdown},
				{Name: EventShowdown, Src: []string{StateShowdown}, Dst: StateGameOver},
				{Name: EventEndGame, Src: []string{StateShowdown, StatePreflop, StateFlop, StateTurn, StateRiver}, Dst: StateGameOver},
			},
			fsm.Callbacks{
				"enter_state":  func(ctx context.Context, evt *fsm.Event) { log.Println(evt.FSM.Current()) },
				"before_event": func(ctx context.Context, evt *fsm.Event) { log.Println(evt.FSM.Current()) },
			},
		),
		potAmount:     0,
		currentBet:    0,
		playersInHand: []string{},
	}
	return g
}

// CurrentState 返回當前遊戲狀態
func (g *Game) CurrentState() string {
	return g.fsm.Current()
}

// TriggerEvent 觸發遊戲事件
func (g *Game) TriggerEvent(event string, args ...interface{}) error {
	err := g.fsm.Event(context.Background(), event, args)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	game := NewGame()

	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: idle

	err := game.TriggerEvent(EventStartGame)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: dealing

	err = game.TriggerEvent(EventDealCards)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: preflop

	err = game.TriggerEvent(EventPreflopBettingRoundStart)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: preflop

	// 模擬玩家行動
	err = game.TriggerEvent(EventPlayerCheck, "SANTHOSH")
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	err = game.TriggerEvent(EventPlayerBet, "PETER", 100) // 假設 PETER 下注 100
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	game.currentBet = 100 // 更新當前下注額
	err = game.TriggerEvent(EventPlayerCall, "STANLEY", 100)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	err = game.TriggerEvent(EventPlayerRaise, "BRANDON", 200) // 假設 BRANDON 加注到 200
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	game.currentBet = 200 // 更新當前下注額
	err = game.TriggerEvent(EventPlayerCall, "KEATING", 200)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	err = game.TriggerEvent(EventPlayerCall, "RAHUL", 200)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}

	err = game.TriggerEvent(EventPreflopBettingRoundEnd)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: flop

	err = game.TriggerEvent(EventFlop)
	if err != nil {
		fmt.Println("觸發事件失敗:", err)
	}
	fmt.Println("當前狀態:", game.CurrentState()) // 輸出: 當前狀態: flop

	// ... 後續遊戲流程
}
