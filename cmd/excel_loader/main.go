package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"invest/internal/config"
	"invest/internal/models"
	"log"
	"math"
	"strconv"
	"strings"
)

type industry struct {
	workers, salary, moscowTax, profitTax float64
	PropertyTax                           float64
	EstateTax                             float64
	PersonalTax                           float64
	TransportTax                          float64
	OtherTax                              float64
}

func main() {
	var cfg config.Config
	err := envconfig.Process("invest", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	var dbDial gorm.Dialector
	switch cfg.DB.Type {
	case "sqlite":
		dbDial = sqlite.Open(cfg.DB.DSN)
	case "mysql":
		dbDial = mysql.Open(cfg.DB.DSN)
	case "postgres":
		dbDial = postgres.Open(cfg.DB.DSN)
	default:
		log.Fatal("wrong type of database")
	}
	db, err := gorm.Open(dbDial)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Industry{}, &models.Patent{}, &models.RegistrationTax{},
		&models.Building{}, &models.Equipment{}, &models.Calculation{}, &models.EquipmentList{}, &models.District{})
	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Transaction(func(tx *gorm.DB) error {
		// Create admin
		tx.Create(&models.User{Fullname: "Administrator", Password: "1", Token: "1", Role: models.Admin})
		if tx.Error != nil {
			return tx.Error
		}
		tx.Create(&models.User{Fullname: "Investor", Password: "2", Token: "2", Role: models.Investor})
		if tx.Error != nil {
			return tx.Error
		}

		var rows [][]string
		// Create a registration taxes
		rts := make([]models.RegistrationTax, 0, 5)
		rows, err = getRows("/home/lenovo/Downloads/Датасеты промышленность/Датасеты/Доп.услуги хакатон.xlsx", "Бух.учет")
		if err != nil {
			return err
		}
		for _, r := range rows {
			switch r[0] {
			case "ООО (АО)":
				from, to := parseBuh(r[1])
				rts = append(rts, models.RegistrationTax{Registration: models.RegOOO, Tax: models.TaxOCH, From: from, To: to, Fee: models.GovTaxOOO})
				from, to = parseBuh(r[2])
				rts = append(rts, models.RegistrationTax{Registration: models.RegOOO, Tax: models.TaxYCH, From: from, To: to, Fee: models.GovTaxOOO})
			case "ИП":
				from, to := parseBuh(r[1])
				rts = append(rts, models.RegistrationTax{Registration: models.RegIP, Tax: models.TaxOCH, From: from, To: to, Fee: models.GovTaxIP})
				from, to = parseBuh(r[2])
				rts = append(rts, models.RegistrationTax{Registration: models.RegIP, Tax: models.TaxYCH, From: from, To: to, Fee: models.GovTaxIP})
				from, to = parseBuh(r[3])
				rts = append(rts, models.RegistrationTax{Registration: models.RegIP, Tax: models.TaxPatent, From: from, To: to, Fee: models.GovTaxIP})
			default:
				continue
			}
		}
		tx.Create(&rts)
		if tx.Error != nil {
			return tx.Error
		}

		// Create a districts
		districts := make([]models.District, 0, 12)
		rows, err = getRows("/home/lenovo/Downloads/Датасеты промышленность/Датасеты/Средняя кадастр стоимость по округам.xlsx", "Среднее по округам")
		if err != nil {
			return err
		}
		for i, r := range rows {
			if i == 0 {
				continue
			}
			price, err := strconv.ParseFloat(r[2], 64)
			if err != nil {
				return err
			}
			districts = append(districts, models.District{Name: r[1], Price: uint32(math.Round(price * 100))})
		}
		tx.Create(&districts)
		if tx.Error != nil {
			return tx.Error
		}

		// Create equipment
		equips := make([]models.EquipmentList, 0, 9)
		rows, err = getRows("/home/lenovo/Downloads/Датасеты промышленность/Датасеты/Станки средняя цена.xlsx", "расчет средней цены")
		if err != nil {
			return err
		}
		for i, r := range rows {
			if i == 0 || len(r) < 4 || r[2] == "" || r[3] == "" {
				continue
			}
			priceUSD, err := strconv.ParseFloat(r[2], 32)
			if err != nil {
				return err
			}
			priceRUB, err := strconv.ParseFloat(r[3], 32)
			if err != nil {
				return err
			}
			equips = append(equips, models.EquipmentList{Name: r[1], PriceUSD: uint32(math.Round(priceUSD * 100)), PriceRUB: uint32(math.Round(priceRUB * 100))})
		}
		tx.Create(&equips)
		if tx.Error != nil {
			return tx.Error
		}

		// Create industry
		inds := make(map[string][]industry)
		rows, err = getRows("/home/lenovo/Downloads/Датасеты промышленность/Датасеты/Обезличенные данные.xlsm", "лист2")
		if err != nil {
			return err
		}
		for i, r := range rows {
			if i == 0 || len(r) < 19 {
				continue
			}
			ind := industry{}
			buff, err := strconv.ParseFloat(strings.Replace(r[3], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.workers = buff
			buff, err = strconv.ParseFloat(strings.Replace(r[5], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.salary = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[6], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.moscowTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[8], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.profitTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[10], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.PropertyTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[12], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.EstateTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[14], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.PersonalTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[16], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.TransportTax = buff * 100000
			buff, err = strconv.ParseFloat(strings.Replace(r[18], "-", "", -1), 64)
			if err != nil {
				continue
			}
			ind.OtherTax = buff * 100000
			inds[r[0]] = append(inds[r[0]], ind)
		}

		industries := make([]models.Industry, 0, len(inds))
		for k, v := range inds {
			var workers, salary, moscowTax, profitTax, propertyTax, estateTax, personalTax, transportTax, otherTax float64
			for _, ind := range v {
				workers += ind.workers
				salary += ind.salary
				moscowTax += ind.moscowTax
				profitTax += ind.profitTax
				propertyTax += ind.PropertyTax
				estateTax += ind.EstateTax
				personalTax += ind.PersonalTax
				transportTax += ind.TransportTax
				otherTax += ind.OtherTax
			}
			industries = append(industries, models.Industry{
				Name:         k,
				Workers:      uint32(workers / float64(len(v))),
				Salary:       uint32(salary / float64(len(v))),
				MoscowTax:    uint32(moscowTax / float64(len(v))),
				ProfitTax:    uint32(profitTax / float64(len(v))),
				PropertyTax:  uint32(propertyTax / float64(len(v))),
				EstateTax:    uint32(estateTax / float64(len(v))),
				PersonalTax:  uint32(personalTax / float64(len(v))),
				TransportTax: uint32(transportTax / float64(len(v))),
				OtherTax:     uint32(otherTax / float64(len(v))),
			})
		}
		tx.Create(&industries)
		if tx.Error != nil {
			return tx.Error
		}

		// Create patent
		patents := make([]models.Patent, 0, 9)
		rows, err = getRows("/home/lenovo/Downloads/Датасеты промышленность/Датасеты/Патентование (потенциальный доход, Москва).xlsx", "Table 1")
		if err != nil {
			return err
		}
		for i, r := range rows {
			if i < 2 {
				continue
			}
			potencialProfit, err := strconv.ParseFloat(r[2], 32)
			if err != nil {
				return err
			}
			tax, err := strconv.ParseFloat(r[3], 32)
			if err != nil {
				return err
			}
			price, err := strconv.ParseFloat(r[4], 32)
			if err != nil {
				return err
			}
			patents = append(patents, models.Patent{Name: r[1], PotencialProfit: uint32(math.Round(potencialProfit * 100)),
				Tax: uint32(tax * 100), Price: uint32(math.Round(price * 100))})
		}
		tx.Create(&patents)
		if tx.Error != nil {
			return tx.Error
		}

		return nil
	}); err != nil {
		log.Fatal(err.Error())
	}
}

func getRows(filename, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatalf("close excel file: %v", err)
		}
	}()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func parseBuh(data string) (from, to uint32) {
	data = strings.Replace(data, "от ", "", -1)
	vals := strings.Split(data, " до ")
	mins := strings.Split(vals[0], "-")
	maxs := strings.Split(vals[1], "-")
	buff1, _ := strconv.ParseUint(mins[0], 10, 32)
	buff2, _ := strconv.ParseUint(mins[1], 10, 32)
	from = uint32((buff1 + buff2) / 2)
	buff1, _ = strconv.ParseUint(maxs[0], 10, 32)
	buff2, _ = strconv.ParseUint(maxs[1], 10, 32)
	to = uint32((buff1 + buff2) / 2)
	return
}
