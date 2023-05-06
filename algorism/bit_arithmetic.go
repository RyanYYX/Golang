package algorism

// add: addendA + addendB = result
func add(addendA, addendB int64) (result int64) {
	if addendB == 0 {
		return addendA
	}
	return add(addendA ^ addendB, (addendA & addendB) << 1)
}

// sub: minuend - subtrahend = result
func sub(minuend, subtrahend int64) (result int64) {
	return add(minuend, add(^subtrahend, 1))
}

// multi: multiplicand * multiplier = product
func multi(multiplicand, multiplier int64) (product int64) {
	var shift int64
	for multiplier != 0 {
		if multiplier & 1 == 1 {
			product = add(product, multiplicand << shift)
		}
		shift = add(shift, 1)
		multiplier >>= 1
	}
	return
}

// div: dividend / divisor = quotient...remainder
func div(dividend, divisor int64) (quotient int64, remainder int64) {
	for i := int64(63); i >= 0; i = sub(i, 1) {
		if (dividend >> i) >= divisor {
			dividend = sub(dividend, divisor << i)
			quotient = add(quotient, 1 << i)
		}
	}
	remainder = dividend
	return
}