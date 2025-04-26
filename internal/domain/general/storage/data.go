package storage

import (
	"time"

	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/catalog"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/user"
	"gorm.io/gorm"
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
		SellerID: 1,
	}
	s.DB.Create(&salePhone)

	userQwe := user.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "qwe",
		Email:    "q@we",
		PasswordHash: func() string {
			hash, _ := user.HashPassword("qwe")
			return hash
		}(),
		LastLogin: time.Now(),
	}
	s.DB.Create(&userQwe)

	userAdmin := user.User{
		Model: gorm.Model{
			ID: 2,
		},
		Username: "ewq",
		Email:    "e@wq",
		PasswordHash: func() string {
			hash, _ := user.HashPassword("qwe")
			return hash
		}(),
		LastLogin: time.Now(),
		Role:      user.RoleAdmin,
	}
	s.DB.Create(&userAdmin)

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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
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
			SellerID: userAdmin.ID,
		},
	}
	for _, phone := range featuredPhones {
		s.DB.Create(&phone)
	}

	saleCatalog := catalog.Catalog{
		Name:      "Sale Phone",
		Phones:    []general.Phone{salePhone},
		CreatorID: 2,
	}
	s.DB.Create(&saleCatalog)

	newPhonesCatalog := catalog.Catalog{
		Name:      "New Phones",
		Phones:    newPhones,
		CreatorID: 2,
	}
	s.DB.Create(&newPhonesCatalog)

	featuredPhonesCatalog := catalog.Catalog{
		Name:      "Featured Phones",
		Phones:    featuredPhones,
		CreatorID: 2,
	}
	s.DB.Create(&featuredPhonesCatalog)
}
