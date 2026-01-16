package utilities

import (
    "math"
)

// ComputeSimilarity calculates the normalized Levenshtein similarity between two strings.
// Returns a value between 0.0 and 1.0.
func ComputeSimilarity(s1, s2 string) float64 {
    if s1 == s2 {
        return 1.0
    }
    if len(s1) == 0 || len(s2) == 0 {
        return 0.0
    }
    
    distance := computeLevenshtein(s1, s2)
    maxLen := float64(math.Max(float64(len(s1)), float64(len(s2))))
    
    if maxLen == 0 {
        return 1.0
    }

    return 1.0 - (float64(distance) / maxLen)
}

func computeLevenshtein(s1, s2 string) int {
    r1, r2 := []rune(s1), []rune(s2)
    n, m := len(r1), len(r2)
    
    // Optimization: Make sure n <= m to use O(min(n,m)) space
    if n > m {
        r1, r2 = r2, r1
        n, m = m, n
    }

    currentRow := make([]int, n+1)
    for i := 0; i <= n; i++ {
        currentRow[i] = i
    }

    for i := 1; i <= m; i++ {
        previousRow := currentRow // Keep reference to previous row
        // Create new row (or could swap buffers if careful)
        // Here, allocating new slice is safer/easier for readability in Go without generics-heavy optimizations
        currentRow = make([]int, n+1) 
        currentRow[0] = i
        
        for j := 1; j <= n; j++ {
            cost := 0
            if r1[j-1] != r2[i-1] {
                cost = 1
            }
            
            add := previousRow[j] + 1
            del := currentRow[j-1] + 1
            change := previousRow[j-1] + cost
            
            currentRow[j] = min(add, min(del, change))
        }
    }
    return currentRow[n]
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
