//**************************************************
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//**************************************************

// Most of the code is a translation to Golang from open office:
// https://github.com/apache/openoffice/blob/trunk/main/scaddins/source/analysis/financial.cxx

package financial

import (
	"errors"
	"math"
)

// Common errors.
var (
	ErrParametersError  = errors.New("invalid data, calculation failed due to missing parameters")
	ErrCalculationError = errors.New("calculation failed")
)

func getCumprinc(fRate float64, nNumPeriods int64, fVal float64, nStartPer int64, nEndPer int64, nPayType int64) (fKapZ float64, e error) {

	var fRmz float64

	if nStartPer < 1 || nEndPer < nStartPer || fRate <= 0.0 || nEndPer > nNumPeriods || nNumPeriods <= 0 || fVal <= 0.0 || (nPayType != 0 && nPayType != 1) {
		e = ErrParametersError
		return
	}
	fRmz = GetRmz(fRate, float64(nNumPeriods), fVal, 0.0, nPayType)

	fKapZ = 0.0

	nStart := uint64(nStartPer)
	nEnd := uint64(nEndPer)

	if nStart == 1 {
		if nPayType <= 0 {
			fKapZ = fRmz + fVal*fRate
		} else {
			fKapZ = fRmz
		}
		nStart++
	}

	for i := nStart; i <= nEnd; i++ {
		if nPayType > 0 {
			fKapZ += fRmz - (GetZw(fRate, float64(i-2), fRmz, fVal, 1)-fRmz)*fRate
		} else {
			fKapZ += fRmz - GetZw(fRate, float64(i-1), fRmz, fVal, 0)*fRate
		}
	}

	return fKapZ, nil
}

// ============================================================================
// XIRR helper functions

/** Calculates the resulting amount for the passed interest rate and the given XIRR parameters. */
func xirrResult(cf CashFlowTab, fRate float64) float64 {

	/*  V_0 ... V_n = input values.
	    D_0 ... D_n = input dates.
	    R           = input interest rate.

	    r   := R+1
	    E_i := (D_i-D_0) / 365

	                n    V_i                n    V_i
	    f(R)  =  SUM   -------  =  V_0 + SUM   ------- .
	               i=0  r^E_i              i=1  r^E_i
	*/

	D_0 := cf.FirstDate()
	r := fRate + 1.0
	fResult := cf.FirstFlow()

	for i, c := range cf {
		if i == 0 {
			continue // Skip first record
		}
		dayPerYear := float64(c.Date.DaysFrom(D_0)) / 365.0
		fResult += c.Flow / math.Pow(r, dayPerYear)

	}
	return fResult
}

/** Calculates the first derivation of lcl_sca_XirrResult(). */
func xirrResult_Deriv1(cf CashFlowTab, fRate float64) float64 {

	/*  V_0 ... V_n = input values.
	    D_0 ... D_n = input dates.
	    R           = input interest rate.

	    r   := R+1
	    E_i := (D_i-D_0) / 365

	                         n    V_i
	    f'(R)  =  [ V_0 + SUM   ------- ]'
	                        i=1  r^E_i

	                     n           V_i                 n    E_i V_i
	           =  0 + SUM   -E_i ----------- r'  =  - SUM   ----------- .
	                    i=1       r^(E_i+1)             i=1  r^(E_i+1)
	*/
	D_0 := cf.FirstDate()
	r := fRate + 1.0
	fResult := 0.0

	for i, c := range cf {
		if i == 0 {
			continue // Skip first record
		}
		dayPerYear := float64(c.Date.DaysFrom(D_0)) / 365.0
		fResult -= dayPerYear * c.Flow / math.Pow(r, dayPerYear+1.0)
	}
	return fResult
}

// Xirr calculation
func Xirr(cf CashFlowTab) (fResultRate float64, e error) {
	if len(cf) <= 2 {
		e = ErrParametersError
		return
	}

	// maximum epsilon for end of iteration
	const fMaxEps float64 = 1e-10
	// maximum number of iterations
	const nMaxIter int64 = 50

	// result interest rate, initialized with passed guessed rate, or 10%
	fResultRate = 0.1

	// Newton's method - try to find a fResultRate, so that lcl_sca_XirrResult() returns 0.
	var (
		fNewRate     float64
		fRateEps     float64
		fResultValue float64

		nIter              int64 = 0
		nIterScan          int64 = 0
		bContLoop          bool  = false
		bResultRateScanEnd bool  = false
	)

	// First the inner while-loop will be executed using the default Value fResultRate
	// or the user guessed fResultRate if those do not deliver a solution for the
	// Newton's method then the range from -0.99 to +0.99 will be scanned with a
	// step size of 0.01 to find fResultRate's value which can deliver a solution
	for {
		if nIterScan >= 1 {
			fResultRate = -0.99 + float64(nIterScan-1)*0.01
		}
		for {
			fResultValue = xirrResult(cf, fResultRate)
			fNewRate = fResultRate - fResultValue/xirrResult_Deriv1(cf, fResultRate)
			fRateEps = math.Abs(fNewRate - fResultRate)
			fResultRate = fNewRate
			bContLoop = (fRateEps > fMaxEps) && (math.Abs(fResultValue) > fMaxEps)
			nIter++
			if !(bContLoop && (nIter < nMaxIter)) {
				break
			}
		}
		nIter = 0
		if math.IsNaN(fResultRate) || math.IsInf(fResultRate, 0) || math.IsNaN(fResultValue) || math.IsInf(fResultValue, 0) {
			bContLoop = true
		}
		nIterScan++
		bResultRateScanEnd = nIterScan >= 200
		if !(bContLoop && !bResultRateScanEnd) {
			break
		}
	}
	if bContLoop {
		e = ErrCalculationError
		return
	}
	return
}

// Irr calculation IRR, neglecting dates.
func Irr(cf CashFlowTab) (guessRate float64, e error) {
	const (
		LowRate      float64 = 0.01
		HighRate     float64 = 0.5
		MaxIteration int     = 150
		PrecisionReq float64 = 1e-10
	)

	var (
		old           float64
		newVal        float64
		newGuessRate  = LowRate
		lowGuessRate  = LowRate
		highGuessRate = HighRate
		npv           float64
		denom         float64
	)
	numOfFlows := len(cf)
	guessRate = LowRate
	if numOfFlows < 2 {
		return 0, ErrParametersError
	}

	for i := 0; i < MaxIteration; i++ {
		npv = 0.00
		for j, c := range cf {
			denom = math.Pow(1+guessRate, float64(j))
			npv = npv + (c.Flow / denom)
		}
		/* Stop checking once the required precision is achieved */
		if (npv > 0) && (npv < PrecisionReq) {
			break
		}
		if old == 0 {
			old = npv
		} else {
			old = newVal
		}
		newVal = npv
		if i > 0 {
			if old < newVal {
				if old < 0 && newVal < 0 {
					highGuessRate = newGuessRate
				} else {
					lowGuessRate = newGuessRate
				}
			} else {
				if old > 0 && newVal > 0 {
					lowGuessRate = newGuessRate
				} else {
					highGuessRate = newGuessRate
				}
			}
		}
		guessRate = (lowGuessRate + highGuessRate) / 2
		newGuessRate = guessRate
	}
	if math.Abs(guessRate-npv) > PrecisionReq {
		e = ErrCalculationError // Even there is a calculation error, return the best guess.
	}
	return
}

func Xnpv(fRate float64, cf CashFlowTab) (fRet float64, e error) {
	nNum := len(cf)

	if nNum < 2 {
		e = ErrParametersError
		return
	}

	fRet = 0.0
	fNull := cf.FirstDate()
	fRate = fRate + 1

	for _, c := range cf {
		d := float64(c.Date.DaysFrom(fNull)) / 365.0
		fRet += c.Flow / math.Pow(fRate, d)
	}
	return
}

// Pv return Present Value
func Pv(rate float64, nper int64, pmt float64, fv float64, t bool) float64 {
	var tVal float64 = 0
	if t {
		tVal = 1
	}

	if rate == 0 {
		return -pmt*float64(nper) - fv
	}

	c := math.Pow(1+rate, float64(nper))
	return (((1-c)/rate)*pmt*(1+rate*tVal) - fv) / c
}

// Fv returns Future Value
func Fv(rate float64, nper int64, pmt float64, pv float64, t bool) float64 {
	var tVal float64 = 0
	if t {
		tVal = 1
	}
	if rate == 0 {
		return -pv - pmt*float64(nper)
	}

	c := math.Pow(1+rate, float64(nper))
	return -pv*c - pmt*(1+rate*tVal)*(c-1)/rate
}

func Xfv(rate float64, cf CashFlowTab) (fRet float64, e error) {
	nNum := len(cf)

	if nNum < 2 {
		e = ErrParametersError
		return
	}

	fRet = 0.0
	fNull := cf.FirstDate()
	rate = rate + 1

	for _, c := range cf {
		d := float64(c.Date.DaysFrom(fNull)) / 365.0
		fRet += c.Flow / math.Pow(rate, d)
	}
	return
}

func getCumipmt(fRate float64, nNumPeriods int64, fVal float64, nStartPer int64, nEndPer int64, nPayType int64) (fZinsZ float64, e error) {
	var (
		fRmz float64
	)

	if nStartPer < 1 || nEndPer < nStartPer || fRate <= 0.0 || nEndPer > nNumPeriods || nNumPeriods <= 0 || fVal <= 0.0 || (nPayType != 0 && nPayType != 1) {
		e = ErrParametersError
		return
	}

	fRmz = GetRmz(fRate, float64(nNumPeriods), fVal, 0.0, nPayType)

	fZinsZ = 0.0

	//sal_uInt32  nStart = sal_uInt32( nStartPer );
	//sal_uInt32  nEnd = sal_uInt32( nEndPer );

	if nStartPer == 1 {
		if nPayType <= 0 {
			fZinsZ = -fVal
		}

		nStartPer++
	}

	for i := nStartPer; i <= nEndPer; i++ {
		if nPayType > 0 {
			fZinsZ += GetZw(fRate, float64(i-2), fRmz, fVal, 1) - fRmz
		} else {
			fZinsZ += GetZw(fRate, float64(i-1), fRmz, fVal, 0)
		}
	}

	fZinsZ *= fRate

	return
}

func Pmt(rate float64, nper int64, pv float64, fv float64, t bool) float64 {
	if rate == 0 {
		return -((pv + fv) / float64(nper))
	}
	term := math.Pow(rate+1, float64(nper))
	if t {
		return -((fv*rate/(term-1) + pv*rate/(1-1/term)) / (1 + rate))
	}

	return -(fv*rate/(term-1) + pv*rate/(1-1/term))
}

// Crf is "Capital recovery factor"
func Crf(rate float64, nper int64) (result float64, e error) {
	if nper == 0 {
		e = ErrParametersError
		return
	}
	t1 := math.Pow(rate+1, float64(nper))
	result = (rate * t1) / (t1 - 1)
	return
}

// ReturnCoefficient is pv/pmt
func ReturnCoefficient(rate float64, nper int64, pv float64, fv float64, t bool) float64 {
	pmt := Pmt(rate/12, nper, pv, fv, t)
	if pmt == 0 {
		panic("pmt can't be zero, will lead to divide by zero")
	}
	return pv / pmt
}

func GetNominal(fRate float64, nPeriods int64) (float64, error) {
	if fRate <= 0.0 || nPeriods < 0 {
		return 0, ErrParametersError
	}

	fPeriods := float64(nPeriods)
	fRet := (math.Pow(fRate+1.0, 1.0/(fPeriods)) - 1.0) * (fPeriods)
	return fRet, nil
}

func GetEffectiveRate(fRate float64, nPeriods int64) (float64, error) {
	if fRate <= 0.0 || nPeriods < 0 {
		return 0, ErrParametersError
	}

	fPeriods := float64(nPeriods)
	fRet := math.Pow(fRate/fPeriods+1.0, fPeriods) - 1.0
	return fRet, nil
}
