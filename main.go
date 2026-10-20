package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

func evaluate(perm []int, t [][]int, d, c []int) int {
	m := len(t)
	n := len(perm)
	machineReady := make([]int, m)
	completion := make([]int, n)
	for _, job := range perm {
		prevFinish := 0
		for j := 0; j < m; j++ {
			start := machineReady[j]
			if prevFinish > start {
				start = prevFinish
			}
			finish := start + t[j][job]
			machineReady[j] = finish
			prevFinish = finish
		}
		completion[job] = prevFinish
	}
	penalty := 0
	for i := 0; i < n; i++ {
		if completion[i] > d[i] {
			penalty += c[i] * (completion[i] - d[i])
		}
	}
	return penalty
}
func F1(perm []int, L int) [][]int {
	if L <= 0 {
		return nil
	}
	n := len(perm)
	result := make([][]int, 0)
	var rec func(cur []int, depth int)
	rec = func(cur []int, depth int) {
		if depth == 0 {
			cp := make([]int, n)
			copy(cp, cur)
			result = append(result, cp)
			return
		}
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				cur[i], cur[j] = cur[j], cur[i]
				rec(cur, depth-1)
				cur[i], cur[j] = cur[j], cur[i]
			}
		}
	}
	rec(perm, L)
	seen := make(map[string]bool, len(result))
	unique := make([][]int, 0, len(result))
	for _, p := range result {
		key := fmt.Sprint(p)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, p)
		}
	}
	return unique
}
func F2(omega [][]int, i int) []int {
	if i < 0 || i >= len(omega) {
		return nil
	}
	return omega[i]
}
func permString(p []int) string {
	var sb strings.Builder
	sb.WriteString("(")
	for i, v := range p {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(strconv.Itoa(v + 1))
	}
	sb.WriteString(")")
	return sb.String()
}
func hillClimb(n int, N int, L int, t [][]int, d, c []int) (bestPerm []int, bestVal int) {
	minS := rand.Perm(n)
	minVal := evaluate(minS, t, d, c)
	fmt.Println("=== Hill Climbing (transposition metric) ===")
	fmt.Printf("Start: minS = %v, min = %v\n", permString(minS), minVal)
	i := 0
	for i < N {
		omega := F1(minS, L)
		bestNeighbor := []int(nil)
		bestNeighborVal := math.MaxInt
		for idx := 0; idx < len(omega); idx++ {
			si := F2(omega, idx)
			if si == nil {
				continue
			}
			v := evaluate(si, t, d, c)
			if v < bestNeighborVal {
				bestNeighborVal = v
				bestNeighbor = si
			}
		}
		if bestNeighbor == nil {
			break
		}
		fmt.Printf("Step %v: current min = %v, minS = %v, best neighbor = %v (val = %v)\n", i+1, minVal, permString(minS), permString(bestNeighbor), bestNeighborVal)
		if bestNeighborVal < minVal {
			minS = bestNeighbor
			minVal = bestNeighborVal
			i++
			fmt.Printf("Update: new min = %v, new minS = %v\n", minVal, permString(minS))
		} else {
			fmt.Println("No improvement, local minimum reached")
			break
		}
	}
	fmt.Println("=== Result ===")
	fmt.Printf("minS = %v\n", permString(minS))
	fmt.Printf("min = %v\n", minVal)
	return minS, minVal
}
func buildSchedule(perm []int, t [][]int) (x, y [][]int) {
	n := len(perm)
	m := len(t)
	x = make([][]int, n)
	y = make([][]int, n)
	for i := range x {
		x[i] = make([]int, m)
		y[i] = make([]int, m)
	}
	machineReady := make([]int, m)
	for _, job := range perm {
		prevFinish := 0
		for j := 0; j < m; j++ {
			start := machineReady[j]
			if prevFinish > start {
				start = prevFinish
			}
			x[job][j] = start
			y[job][j] = start + t[j][job]
			machineReady[j] = y[job][j]
			prevFinish = y[job][j]
		}
	}
	return x, y
}

func writeResults(filename string, perm []int, t [][]int, d, c []int, x, y [][]int) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, "\n\n===================")
	n := len(perm)
	m := len(t)
	fmt.Fprintf(f, "Permutation: %v\n\n", permString(perm))
	fmt.Fprintf(f, "%-8v", "Job")
	for i := 0; i < m; i++ {
		fmt.Fprintf(f, "%-18v", fmt.Sprintf("Machine %v", i+1))
	}
	fmt.Fprintf(f, "%-14v%-12v%-8v\n", "Completion", "Due", "Fine")
	for job := 0; job < n; job++ {
		fmt.Fprintf(f, "%-8v", job+1)
		for j := 0; j < m; j++ {
			cell := fmt.Sprintf("%v-%v", x[job][j], y[job][j])
			fmt.Fprintf(f, "%-18v", cell)
		}
		finish := y[job][m-1]
		tard := finish - d[job]
		if tard < 0 {
			tard = 0
		}
		fmt.Fprintf(f, "%-14v%-12v%-8v\n", finish, d[job], c[job]*tard)
	}
	return nil
}
func main() {
	data1, err1 := os.ReadFile("param.txt")
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
		return
	}
	input1 := strings.Fields(string(data1))
	if len(input1) != 2 {
		fmt.Println("param.txt: need two numbers (m n)")
		return
	}
	m, err2 := strconv.Atoi(input1[0])
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		return
	}
	n, err3 := strconv.Atoi(input1[1])
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
		return
	}
	data2, err4 := os.ReadFile("time_w.txt")
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
		return
	}
	lines1 := strings.Split(string(data2), "\n")
	var t [][]int
	for _, line := range lines1 {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		row := make([]int, n)
		for i, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			row[i] = v
		}
		t = append(t, row)
	}
	if len(t) != m {
		fmt.Printf("Error: the file contains %v lines, but %v were expected\n", len(t), m)
		return
	}
	for i, row := range t {
		if len(row) != n {
			fmt.Printf("Error: in row %v, there are %v columns, but %v were expected\n", i, len(row), n)
			return
		}
	}
	data3, err5 := os.ReadFile("time_st.txt")
	if err5 != nil {
		fmt.Printf("Error: %v\n", err5)
		return
	}
	lines2 := strings.Split(string(data3), "\n")
	d := make([]int, n)
	c := make([]int, n)
	for i, line := range lines2 {
		if i == 1 || i == 3 {
			line = strings.TrimSpace(line)
			parts := strings.Fields(line)
			for j, p := range parts {
				v, err := strconv.Atoi(p)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					return
				}
				if i == 1 {
					d[j] = v
				} else {
					c[j] = v
				}
			}
		}
	}
	fmt.Println("Enter search depth N (number of steps):")
	var N int
	if _, err := fmt.Scan(&N); err != nil || N <= 0 {
		fmt.Println("Invalid N, using 100")
		N = 100
	}
	fmt.Println("Enter neighborhood radius L (transposition distance, usually 1):")
	var L int
	if _, err := fmt.Scan(&L); err != nil || L <= 0 || L > 2 {
		fmt.Println("Invalid L, using 1")
		L = 1
	}
	bestPerm, bestVal := hillClimb(n, N, L, t, d, c)
	x, y := buildSchedule(bestPerm, t)
	if err := writeResults("result.txt", bestPerm, t, d, c, x, y); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Saved %v with fine = %v\n", permString(bestPerm), bestVal)
}
