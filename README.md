# go-finance-math
Excel like finance math in Go language

Most of the code is copy-paste-refactor from:
https://github.com/apache/openoffice/blob/trunk/main/scaddins/source/analysis/financial.cxx

## Goal
The goal of that repository is to be a solid base for larger projects where finance math is demand.
Currently, the code is ok, but still in development process.

## Structures
In that repository, 3 new structures:
* Date - Wrapper of time.Time, but used for date only. time is always zero.
* CashFlow - holds: date and flow. (flow is float64)
* CashFlowTab - is just slice of cashFlow
