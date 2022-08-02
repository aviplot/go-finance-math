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

// Most of the code is a translation to Golang from:
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
func xirrResult(cf cashFlowTab, fRate float64) float64 {

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
func xirrResult_Deriv1(cf cashFlowTab, fRate float64) float64 {

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
func Xirr(cf cashFlowTab) (fResultRate float64, e error) {
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

func Xnpv(fRate float64, cf cashFlowTab) (fRet float64, e error) {
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

func Xfv(rate float64, cf cashFlowTab) (fRet float64, e error) {
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

//
//float64 getPrice( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fRate, float64 fYield, float64 fRedemp, int64 nFreq,
//const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fYield < 0.0 || fRate < 0.0 || fRedemp <= 0 || CHK_Freq || nSettle >= nMat )
//THROW_IAE;
//
//float64 fRet = getPrice_( GetNullDate( xOpt ), nSettle, nMat, fRate, fYield, fRedemp, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getPricedisc( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fDisc, float64 fRedemp, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fDisc <= 0.0 || fRedemp <= 0 || nSettle >= nMat )
//THROW_IAE;
//
//float64 fRet = fRedemp * ( 1.0 - fDisc * GetYearDiff( GetNullDate( xOpt ), nSettle, nMat, getDateMode( xOpt, rOB ) ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getPricemat( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nIssue, float64 fRate, float64 fYield, const ANY& rOB )
//THROWDEF_RTE_IAE
//{
//if( fRate < 0.0 || fYield < 0.0 || nSettle >= nMat )
//THROW_IAE;
//
//int64	nNullDate = GetNullDate( xOpt );
//int64   nBase = getDateMode( xOpt, rOB );
//
//float64		fIssMat = GetYearFrac( nNullDate, nIssue, nMat, nBase );
//float64		fIssSet = GetYearFrac( nNullDate, nIssue, nSettle, nBase );
//float64		fSetMat = GetYearFrac( nNullDate, nSettle, nMat, nBase );
//
//float64		fRet = 1.0 + fIssMat * fRate;
//fRet /= 1.0 + fSetMat * fYield;
//fRet -= fIssSet * fRate;
//fRet *= 100.0;
//
//RETURN_FINITE( fRet );
//}
//
//
//float64 getMduration( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fCoup, float64 fYield, int64 nFreq, const ANY& rOB )
//THROWDEF_RTE_IAE
//{
//if( fCoup < 0.0 || fYield < 0.0 || CHK_Freq )
//THROW_IAE;
//
//float64      fRet = GetDuration( GetNullDate( xOpt ),  nSettle, nMat, fCoup, fYield, nFreq, getDateMode( xOpt, rOB ) );
//fRet /= 1.0 + ( fYield / float64( nFreq ) );
//RETURN_FINITE( fRet );
//}
//
//

//
//
//float64 getDollarfr( float64 fDollarDec, int64 nFrac ) THROWDEF_RTE_IAE
//{
//if( nFrac <= 0 )
//THROW_IAE;
//
//float64	fInt;
//float64	fFrac = nFrac;
//
//float64	fRet = modf( fDollarDec, &fInt );
//
//fRet *= fFrac;
//
//fRet *= pow( 10.0, -ceil( log10( fFrac ) ) );
//
//fRet += fInt;
//
//RETURN_FINITE( fRet );
//}
//
//
//float64 getDollarde( float64 fDollarFrac, int64 nFrac ) THROWDEF_RTE_IAE
//{
//if( nFrac <= 0 )
//THROW_IAE;
//
//float64	fInt;
//float64	fFrac = nFrac;
//
//float64	fRet = modf( fDollarFrac, &fInt );
//
//fRet /= fFrac;
//
//fRet *= pow( 10.0, ceil( log10( fFrac ) ) );
//
//fRet += fInt;
//
//RETURN_FINITE( fRet );
//}
//
//
//float64 getYield( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fCoup, float64 fPrice, float64 fRedemp, int64 nFreq, const ANY& rOB )
//THROWDEF_RTE_IAE
//{
//if( fCoup < 0.0 || fPrice <= 0.0 || fRedemp <= 0.0 || CHK_Freq || nSettle >= nMat )
//THROW_IAE;
//
//float64 fRet = getYield_( GetNullDate( xOpt ), nSettle, nMat, fCoup, fPrice, fRedemp, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getYielddisc( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fPrice, float64 fRedemp, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fPrice <= 0.0 || fRedemp <= 0.0 || nSettle >= nMat )
//THROW_IAE;
//
//int64	nNullDate = GetNullDate( xOpt );
//
//#if 0
//float64		fRet = 1.0 - fPrice / fRedemp;
//fRet /= GetYearFrac( nNullDate, nSettle, nMat, getDateMode( xOpt, rOB ) );
//fRet /= 0.99795;  // don't know what this constant means in original
//#endif
//
//float64 fRet = ( fRedemp / fPrice ) - 1.0;
//fRet /= GetYearFrac( nNullDate, nSettle, nMat, getDateMode( xOpt, rOB ) );
//
//RETURN_FINITE( fRet );
//}
//
//
//float64 getYieldmat( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nIssue, float64 fRate, float64 fPrice, const ANY& rOB )
//THROWDEF_RTE_IAE
//{
//if( fRate < 0.0 || fRate <= 0.0 || nSettle >= nMat )
//THROW_IAE;
//
//float64 fRet = GetYieldmat( GetNullDate( xOpt ),  nSettle, nMat, nIssue, fRate, fPrice, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getTbilleq( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fDisc ) THROWDEF_RTE_IAE
//{
//nMat++;
//
//int64	nDiff = GetDiffDate360( xOpt, nSettle, nMat, sal_True );
//
//if( fDisc <= 0.0 || nSettle >= nMat || nDiff > 360 )
//THROW_IAE;
//
//float64 fRet = ( 365 * fDisc ) / ( 360 - ( fDisc * float64( nDiff ) ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getTbillprice( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fDisc ) THROWDEF_RTE_IAE
//{
//if( fDisc <= 0.0 || nSettle > nMat )
//THROW_IAE;
//
//nMat++;
//
//float64	fFraction = GetYearFrac( xOpt, nSettle, nMat, 0 );	// method: USA 30/360
//
//float64	fDummy;
//if( modf( fFraction, &fDummy ) == 0.0 )
//THROW_IAE;
//
//float64 fRet = 100.0 * ( 1.0 - fDisc * fFraction );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getTbillyield( constREFXPS& xOpt, int64 nSettle, int64 nMat, float64 fPrice )
//THROWDEF_RTE_IAE
//{
//int64	nDiff = GetDiffDate360( xOpt, nSettle, nMat, sal_True );
//nDiff++;
//
//if( fPrice <= 0.0 || nSettle >= nMat || nDiff > 360 )
//THROW_IAE;
//
//float64		fRet = 100.0;
//fRet /= fPrice;
//fRet--;
//fRet /= float64( nDiff );
//fRet *= 360.0;
//
//RETURN_FINITE( fRet );
//}
//
//
//float64 getOddfprice( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nIssue, int64 nFirstCoup,
//float64 fRate, float64 fYield, float64 fRedemp, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fRate < 0 || fYield < 0 || CHK_Freq || nMat <= nFirstCoup || nFirstCoup <= nSettle || nSettle <= nIssue )
//THROW_IAE;
//
//float64 fRet = GetOddfprice( GetNullDate( xOpt ), nSettle, nMat, nIssue, nFirstCoup, fRate, fYield, fRedemp, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getOddfyield( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nIssue, int64 nFirstCoup,
//float64 fRate, float64 fPrice, float64 fRedemp, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fRate < 0 || fPrice <= 0 || CHK_Freq || nMat <= nFirstCoup || nFirstCoup <= nSettle || nSettle <= nIssue )
//THROW_IAE;
//
//float64 fRet = GetOddfyield( GetNullDate( xOpt ), nSettle, nMat, nIssue, nFirstCoup, fRate, fPrice, fRedemp, nFreq,
//getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getOddlprice( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nLastInterest,
//float64 fRate, float64 fYield, float64 fRedemp, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fRate < 0 || fYield < 0 || CHK_Freq || nMat <= nSettle || nSettle <= nLastInterest )
//THROW_IAE;
//
//float64 fRet = GetOddlprice( GetNullDate( xOpt ), nSettle, nMat, nLastInterest, fRate, fYield, fRedemp, nFreq,
//getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getOddlyield( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nLastInterest,
//float64 fRate, float64 fPrice, float64 fRedemp, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fRate < 0 || fPrice <= 0 || CHK_Freq || nMat <= nSettle || nSettle <= nLastInterest )
//THROW_IAE;
//
//float64 fRet = GetOddlyield( GetNullDate( xOpt ), nSettle, nMat, nLastInterest, fRate, fPrice, fRedemp, nFreq,
//getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//// ============================================================================
//

//
//
//float64 getIntrate( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, float64 fInvest, float64 fRedemp, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//if( fInvest <= 0.0 || fRedemp <= 0.0 || nSettle >= nMat )
//THROW_IAE;
//
//float64 fRet = ( ( fRedemp / fInvest ) - 1.0 ) / GetYearDiff( GetNullDate( xOpt ), nSettle, nMat, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCoupncd( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCoupncd( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCoupdays( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCoupdays( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCoupdaysnc( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCoupdaysnc( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCoupdaybs( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCoupdaybs( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCouppcd( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCouppcd( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getCoupnum( constREFXPS& xOpt,
//int64 nSettle, int64 nMat, int64 nFreq, const ANY& rOB ) THROWDEF_RTE_IAE
//{
//float64 fRet = GetCoupnum( GetNullDate( xOpt ), nSettle, nMat, nFreq, getDateMode( xOpt, rOB ) );
//RETURN_FINITE( fRet );
//}
//
//
//float64 getFvschedule( float64 fPrinc, const SEQSEQ( float64 )& rSchedule ) THROWDEF_RTE_IAE
//{
//Scafloat64List aSchedList;
//
//aSchedList.Append( rSchedule );
//
//for( const float64* p = aSchedList.First() ; p ; p = aSchedList.Next() )
//fPrinc *= 1.0 + *p;
//
//RETURN_FINITE( fPrinc );
//}
//
