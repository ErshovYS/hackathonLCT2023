package models

import (
	"gorm.io/gorm"
)

type Role uint
type Tax uint
type Registration uint

type User struct {
	gorm.Model
	FirstName    string
	MiddleName   string
	LastName     string
	Email        string
	Organization string
	INN          string
	Site         string
	IndustryID   uint
	Country      string
	City         string
	Job          string
	Password     string
	Token        string
	Role         Role
	Calculations []Calculation
}

type Industry struct {
	gorm.Model
	Name         string
	Workers      uint32
	Salary       uint32
	MoscowTax    uint32
	ProfitTax    uint32
	PropertyTax  uint32
	EstateTax    uint32
	PersonalTax  uint32
	TransportTax uint32
	OtherTax     uint32
	Users        []User
	Calculations []Calculation
}

type District struct {
	gorm.Model
	Name         string
	Price        uint32
	Calculations []Calculation
}

type EquipmentList struct {
	gorm.Model
	Name     string
	PriceUSD uint32
	PriceRUB uint32
}

type Equipment struct {
	gorm.Model
	Name     string
	PriceRUB uint32
	Count    uint32
}

type RegistrationTax struct {
	gorm.Model
	Registration Registration
	Tax          Tax
	From         uint32
	To           uint32
	Fee          uint32
}

type Calculation struct {
	gorm.Model
	UserID            uint
	IndustryID        *uint
	WorkerCount       uint32
	DistrictID        uint
	LandArea          float32
	CapRebuildingArea float32
	CapBuildingArea   float32
	Equipments        []Equipment `gorm:"many2many:calculation_equipment;"`
	Buildings         []Building  `gorm:"many2many:calculation_building;"`
	RegistrationTaxID uint
	PatentID          *uint
	OtherPayments     uint64
	PersonalFrom      float32
	PersonalTo        float32
	EstateFrom        float32
	EstateTo          float32
	TaxFrom           float32
	TaxTo             float32
	ServiceFrom       float32
	ServiceTo         float32
	ResultFrom        float64
	ResultTo          float64
	ReportLink        string
}

type Building struct {
	gorm.Model
	Name string
	Area float32
}

type Patent struct {
	gorm.Model
	Name            string
	PotencialProfit uint32
	Tax             uint32
	Price           uint32
	Calculations    []Calculation
}
