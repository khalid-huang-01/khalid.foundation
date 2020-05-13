// 415. Add Strings
// type byte = uint8
func addStrings(num1 string, num2 string) string {
	var result string
	cur1 := len(num1) - 1
	cur2 := len(num2) - 1
	var remain, carry byte

	for cur1 >= 0 || cur2 >= 0 {
		tmp := carry
		carry = 0
		if cur1 >= 0 {
			tmp += num1[cur1] - '0'
			cur1--
		}
		if cur2 >= 0 {
			tmp += num2[cur2] - '0'
			cur2--
		}
		remain = tmp % 10
		carry = tmp / 10
		result = string(remain+'0') + result
	}
	if carry > 0 {
		result = string(carry+'0') + result
	}
	return result
}