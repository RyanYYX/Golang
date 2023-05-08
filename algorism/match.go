package algorism

// BF Brute-Force
func BF(s, p string) (index int) {
	sl, pl := len(s), len(p)
	for i := 0; i < sl-pl+1; i++ {
		if s[i:i+pl] == p {
			return i
		}
	}
	return -1
}

// RK Rabin-Karp
func RK(s, p string) (index int) {
	sl, pl := len(s), len(p)
	var hash = make(map[string][]int)
	num := ""
	for i := 0; i < sl; i++ {
		if i < pl {
			num += string(s[i])
			if i == pl-1 {
				hash[num] = append(hash[num], 0)
			}
		} else {
			num = num[1:] + string(s[i])
			hash[num] = append(hash[num], i-pl+1)
		}
	}

	for _, index := range hash[num] {
		if s[index:index+pl] == p {
			return index
		}
	}
	return -1
}

// BM Boyer-Moore
func BM(s, p string) (index int) {
	return
}

func goodSuffixRule(s string) (suffix []int, prefix []bool) {
	suffix, prefix = make([]int, len(s)), make([]bool, len(s))
	for i := 0; i < len(s); i++ {
		j, k := i, len(s)-1
		for j >= 0 && s[j] == s[k] {
			j, k = j-1, k-1
		}
		if j == -1 {
			prefix[i] = true
		}
		suffix[k] = j
	}
	return
}

// KMP Knuth-Morris-Pratt
func KMP(s, p string) (index int) {
	next, i, j := getNext(p), 0, -1
	for ; i<len(s) && j+1<len(p); i++ {
		for j > -1 && s[i] != p[j+1] {
			j = next[j]
		}
		if s[i] == p[j+1] {
			j = j + 1
		}
	}

	if j+1 == len(p) {
		return i - len(p)
	}
	return -1
}

func getNext(s string) (next []int) {
	next = make([]int, len(s))
	next[0]= -1
	for i,j:=1,-1; i<len(s); i++ {
		for j > -1 && s[i] != s[j+1] {
			j = next[j]
		}
		if s[i] == s[j+1] {
			j = j + 1
		}
		next[i] = j
	}
	return
}