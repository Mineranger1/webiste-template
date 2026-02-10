package database

import (
	"fmt"
	"log"
)

type SeedProduct struct {
	Name        string
	Description string
	Dosage      string
	Data        string
	Category    string
}

func Seed() error {
	// Check if data exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	categories := map[string]string{
		"krowy":         "Krowy Zasuszone i w Laktacji",
		"bydlo-opasowe": "Bydło Opasowe",
		"dodatki":       "Dodatki Paszowe i Specjalistyczne",
		"cieleta":       "Cielęta",
		"specjalne":     "Produkty Specjalne - Higiena i Tłuszcze",
	}

	categoryIDs := make(map[string]int)

	for slug, name := range categories {
		res, err := DB.Exec("INSERT INTO categories (name, slug) VALUES (?, ?)", name, slug)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", name, err)
		}
		id, _ := res.LastInsertId()
		categoryIDs[slug] = int(id)
	}

	products := []SeedProduct{
		// Krowy
		{
			Name:        "Biomix DRY (Krowy Zasuszone)",
			Category:    "krowy",
			Description: "Mieszanka mineralno - witaminowa o odpowiednich proporcjach wapnia, fosforu i magnezu oraz wysokiej dawce witaminy E. Zapewnia bezpieczny okres zasuszenia i świetny start w laktację.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Beta Forte (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka zapewniająca wysokie wskaźniki rozrodu wysokowydajnych krów. Poprawia gospodarkę hormonalną. Dodatek chronionej choliny dodatkowo zabezpiecza wątrobę przed otłuszczeniem.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Energy Drink (Krowy Zasuszone)",
			Category:    "krowy",
			Description: "Mieszanka stymuluje apetyt zaraz po wycieleniu. Uzupełnia gospodarkę wodną i elektrolitową. Zapobiega niedoborom energii i występowaniu porażeń poporodowych. Stymuluje rozwój brodawek żwacza.",
			Dosage:      "Według etykiety",
		},
		{
			Name:        "Biomix Kation / Anion (Krowy Zasuszone)",
			Category:    "krowy",
			Description: "Sole gorzkie w formie chronionej zapewniają stymulację wchłaniania wapnia przez krowy 2-3 tygodnie przed porodem. Zapobiega to zaleganiu poporodowemu i umożliwia zdrową kolejną laktację.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Beta (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka poprawiająca wskaźniki rozrodu bydła. Wysoki poziom witaminy A, D3, E oraz obecność beta karotenu i selen wspiera funkcje rozrodcze organizmu. Wpływa też na zawartość przeciwciał w siarze.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Somatic (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka podnosi odporność przez obecność witaminy E oraz cynku. Efektywnie obniża poziom komórek somatycznych w mleku. Zabezpiecza wątrobę, a dawka biotyny redukuje występowanie kulawizn.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Lacto (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka z wysoką zawartością magnezu zapobiega tężyczce pastwiskowej. Uzupełnia dawkę pokarmową w niezbędne mikro i makro składniki w funkcjonowaniu zwierząt laktacyjnych.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Extra KM (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka z witaminą E oraz wysoko przyswajalnymi formami cynku, miedzi, manganu i selenu. Kwas foliowy i niacyna wpływają na metabolizm krów, a wysoki poziom biotyny wspomaga racice.",
			Dosage:      "200-250g dziennie",
		},
		{
			Name:        "Biomix Racice (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka z najwyższym poziomem biotyny oraz wysokim poziomem witaminy E i łatwo przyswajalnym chelacie cynku, miedzi i manganu - działa wspomagająco na racice, skórę i sierść.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix HP Active Max (Krowy w Laktacji)",
			Category:    "krowy",
			Description: "Mieszanka z żywymi kulturami drożdży zapewnia wysoką strawność dawki pokarmowej. Stabilizuje pH żwacza, a dodatek niacyny, biotyny oraz kwasu foliowego stymulują metabolizm do wysokiej i bezpiecznej produkcji mleka.",
			Dosage:      "200-250g dziennie",
		},
		{
			Name:        "BIOMIX ROBOT (Pasze)",
			Category:    "krowy",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Stanowi uzupełnienie dawki pasz objętościowych o średniej zawartości energii i białka.",
			Data:        "Białko surowe 19.2% / Tłuszcz surowy 3.2% / Energia 7.7MJ",
			Dosage:      "Krowy mleczne 1-10kg/szt./dzień",
		},
		{
			Name:        "BIOMIX ENERGYFAT (Pasze)",
			Category:    "krowy",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Przeznaczona do pokrycia niedoborów energetycznych w szczycie laktacji. Zawiera metabolity drożdży oraz tłuszcz chroniony.",
			Data:        "Białko surowe 19% / Tłuszcz surowy 3.4% / Energia 8.0MJ",
			Dosage:      "Krowy mleczne 1-4kg/szt./dzień",
		},
		{
			Name:        "BIOMIX BIOMLEK 24 (Pasze)",
			Category:    "krowy",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Mieszanka stanowi uzupełnienie dawki pasz objętościowych o niskiej zawartości białka.",
			Data:        "Białko surowe 24% / Tłuszcz surowy 3.8% / Energia 7.1MJ",
			Dosage:      "Krowy mleczne 1-10kg/szt./dzień",
		},
		// Opasy
		{
			Name:        "Biomix Super Opas",
			Category:    "bydlo-opasowe",
			Description: "Mieszanka wspomagająca tucz opasów. Siarka umożliwia efektywne wykorzystanie składników dawki. Dodatek drożdży poprawia strawność, a zastosowanie witamin w formach chelatowych zwiększa ich przyswajalność. Dla alkalizowanego zboża.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Opas Max",
			Category:    "bydlo-opasowe",
			Description: "Mieszanka wspomagająca intensywny tucz dla opasów. Siarka pozwala na efektywne wykorzystanie składników dawki. Obecność drożdży wspomaga strawność i wykorzystanie dawki pokarmowej.",
			Dosage:      "100-150g dziennie",
		},
		// Dodatki
		{
			Name:        "Biomix Drożdże MAX",
			Category:    "dodatki",
			Description: "Mieszanka żywych drożdży oraz metabolitów drożdżowych. Produkt został uzupełniony drożdżami piwnymi oraz (MOS), które wspomagają mikroflorę przewodu pokarmowego i zapewniają ochronę jelit.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Drożdże",
			Category:    "dodatki",
			Description: "Mieszanka z żywymi kulturami drożdży. Zapewnia prawidłową fermentację żwacza. Zwiększa produkcję mleka. Poprawia stan zdrowotny stada i zabezpiecza krowy przed stresem cieplnym.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Bufor Alg",
			Category:    "dodatki",
			Description: "Kompleks buforujący o połączonym działaniu alg morskich, tlenku magnezu i kwaśnego węglanu. Zapobiega kwasicy przez odpowiednią stabilizację pH żwacza. Zalecany podczas stresu cieplnego.",
			Dosage:      "100-200g dziennie",
		},
		{
			Name:        "Biomix Immuno",
			Category:    "dodatki",
			Description: "Mieszanka dla krów w okresie zmniejszonej odporności redukująca skutki stresu cieplnego. Dzięki wysokiej dawce witaminy C i E znacząco wpływa na wzrost układu immunologicznego redukując wpływ patogenów.",
			Dosage:      "10g dziennie co drugi tydzień",
		},
		{
			Name:        "Biomix Rumen Gold",
			Category:    "dodatki",
			Description: "Kompleks buforujący, ale z dodatkiem żywych drożdży i fosforanu jednowapniowego. Zapobiega kwasicy, utrzymując odpowiednie pH w żwaczu. Zwiększa skuteczność trawienia.",
			Dosage:      "300-400g dziennie",
		},
		{
			Name:        "Biomix Sorbent Tox",
			Category:    "dodatki",
			Description: "Mieszanka paszowa eliminująca toksyny. Redukuje mikotoksyny i endotoksyny występujące w paszach i przewodzie pokarmowym. Kompleks składników potęguje skuteczność preparatu.",
			Dosage:      "15-50g dziennie",
		},
		{
			Name:        "Biomix Sorbacid",
			Category:    "dodatki",
			Description: "Mieszanka ograniczająca fermentację oraz zabezpiecza kiszonki po otwarciu pryzmy przed powstawaniem mykotoksyn. TMR/PMR skutecznie jest chroniony przed zagrzaniem się na stole paszowym.",
			Dosage:      "1-2 kg/ tonę TMR/PMR",
		},
		{
			Name:        "Biomix Soda Lux",
			Category:    "dodatki",
			Description: "Preparat buforujący o połączonym działaniu. Skutecznie zapobiega kwasicy żwacza oraz stabilizuje pH przewodu pokarmowego, wspierając prawidłowe trawienie i zdrowie zwierząt.",
			Dosage:      "300g dziennie",
		},
		{
			Name:        "Biomix Inokulant (Zakiszacz)",
			Category:    "dodatki",
			Description: "Kompleks mikroorganizmów, który zawiera kultury bakterii kwasu propionowego oraz mlekowego. Kwasy te hamują procesy gnilne, ograniczają grzanie się kiszonki i znacznie wydłużają ich trwałość.",
			Dosage:      "Według etykiety",
		},
		// Cielęta
		{
			Name:        "Biomix Calf",
			Category:    "cieleta",
			Description: "Mieszanka dla cieląt i młodzieży. Gwarantuje prawidłowy rozwój układu rozrodczego i hormonalnego. Wspomaga właściwy rozwój gruczołów płciowych. U zacielonych jałówek wspiera utrzymanie płodu.",
			Dosage:      "Zależne od wieku",
		},
		{
			Name:        "BIOMIX STARTER CALF (Pasza)",
			Category:    "cieleta",
			Description: "Mieszanka paszowa uzupełniająca dla cieląt. Świetny start to świetna produkcja potem. Zadbaj o swoje cielęta najlepiej jak możesz.",
			Data:        "Białko surowe 18% / Tłuszcz surowy 3.4%",
			Dosage:      "Cielęta 2tyg. - 6 miesięcy / 1-3kg/szt./dzień",
		},
		{
			Name:        "BioMilk GOLD (Preparat mlekozastępczy)",
			Category:    "cieleta",
			Description: "Oparty w 30% na odtłuszczonym mleku w proszku, serwatce oraz tłuszczach roślinnych. Innowacyjna receptura pozwala na zwiększenie rozpuszczalności oraz strawności produktu.",
			Data:        "0.5% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},
		{
			Name:        "BioMilk SILVER (Preparat mlekozastępczy)",
			Category:    "cieleta",
			Description: "Oparty o serwatkę w proszku, tłuszcze roślinne oraz odtłuszczone mleko w proszku z dodatkiem drożdży.",
			Data:        "0.9% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},
		{
			Name:        "BioMilk PLATINUM (Preparat mlekozastępczy)",
			Category:    "cieleta",
			Description: "Najwyższej jakości preparat, oparty w 40% na odtłuszczonym mleku w proszku oraz tłuszczach roślinnych. Innowacyjna metoda produkcji pozwala zwiększyć rozpuszczalność i strawność produktu.",
			Data:        "0% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},
		{
			Name:        "BIOMIX TMR FULL (Pasza startowa)",
			Category:    "cieleta",
			Description: "Kompletna pasza startowa dla cieląt. Zapewnia stabilne żywienie przez cały okres odchowu, eliminując stres związany ze zmianą pasz. Wspiera prawidłowy rozwój przewodu pokarmowego oraz mikroflory jelitowej cieląt. Poprawia pobranie paszy, wspomaga rozwój żwacza i umożliwia wcześniejsze ograniczenie odpajania mlekiem.",
		},
		// Specjalne
		{
			Name:        "BIOMIX DEZODRY (Dezynfekcja)",
			Category:    "specjalne",
			Description: "Preparat do suchej dezynfekcji legowisk, ściółki oraz pomieszczeń inwentarskich. Skutecznie ogranicza rozwój bakterii, grzybów i drobnoustrojów chorobotwórczych. Wiąże nadmiar wilgoci, redukuje nieprzyjemne zapachy oraz obniża poziom amoniaku. Zmniejsza ryzyko chorób racic, skóry i układu oddechowego.",
		},
		{
			Name:        "SUPER FAT (Suplement tłuszczowy)",
			Category:    "specjalne",
			Description: "Sprawdzony suplement tłuszczowy dla wysokowydajnych krów, który zwiększa energetyczność dawki. Dzięki unikalnej kombinacji kwasów tłuszczowych skutecznie zaspokaja zapotrzebowanie energetyczne. Kluczowe cechy: - Zmniejsza ryzyko stłuszczenia wątroby. - Zawiera kwas oleinowy w celu poprawy strawności. - Dostarcza nienasycone kwasy tłuszczowe i kwasy Omega 3. - Wspomaga układ odpornościowy i płodność. - Wysoka strawność, brak wpływu na pH żwacza (w przeciwieństwie do mydeł wapniowych).",
		},
	}

	for _, p := range products {
		catID := categoryIDs[p.Category]
		_, err := DB.Exec(`INSERT INTO products (name, description, dosage, data, category_id) VALUES (?, ?, ?, ?, ?)`,
			p.Name, p.Description, p.Dosage, p.Data, catID)
		if err != nil {
			log.Printf("Failed to insert product %s: %v", p.Name, err)
		}
	}

	return nil
}
