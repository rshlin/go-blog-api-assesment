package logical

import (
	"log"
	"unicode"
)

// NumDecodings
//
// # Dynamic Programming Solution
//
// Space Complexity: O(n)
//
// Time Complexity: O(n)
func NumDecodings(s string) int {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			log.Fatalf("Error: The input string must contain only digits.")
		}
	}

	n := len(s)
	if n == 0 {
		return 0
	}

	dp := make([]int, n+1)
	dp[0] = 1
	if s[0] == '0' {
		dp[1] = 0
	} else {
		dp[1] = 1
	}

	for i := 2; i <= n; i++ {
		if s[i-1] != '0' {
			dp[i] += dp[i-1]
		}

		if s[i-2] == '1' || (s[i-2] == '2' && s[i-1] <= '6') {
			dp[i] += dp[i-2]
		}
	}

	return dp[n]
}
