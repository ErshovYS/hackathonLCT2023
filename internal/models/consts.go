package models

const (
	Admin Role = iota
	Investor
)

const (
	TaxOCH Tax = iota
	TaxYCH
	TaxPatent
)

const (
	RegOOO Registration = iota
	RegIP
)

const (
	GovTaxOOO = 4000
	GovTaxIP  = 800

	PatentCoef = 0.06

	CapBuildingFrom   = 8000000
	CapBuildingTo     = 12000000
	CapRebuildingFrom = 500000
	CapRebuildingTo   = 1200000
)
