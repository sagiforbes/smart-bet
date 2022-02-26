package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type GamblingInfo struct {
	NumberOfOptions          int64
	WiningPossibilities      int64
	AlwaysLosingPossiblities int64
	PercentToBet             float32
	MoneyInWallet            float32
	SimulatorCount           int64
	Odds                     float32
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println("Failed to calculate: ", err)
		os.Exit(1)
	}

}

func readFromStdIn(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	exitOnError(err)
	return strings.TrimSpace(string(text))
}

func textToInt(text string, def ...int) int64 {
	text = strings.TrimSpace(text)
	if text == "" && len(def) > 0 {
		return int64(def[0])
	}

	num, err := strconv.Atoi(string(text))
	exitOnError((err))
	return int64(num)
}

func textToOdd(text string, def ...float32) float32 {
	text = strings.TrimSpace(text)
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		exitOnError(fmt.Errorf("Invalid Odd string"))
	}

	casino, err := strconv.ParseFloat(parts[1], 32)
	exitOnError(err)
	player, err := strconv.ParseFloat(parts[0], 32)
	exitOnError(err)

	return (float32(casino) / float32(player))
}

func playGame(gamblingInfo *GamblingInfo) []string {
	var moneyToRisk = float32(gamblingInfo.MoneyInWallet * (gamblingInfo.PercentToBet / 100.0))
	var winVal = gamblingInfo.WiningPossibilities - 1
	var totalPossibilities = gamblingInfo.NumberOfOptions
	rand.Seed(time.Now().UnixNano())
	randResult := rand.Int63n(totalPossibilities)
	randResult += gamblingInfo.AlwaysLosingPossiblities
	var win bool

	if randResult <= winVal {
		win = true
	}

	ret := make([]string, 0)
	ret = append(ret, fmt.Sprintf("%f", gamblingInfo.MoneyInWallet)) //before bet money in wallet
	ret = append(ret, fmt.Sprintf("%f", moneyToRisk))                //how much money was risked
	ret = append(ret, fmt.Sprintf("%v", win))                        //true if win
	if win {
		gamblingInfo.MoneyInWallet += moneyToRisk * gamblingInfo.Odds
	} else {
		gamblingInfo.MoneyInWallet -= moneyToRisk
	}
	ret = append(ret, fmt.Sprintf("%f", gamblingInfo.MoneyInWallet)) //how much money after gamble

	return ret
}

func runSimulator(gamblingInfo *GamblingInfo, count int64) [][]string {
	ret := make([][]string, 0)
	if count >= gamblingInfo.SimulatorCount || gamblingInfo.MoneyInWallet < 1.0 {
		return ret
	}

	gameRes := playGame(gamblingInfo)

	ret = append(ret, gameRes)

	ret = append(ret, runSimulator(gamblingInfo, count+1)...)
	return ret
}

func main() {
	gamblingInfo := GamblingInfo{}

	const (
		dOptions = 2
		dWin     = 1
		dLose    = 0

		dDefaultOdd = "1:2"
		dStart      = 100
		dSimCount   = 200
	)

	gamblingInfo.NumberOfOptions = *flag.Int64("p", dOptions, "How many possitions in the gamble")
	gamblingInfo.WiningPossibilities = *flag.Int64("w", dWin, "How many wining possitions exists")
	gamblingInfo.AlwaysLosingPossiblities = *flag.Int64("l", dLose, "How many wining possitions always lose")
	var oddsIn = *flag.String("odd", dDefaultOdd, "Odds in the form of user:casino")
	gamblingInfo.Odds = textToOdd(oddsIn)

	pWin := float32(float64(gamblingInfo.WiningPossibilities) / float64(gamblingInfo.AlwaysLosingPossiblities+gamblingInfo.NumberOfOptions))
	pLose := 1.0 - pWin
	percentageToBet := float32(float32(gamblingInfo.Odds*(pWin)-pLose)/float32(gamblingInfo.Odds)) * 100

	gamblingInfo.PercentToBet = float32(*flag.Float64("per", float64(percentageToBet), "Percentage of wallet to bet"))

	if gamblingInfo.PercentToBet < 0.1 {
		exitOnError(fmt.Errorf("its not worth to play if the gambling pecentage is %f", gamblingInfo.PercentToBet*float32(100.0)))
	}

	gamblingInfo.MoneyInWallet = float32(*flag.Float64("wallet", dStart, "How much money in wallet"))

	gamblingInfo.SimulatorCount = *flag.Int64("count", dSimCount, "How many games to play")

	simulatorResults := runSimulator(&gamblingInfo, 0)

	fmt.Println("Wallet before,Amout to risk,Win,Wallet after play")
	for _, col := range simulatorResults {
		fmt.Println(strings.Join(col, ","))
	}

}
