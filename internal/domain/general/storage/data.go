package storage

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
)

func (s *Storage) AddData() {
	salePhone := general.Phone{
		Name:        "Sharp Aquos wish4",
		Price:       100,
		Description: "Refreshingly simple",
		Brand:       "Sharp",
		Condition:   general.Excellent,
		Image: general.Image{
			URL: "/static/images/aquos-wish-4-1.jpg",
		},
	}
	s.DB.Create(&salePhone)

	newPhones := []general.Phone{
		{
			Name:        "Xiaomi Redmi 12 5G",
			Description: "Affordable swiftness",
			Price:       499.00,
			Brand:       "Xiaomi",
			Condition:   general.Good,
			Image: general.Image{
				URL: "/static/images/redmi-12-5g.webp",
			},
		},
		{
			Name:        "Samsung Galaxy Z Fold6",
			Description: "Unfold possibilities",
			Brand:       "Samsung",
			Price:       499.00,
			Condition:   general.Excellent,
			Image: general.Image{
				URL: "/static/images/galaxy-z-fold6.jpg",
			},
		},
		{
			Name:        "Sharp Aquos Wish 4",
			Description: "Refreshingly simple",
			Price:       499.00,
			Brand:       "Sharp",
			Condition:   general.Good,
			Image: general.Image{
				URL: "/static/images/aquos-wish-4-2.jpg",
			},
		},
		{
			Name:        "Google Pixel 9",
			Description: "Oxydizingly cutting-edge",
			Brand:       "Google",
			Price:       499.00,
			Condition:   general.Excellent,
			Image: general.Image{
				URL: "/static/images/pixel-9.webp",
			},
		},
	}
	for _, phone := range newPhones {
		s.DB.Create(&phone)
	}

	featuredPhones := []general.Phone{
		{
			Name:        "Kyocera 902KC",
			Description: "A phone. Buy it.",
			Brand:       "Kyocera",
			Condition:   general.Fair,
			Price:       399.00,
			Image: general.Image{
				URL: "/static/images/digno-902kc.jpeg",
			},
		},
		{
			Name:        "Sharp Aquos 601SH",
			Description: "A phone. Buy it.",
			Brand:       "Sharp",
			Condition:   general.Good,
			Price:       499.00,
			Image: general.Image{
				URL: "/static/images/aquos-601sh.jpg",
			},
		},
		{
			Name:        "Siemens A50",
			Description: "A phone. Buy it.",
			Brand:       "Siemens",
			Condition:   general.Fair,
			Price:       499.00,
			Image: general.Image{
				URL: "/static/images/siemens-a50.jpeg",
			},
		},
		{
			Name:        "Kyocera KYF43",
			Description: "A phone. Buy it.",
			Brand:       "Kyocera",
			Condition:   general.Good,
			Price:       499.00,
			Image: general.Image{
				URL: "/static/images/kyocera-kyf43.jpg",
			},
		},
		{
			Name:        "Kyocera 414 Au Marvera",
			Description: "A phone. Buy it.",
			Brand:       "Kyocera",
			Condition:   general.Fair,
			Price:       499.00,
			Image: general.Image{
				URL: "/static/images/au-marvera-kyy08.jpg",
			},
		},
		{
			Name:        "Blackberry Key2",
			Description: "A phone. Buy it.",
			Brand:       "Blackberry",
			Condition:   general.Good,
			Price:       499.00,
			Image: general.Image{
				URL: "/static/images/blackberrry-key2.jpg",
			},
		},
	}
	for _, phone := range featuredPhones {
		s.DB.Create(&phone)
	}

	saleCatalog := catalog.Catalog{
		Name:   "Sale Phone",
		Phones: []general.Phone{salePhone},
	}
	s.DB.Create(&saleCatalog)

	newPhonesCatalog := catalog.Catalog{
		Name:   "New Phones",
		Phones: newPhones,
	}
	s.DB.Create(&newPhonesCatalog)

	featuredPhonesCatalog := catalog.Catalog{
		Name:   "Featured Phones",
		Phones: featuredPhones,
	}
	s.DB.Create(&featuredPhonesCatalog)
}
