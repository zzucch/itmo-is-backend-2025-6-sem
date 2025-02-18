package data

type Phone struct {
	Name        string
	Image       string
	Description string
	Price       string
}

type Review struct {
	Text   string
	Author string
}

type IndexPageData struct {
	Title          string
	SalePhone      Phone
	NewPhones      []Phone
	FeaturedPhones []Phone
}

var Data = IndexPageData{
	Title: "The Phone Marketplace",
	SalePhone: Phone{
		Name:  "Sharp Aquos wish4",
		Image: "/static/images/aquos-wish-4-1.jpg",
	},
	NewPhones: []Phone{
		{
			Name:        "Xiaomi Redmi 12 5G",
			Image:       "/static/images/redmi-12-5g.webp",
			Description: "Affordable swiftness",
		},
		{
			Name:        "Samsung Galaxy Z Fold6",
			Image:       "/static/images/galaxy-z-fold6.jpg",
			Description: "Unfold possibilities",
		},
		{
			Name:        "Sharp Aquos Wish 4",
			Image:       "/static/images/aquos-wish-4-2.jpg",
			Description: "Refreshingly simple",
		},
		{
			Name:        "Google Pixel 9",
			Image:       "/static/images/pixel-9.webp",
			Description: "Oxydizingly cutting-edge",
		},
	},
	FeaturedPhones: []Phone{
		{
			Name:        "Kyocera 902KC",
			Image:       "/static/images/digno-902kc.jpeg",
			Description: "A phone. Buy it.",
			Price:       "$399",
		},
		{
			Name:        "Sharp Aquos 601SH",
			Image:       "/static/images/aquos-601sh.jpg",
			Description: "A phone. Buy it.",
			Price:       "$499",
		},
		{
			Name:        "Siemens A50",
			Image:       "/static/images/siemens-a50.jpeg",
			Description: "A phone. Buy it.",
			Price:       "$499",
		},
		{
			Name:        "Kyocera KYF43",
			Image:       "/static/images/kyocera-kyf43.jpg",
			Description: "A phone. Buy it.",
			Price:       "$499",
		},
		{
			Name:        "Kyocera 414 Au Marvera",
			Image:       "/static/images/au-marvera-kyy08.jpg",
			Description: "A phone. Buy it.",
			Price:       "$499",
		},
		{
			Name:        "Blackberry Key2",
			Image:       "/static/images/blackberrry-key2.jpg",
			Description: "A phone. Buy it.",
			Price:       "$499",
		},
	},
}
