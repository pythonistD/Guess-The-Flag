package algo

func WordDistance(word1, word2 string) int {
	// Переводим в руны, так как строка индексируется по байтам
	w1Runes := []rune{
		' ',
	}
	w1Runes = append(w1Runes, []rune(word1)...)
	// добавляем нулевой элемент, чтобы не делать вычитания индексов
	w2Runes := []rune{
		' ',
	}
	w2Runes = append(w2Runes, []rune(word2)...)

	dp := make([][]int, len(w1Runes))

	for i := range dp {
		dp[i] = make([]int, len(w2Runes))
		dp[i][0] = i
	}
	for i := range dp[0] {
		dp[0][i] = i
	}

	for i := 1; i < len(dp); i++ {
		for j := 1; j < len(dp[0]); j++ {
			if w1Runes[i] == w2Runes[j] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				insert := dp[i][j-1] + 1
				del := dp[i-1][j] + 1
				replace := dp[i-1][j-1] + 1
				dp[i][j] = minOf3(insert, del, replace)
			}
		}
	}
	return dp[len(w1Runes)-1][len(w2Runes)-1]
}

func minOf3(a, b, c int) int {
	res := a
	if res > b {
		res = b
	}
	if res > c {
		res = c
	}
	return res
}
