package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"slices"
	"strconv"
	"sync"
)

type Player struct {
	Conn    net.Conn
	Name    string
	Symbol  string
	Scanner *bufio.Scanner
}

type GameSession struct {
	Players     []*Player
	Board       []string
	PlayerIndex int // Index of the player whose turn it is
	wg          sync.WaitGroup
}

var (
	mut                 sync.Mutex
	sessions            []*GameSession
	xSymbol             = "X"
	oSymbol             = "O"
	winningCombinations = [][3]int{
		{0, 1, 2}, // Horizontal
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6}, // Vertical
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8}, // Diagonal
		{2, 4, 6},
	}
)

// WaitingPlayer informs a player to wait for their opponent's move.
func (session *GameSession) WaitingPlayer(indexPlayer int) {
	player := session.Players[indexPlayer]
	log.Printf("Player %s is waiting for their opponent's move", player.Name)
	player.Conn.Write([]byte("---\nWaiting for opponent's move\n"))
	session.wg.Done()
}

// CheckWin checks if the specified player has won.
func (session *GameSession) CheckWin(indexPlayer int) bool {
	symbolPlayer := session.Players[indexPlayer].Symbol
	for _, combination := range winningCombinations {
		if session.Board[combination[0]] == symbolPlayer &&
			session.Board[combination[1]] == symbolPlayer &&
			session.Board[combination[2]] == symbolPlayer {
			log.Printf("Player %s with symbol %s has won", session.Players[indexPlayer].Name, symbolPlayer)
			return true
		}
	}
	return false
}

// ViewBoard formats the game board for display.
func (session *GameSession) ViewBoard() string {
	var sBoard string
	for ind, cell := range session.Board {
		if (ind+1)%3 == 0 {
			sBoard += fmt.Sprintf("%s\n", cell)
		} else {
			sBoard += fmt.Sprintf("%s |", cell)
		}
	}
	return "---\n" + sBoard
}

// PlayerMove handles a player's move.
func (session *GameSession) PlayerMove(indexPlayer int) {
	player := session.Players[indexPlayer]
	message := fmt.Sprintf("Your turn. You are playing as %s\nEnter the number of an empty cell:\n", player.Symbol)
	player.Conn.Write([]byte(message))
	player.Conn.Write([]byte(session.ViewBoard()))
	for {
		player.Scanner.Scan()
		cell := player.Scanner.Text()
		cellInt, err := strconv.Atoi(cell)
		if err != nil {
			player.Conn.Write([]byte("Invalid input. Try again:\n"))
			continue
		}
		if slices.Contains(session.Board, cell) {
			session.Board[cellInt-1] = player.Symbol
			player.Conn.Write([]byte(session.ViewBoard()))
			break
		} else {
			player.Conn.Write([]byte("Cell is unavailable. Try again:\n"))
		}
	}
	session.wg.Done()
}

// GetNewBoard creates a new game board.
func GetNewBoard() []string {
	return []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
}

// addPlayerToSession adds a player to a session or creates a new one
func addPlayerToSession(player *Player) *GameSession {
	mut.Lock()
	defer mut.Unlock()

	if sessions == nil || len(sessions[len(sessions)-1].Players) == 2 {
		session := &GameSession{Board: GetNewBoard()}
		session.Players = append(session.Players, player)
		sessions = append(sessions, session)
		player.Conn.Write([]byte("Waiting for a second player\n"))
		return session
	}

	session := sessions[len(sessions)-1]
	session.Players = append(session.Players, player)
	return session
}

// endGame is responsible for ending the game
func endGame(session *GameSession, winnerIndex int) {
	if winnerIndex >= 0 {
		winner := session.Players[winnerIndex]
		loser := session.Players[1-winnerIndex]

		message := fmt.Sprintf("---\n%s\nPlayer %s wins! Game over.\n", session.ViewBoard(), winner.Name)
		log.Printf("Game over. Player %s wins.", winner.Name)

		winner.Conn.Write([]byte(message))
		loser.Conn.Write([]byte(message))
	} else {
		drawMessage := "---\nIt's a draw! Game over.\n"
		log.Printf("Game ended in a draw.")
		for _, player := range session.Players {
			player.Conn.Write([]byte(drawMessage))
		}
	}
}

// startGame
func startGame(session *GameSession) {
	player1 := session.Players[0]
	player2 := session.Players[1]
	defer player1.Conn.Close()
	defer player2.Conn.Close()

	firstMove := rand.Intn(2)
	session.PlayerIndex = firstMove
	player1.Symbol = xSymbol
	player2.Symbol = oSymbol

	for i := 0; i < 9; i++ {
		playerMoveIndex := session.PlayerIndex
		playerWaitIndex := 1 - playerMoveIndex

		session.wg.Add(2)
		go session.PlayerMove(playerMoveIndex)
		go session.WaitingPlayer(playerWaitIndex)
		session.wg.Wait()

		if session.CheckWin(playerMoveIndex) {
			endGame(session, playerMoveIndex)
			return
		}
		session.PlayerIndex = 1 - session.PlayerIndex
	}
	endGame(session, -1)
}

func main() {
	ip := "192.168.1.79"
	port := "8080"
	host := ip + ":" + port
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error creating server: %s\n", err)
		return
	}
	log.Println("Server started on port 8080")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %s\n", err)
			continue
		}
		scanner := bufio.NewScanner(conn)
		conn.Write([]byte("---\nEnter your name: "))
		scanner.Scan()
		playerName := scanner.Text()

		player := &Player{Conn: conn, Name: playerName, Scanner: bufio.NewScanner(conn)}
		log.Printf("Player %s (%s) connected", playerName, conn.RemoteAddr())

		session := addPlayerToSession(player)
		if len(session.Players) == 2 {
			log.Printf("Starting game between %s and %s", session.Players[0].Name, session.Players[1].Name)
			go startGame(session)
		}
	}
}
