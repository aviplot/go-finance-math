# go-finance-math
Excel like finance math in Go language

Most of the code is copy-paste-refactor from:
https://github.com/apache/openoffice/blob/trunk/main/scaddins/source/analysis/financial.cxx

## Goal
The goal of that repository is to be a solid base for larger projects where finance math is demand.

## Structures
In that repository, 3 new structures:
* date - Wrapper of time.Time, but used for date only. time is always zero.
* cashFlow - holds: date and flow. (flow is float64)
* cashFlowTab - is just slice of cashFlow
