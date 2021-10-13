package others

import (
	"fmt"
	"strconv"
	"strings"
)

// leetcode 811, 直接使用切分后利用哈希表就可以了
func subdomainVisits(cpdomains []string) []string {
	counts := make(map[string]int)
	for _, cpdomain := range cpdomains {
		tmp := strings.Split(cpdomain, " ")
		count, _ :=  strconv.Atoi(tmp[0])
		domains := strings.Split(tmp[1], ".")
		cur := ""
		for i := len(domains) - 1; i  >= 0; i-- {
			if cur == "" {
				cur = domains[i]
			} else {
				cur = domains[i] + "." + cur
			}
			counts[cur] += count
		}
	}
	rsl := make([]string, 0)
	for domain, count := range counts {
		rsl = append(rsl, fmt.Sprintf("%d %s", count, domain))
	}
	return rsl
}