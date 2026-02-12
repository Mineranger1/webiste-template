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
	Content     string
	ImagePath   string
}

func Seed() error {
	if err := seedProducts(); err != nil {
		return err
	}
	return seedEmployees()
}

func seedEmployees() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	employees := []struct {
		Name      string
		Position  string
		Phone     string
		Email     string
		ImagePath string
	}{
		{
			Name:      "Wojciech Szwarc",
			Position:  "Prezes // Doradca Żywieniowy",
			Phone:     "+48 603 789 661",
			Email:     "wojciech.szwarc@biomixpoland.pl",
			ImagePath: "/static/images/wojciech_szwarc.webp",
		},
		{
			Name:      "Adrian Adamczak",
			Position:  "Dyrektor Handlowy ds. Bydła // Doradca Żywieniowy",
			Phone:     "+48 787 589 927",
			Email:     "adrian.adamczak@biomixpoland.pl",
			ImagePath: "/static/images/adrian_adamczak.webp",
		},
		{
			Name:      "Łukasz Zamojdzin",
			Position:  "Specjalista ds. żywienia",
			Phone:     "+48 884 259 025",
			Email:     "lukasz.zamojdzin@biomixpoland.pl",
			ImagePath: "/static/images/lukasz_zamojdzin.webp",
		},

		{
			Name:      "Jarosław Długołęcki",
			Position:  "Specjalista ds. żywienia",
			Phone:     "+48 609 095 318",
			Email:     "jaroslaw.dlugolecki@biomixpoland.pl",
			ImagePath: "/static/images/jaroslaw_dlugolecki.webp",
		},
		{
			Name:      "Ewelina Nowicka",
			Position:  "Specjalista ds. żywienia",
			Phone:     "+48 665 003 160",
			Email:     "ewelina.nowicka@biomixpoland.pl",
			ImagePath: "/static/images/wojciech_szwarc.webp",
		},
		{
			Name:      "Sandra Hemerling",
			Position:  "Kierownik Zakładu",
			Phone:     "+48 661 017 021",
			Email:     "sandra.hemerling@biomixpoland.pl",
			ImagePath: "/static/images/wojciech_szwarc.webp",
		},
		{
			Name:      "Mirosława Kamińska",
			Position:  "Specjalista ds. obsługi klienta i jakości",
			Phone:     "+48 661 026 133",
			Email:     "miroslawa.kaminska@biomixpoland.pl",
			ImagePath: "/static/images/wojciech_szwarc.webp",
		},
	}

	for _, e := range employees {
		_, err := DB.Exec("INSERT INTO employees (name, position, phone, email, image_path) VALUES (?, ?, ?, ?, ?)",
			e.Name, e.Position, e.Phone, e.Email, e.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to insert employee %s: %w", e.Name, err)
		}
	}
	return nil
}

func seedProducts() error {
	// Check if data exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		// If table doesn't exist, InitDB should have created it, but let's be safe
		return err
	}
	if count > 0 {
		return nil // Already seeded
	}

	// 1. Categories
	categories := []struct {
		Name        string
		Slug        string
		Description string
		ImagePath   string
	}{
		{
			Name:        "Premixy",
			Slug:        "premixy",
			Description: "Szeroka gama mieszanek mineralno-witaminowych dla bydła w każdym wieku.",
			ImagePath:   "/static/images/premixy.webp",
		},
		{
			Name:        "Preparaty mlekozastępcze",
			Slug:        "preparaty-mlekozastepcze",
			Description: "Wysokiej jakości preparaty dla cieląt zapewniające zdrowy start.",
			ImagePath:   "/static/images/mlekozastepcze.webp",
		},
		{
			Name:        "Pasze",
			Slug:        "pasze",
			Description: "Pełnoporcjowe i uzupełniające mieszanki paszowe.",
			ImagePath:   "/static/images/pasze.webp",
		},
		{
			Name:        "Specjalne",
			Slug:        "specjalne",
			Description: "Produkty specjalistyczne do zadań specjalnych, higieny i suplementacji.",
			ImagePath:   "/static/images/specjalne.webp",
		},
	}

	categoryIDs := make(map[string]int)

	for _, cat := range categories {
		res, err := DB.Exec("INSERT INTO categories (name, slug, description, image_path) VALUES (?, ?, ?, ?)",
			cat.Name, cat.Slug, cat.Description, cat.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", cat.Name, err)
		}
		id, _ := res.LastInsertId()
		categoryIDs[cat.Slug] = int(id)
	}

	// 2. Products
	// Content for SuperFat
	superFatContent := `
<div class="space-y-6">
    <h2 class="text-2xl font-bold text-primary">SUPER FAT - Suplement tłuszczowy</h2>
    <p class="text-lg">Sprawdzony suplement tłuszczowy dla wysokowydajnych krów, który zwiększa energetyczność dawki, pomagając ograniczyć utratę masy ciała, poprawić parametry mleka oraz wspierać odporność i płodność.</p>

    <p>Dzięki unikalnej kombinacji kwasów tłuszczowych skutecznie zaspokaja zapotrzebowanie energetyczne krów mlecznych, dostarczając niezbędnych składników odżywczych. Zapewnia optymalną wydajność, szczególnie w okresach ograniczonego pobrania paszy.</p>

    <p>Produkt jest obojętny dla żwacza, charakteryzuje się wysoką strawnością, stanowiąc doskonałe źródło energii. Idealnie sprawdza się w każdej diecie mlecznej, niezależnie od etapu laktacji.</p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 my-8">
        <div class="bg-gray-50 p-6 rounded-lg shadow-sm">
            <h3 class="text-xl font-bold text-secondary mb-4">Dlaczego warto?</h3>
            <ul class="space-y-2 list-disc list-inside">
                <li>Zmniejsza ślad węglowy (CO2)</li>
                <li>Poprawia zdrowie i odporność</li>
                <li>Zwiększa energię</li>
                <li>Zawiera Kwasy Omega 3</li>
                <li>Zmniejsza ryzyko stłuszczenia wątroby</li>
                <li>Wspomaga układ odpornościowy</li>
                <li>Zwiększa ogólną podaż energii</li>
            </ul>
        </div>
        <div class="bg-gray-50 p-6 rounded-lg shadow-sm">
            <h3 class="text-xl font-bold text-secondary mb-4">Analiza</h3>
            <ul class="space-y-2">
                <li><strong>Sucha masa:</strong> 94.0%</li>
                <li><strong>Olej:</strong> 50.0%</li>
                <li><strong>Białko:</strong> 1.4%</li>
                <li><strong>Włókno:</strong> 27.0%</li>
                <li><strong>Popiół:</strong> 3.0%</li>
                <li><strong>ME (MJ / kg DM):</strong> 27.0</li>
            </ul>
        </div>
    </div>

    <div class="overflow-x-auto">
        <h3 class="text-xl font-bold text-primary mb-4">Profil kwasu tłuszczowego</h3>
        <table class="min-w-full bg-white border border-gray-200 rounded-lg overflow-hidden">
            <thead class="bg-gray-100">
                <tr>
                    <th class="py-2 px-4 text-left">Kwas tłuszczowy</th>
                    <th class="py-2 px-4 text-right">% TFA</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
                <tr><td>≤C14:0 Mirystynowy</td><td class="text-right">4.1</td></tr>
                <tr><td>C16:0 Palmitynowy</td><td class="text-right">14.9</td></tr>
                <tr><td>C18:0 Stearynowy</td><td class="text-right">3.1</td></tr>
                <tr><td>C18:1 Oleinowy</td><td class="text-right">41.6</td></tr>
                <tr><td>Suma kwasów Omega 3</td><td class="text-right">6.1</td></tr>
                <tr><td>C18:2 Linolowy</td><td class="text-right">23.5</td></tr>
                <tr><td>C18:3 Kwas Linoleowy</td><td class="text-right">1.7</td></tr>
                <tr><td>C20: 5 EPA</td><td class="text-right">1.7</td></tr>
                <tr><td>C22: 5 DPA</td><td class="text-right">0.2</td></tr>
                <tr><td>C22: 6 DHA</td><td class="text-right">2.5</td></tr>
            </tbody>
        </table>
    </div>

    <div class="my-8">
        <h3 class="text-xl font-bold text-primary mb-4">Porównanie SUPER FAT i mydła wapniowego</h3>
        <div class="overflow-x-auto">
            <table class="min-w-full bg-white border border-gray-200 rounded-lg">
                <thead class="bg-gray-100">
                    <tr>
                        <th class="py-2 px-4 text-left">Parametr</th>
                        <th class="py-2 px-4 text-center">SUPER FAT</th>
                        <th class="py-2 px-4 text-center">Mydło wapniowe</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-gray-200">
                    <tr><td>Energia metaboliczna (MJ/kg)</td><td class="text-center font-bold text-green-600">27.0</td><td class="text-center">23.8</td></tr>
                    <tr><td>Wymóg pH dla stabilności</td><td class="text-center font-bold text-green-600">ŻADNE</td><td class="text-center">6.0-6.8</td></tr>
                    <tr><td>Kompozycja kwasów tłuszczowych</td><td class="text-center font-bold text-green-600">NIENASYCONE</td><td class="text-center">GŁÓWNIE NASYCONE</td></tr>
                    <tr><td>Strawność</td><td class="text-center font-bold text-green-600">WYSOKA</td><td class="text-center">NISKA</td></tr>
                    <tr><td>Pobór suchej masy</td><td class="text-center font-bold text-green-600">ULEPSZONY</td><td class="text-center">ZREDUKOWANY O 2.5%</td></tr>
                    <tr><td>Kwasy tłuszczowe Omega 3</td><td class="text-center font-bold text-green-600">TAK</td><td class="text-center">NIE</td></tr>
                </tbody>
            </table>
        </div>
    </div>

    <div class="bg-blue-50 p-6 rounded-lg border-l-4 border-blue-500">
        <h3 class="text-xl font-bold text-blue-800 mb-2">Punkty kluczowe produktu</h3>
        <ul class="space-y-2 list-disc list-inside text-blue-900">
            <li>Super Fat zawiera lekkostrawne oleje roślinne i morskie, mydło wapniowe – trudniej przyswajalne tłuszcze nasycone.</li>
            <li>Super Fat szybko opuszcza żwacz, chroniąc przed fermentacją, niezależnie od pH. Sole wapnia wymagają stabilnego pH 6,0–6,8, co bywa trudne.</li>
            <li>Super Fat dostarcza Omega-3, wspierające płodność i zdrowie krów, których brak w mydle wapniowym.</li>
            <li>Zawiera także białko, skrobię i cukry, niedostępne w mydle wapniowym.</li>
        </ul>
    </div>

    <div class="mt-4">
        <strong>Dawka:</strong> 300-750g / dzień
    </div>
</div>
`

	products := []SeedProduct{
		// --- PASZE (robot, energyfat, biomlek24, starter calf) ---
		{
			Name:        "BIOMIX ROBOT",
			Category:    "pasze",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Stanowi uzupełnienie dawki pasz objętościowych o średniej zawartości energii i białka.",
			Data:        "Białko surowe 19.2% / Tłuszcz surowy 3.2% / Energia 7.7MJ",
			Dosage:      "Krowy mleczne 1-10kg/szt./dzień",
		},
		{
			Name:        "BIOMIX ENERGYFAT",
			Category:    "pasze",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Przeznaczona do pokrycia niedoborów energetycznych w szczycie laktacji. Zawiera metabolity drożdży oraz tłuszcz chroniony.",
			Data:        "Białko surowe 19% / Tłuszcz surowy 3.4% / Energia 8.0MJ",
			Dosage:      "Krowy mleczne 1-4kg/szt./dzień",
		},
		{
			Name:        "BIOMIX BIOMLEK 24",
			Category:    "pasze",
			Description: "Mieszanka paszowa uzupełniająca dla krów wysokomlecznych. Mieszanka stanowi uzupełnienie dawki pasz objętościowych o niskiej zawartości białka.",
			Data:        "Białko surowe 24% / Tłuszcz surowy 3.8% / Energia 7.1MJ",
			Dosage:      "Krowy mleczne 1-10kg/szt./dzień",
		},
		{
			Name:        "BIOMIX STARTER CALF",
			Category:    "pasze",
			Description: "Mieszanka paszowa uzupełniająca dla cieląt. Świetny start to świetna produkcja potem. Zadbaj o swoje cielęta najlepiej jak możesz.",
			Data:        "Białko surowe 18% / Tłuszcz surowy 3.4%",
			Dosage:      "Cielęta 2tyg. - 6 miesięcy / 1-3kg/szt./dzień",
		},

		// --- PREPARATY MLEKOZASTĘPCZE (biomil gold, silver i platinum) ---
		{
			Name:        "BioMilk GOLD",
			Category:    "preparaty-mlekozastepcze",
			Description: "Oparty w 30% na odtłuszczonym mleku w proszku, serwatce oraz tłuszczach roślinnych. Innowacyjna receptura pozwala na zwiększenie rozpuszczalności oraz strawności produktu.",
			Data:        "0.5% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},
		{
			Name:        "BioMilk SILVER",
			Category:    "preparaty-mlekozastepcze",
			Description: "Oparty o serwatkę w proszku, tłuszcze roślinne oraz odtłuszczone mleko w proszku z dodatkiem drożdży.",
			Data:        "0.9% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},
		{
			Name:        "BioMilk PLATINUM",
			Category:    "preparaty-mlekozastepcze",
			Description: "Najwyższej jakości preparat, oparty w 40% na odtłuszczonym mleku w proszku oraz tłuszczach roślinnych. Innowacyjna metoda produkcji pozwala zwiększyć rozpuszczalność i strawność produktu.",
			Data:        "0% włókna / 21% białko / probiotyk + witamina C. Okres skarmiania od 3 dnia.",
		},

		// --- SPECJALNE (dezodry, tmr full i superfat) ---
		{
			Name:        "BIOMIX DEZODRY",
			Category:    "specjalne",
			Description: "Preparat do suchej dezynfekcji legowisk, ściółki oraz pomieszczeń inwentarskich. Skutecznie ogranicza rozwój bakterii, grzybów i drobnoustrojów chorobotwórczych. Wiąże nadmiar wilgoci, redukuje nieprzyjemne zapachy oraz obniża poziom amoniaku. Zmniejsza ryzyko chorób racic, skóry i układu oddechowego.",
		},
		{
			Name:        "BIOMIX TMR FULL",
			Category:    "specjalne",
			Description: "Kompletna pasza startowa dla cieląt. Zapewnia stabilne żywienie przez cały okres odchowu, eliminując stres związany ze zmianą pasz. Wspiera prawidłowy rozwój przewodu pokarmowego oraz mikroflory jelitowej cieląt. Poprawia pobranie paszy, wspomaga rozwój żwacza i umożliwia wcześniejsze ograniczenie odpajania mlekiem.",
		},
		{
			Name:        "SUPER FAT",
			Category:    "specjalne",
			Description: "Sprawdzony suplement tłuszczowy dla wysokowydajnych krów, który zwiększa energetyczność dawki, pomagając ograniczyć utratę masy ciała, poprawić parametry mleka oraz wspierać odporność i płodność.",
			Content:     superFatContent,
			Data:        "Energia 27 MJ/kg / Tłuszcz 50%",
			Dosage:      "300-750g / dzień",
			ImagePath:   "/static/images/specjalne.jpg",
		},

		// --- PREMIXY (Reszta) ---
		{
			Name:        "Biomix DRY (Krowy Zasuszone)",
			Category:    "premixy",
			Description: "Mieszanka mineralno - witaminowa o odpowiednich proporcjach wapnia, fosforu i magnezu oraz wysokiej dawce witaminy E. Zapewnia bezpieczny okres zasuszenia i świetny start w laktację.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Beta Forte (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka zapewniająca wysokie wskaźniki rozrodu wysokowydajnych krów. Poprawia gospodarkę hormonalną. Dodatek chronionej choliny dodatkowo zabezpiecza wątrobę przed otłuszczeniem.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Energy Drink (Krowy Zasuszone)",
			Category:    "premixy",
			Description: "Mieszanka stymuluje apetyt zaraz po wycieleniu. Uzupełnia gospodarkę wodną i elektrolitową. Zapobiega niedoborom energii i występowaniu porażeń poporodowych. Stymuluje rozwój brodawek żwacza.",
			Dosage:      "Według etykiety",
		},
		{
			Name:        "Biomix Kation / Anion (Krowy Zasuszone)",
			Category:    "premixy",
			Description: "Sole gorzkie w formie chronionej zapewniają stymulację wchłaniania wapnia przez krowy 2-3 tygodnie przed porodem. Zapobiega to zaleganiu poporodowemu i umożliwia zdrową kolejną laktację.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Beta (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka poprawiająca wskaźniki rozrodu bydła. Wysoki poziom witaminy A, D3, E oraz obecność beta karotenu i selen wspiera funkcje rozrodcze organizmu. Wpływa też na zawartość przeciwciał w siarze.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Somatic (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka podnosi odporność przez obecność witaminy E oraz cynku. Efektywnie obniża poziom komórek somatycznych w mleku. Zabezpiecza wątrobę, a dawka biotyny redukuje występowanie kulawizn.",
			Dosage:      "200-300g dziennie",
		},
		{
			Name:        "Biomix Lacto (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka z wysoką zawartością magnezu zapobiega tężyczce pastwiskowej. Uzupełnia dawkę pokarmową w niezbędne mikro i makro składniki w funkcjonowaniu zwierząt laktacyjnych.",
			Dosage:      "200g dziennie",
		},
		{
			Name:        "Biomix Extra KM (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka z witaminą E oraz wysoko przyswajalnymi formami cynku, miedzi, manganu i selenu. Kwas foliowy i niacyna wpływają na metabolizm krów, a wysoki poziom biotyny wspomaga racice.",
			Dosage:      "200-250g dziennie",
		},
		{
			Name:        "Biomix Racice (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka z najwyższym poziomem biotyny oraz wysokim poziomem witaminy E i łatwo przyswajalnym chelacie cynku, miedzi i manganu - działa wspomagająco na racice, skórę i sierść.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix HP Active Max (Krowy w Laktacji)",
			Category:    "premixy",
			Description: "Mieszanka z żywymi kulturami drożdży zapewnia wysoką strawność dawki pokarmowej. Stabilizuje pH żwacza, a dodatek niacyny, biotyny oraz kwasu foliowego stymulują metabolizm do wysokiej i bezpiecznej produkcji mleka.",
			Dosage:      "200-250g dziennie",
		},
		{
			Name:        "Biomix Super Opas",
			Category:    "premixy",
			Description: "Mieszanka wspomagająca tucz opasów. Siarka umożliwia efektywne wykorzystanie składników dawki. Dodatek drożdży poprawia strawność, a zastosowanie witamin w formach chelatowych zwiększa ich przyswajalność. Dla alkalizowanego zboża.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Opas Max",
			Category:    "premixy",
			Description: "Mieszanka wspomagająca intensywny tucz dla opasów. Siarka pozwala na efektywne wykorzystanie składników dawki. Obecność drożdży wspomaga strawność i wykorzystanie dawki pokarmowej.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Drożdże MAX",
			Category:    "premixy",
			Description: "Mieszanka żywych drożdży oraz metabolitów drożdżowych. Produkt został uzupełniony drożdżami piwnymi oraz (MOS), które wspomagają mikroflorę przewodu pokarmowego i zapewniają ochronę jelit.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Drożdże",
			Category:    "premixy",
			Description: "Mieszanka z żywymi kulturami drożdży. Zapewnia prawidłową fermentację żwacza. Zwiększa produkcję mleka. Poprawia stan zdrowotny stada i zabezpiecza krowy przed stresem cieplnym.",
			Dosage:      "100-150g dziennie",
		},
		{
			Name:        "Biomix Bufor Alg",
			Category:    "premixy",
			Description: "Kompleks buforujący o połączonym działaniu alg morskich, tlenku magnezu i kwaśnego węglanu. Zapobiega kwasicy przez odpowiednią stabilizację pH żwacza. Zalecany podczas stresu cieplnego.",
			Dosage:      "100-200g dziennie",
		},
		{
			Name:        "Biomix Immuno",
			Category:    "premixy",
			Description: "Mieszanka dla krów w okresie zmniejszonej odporności redukująca skutki stresu cieplnego. Dzięki wysokiej dawce witaminy C i E znacząco wpływa na wzrost układu immunologicznego redukując wpływ patogenów.",
			Dosage:      "10g dziennie co drugi tydzień",
		},
		{
			Name:        "Biomix Rumen Gold",
			Category:    "premixy",
			Description: "Kompleks buforujący, ale z dodatkiem żywych drożdży i fosforanu jednowapniowego. Zapobiega kwasicy, utrzymując odpowiednie pH w żwaczu. Zwiększa skuteczność trawienia.",
			Dosage:      "300-400g dziennie",
		},
		{
			Name:        "Biomix Sorbent Tox",
			Category:    "premixy",
			Description: "Mieszanka paszowa eliminująca toksyny. Redukuje mikotoksyny i endotoksyny występujące w paszach i przewodzie pokarmowym. Kompleks składników potęguje skuteczność preparatu.",
			Dosage:      "15-50g dziennie",
		},
		{
			Name:        "Biomix Sorbacid",
			Category:    "premixy",
			Description: "Mieszanka ograniczająca fermentację oraz zabezpiecza kiszonki po otwarciu pryzmy przed powstawaniem mykotoksyn. TMR/PMR skutecznie jest chroniony przed zagrzaniem się na stole paszowym.",
			Dosage:      "1-2 kg/ tonę TMR/PMR",
		},
		{
			Name:        "Biomix Soda Lux",
			Category:    "premixy",
			Description: "Preparat buforujący o połączonym działaniu. Skutecznie zapobiega kwasicy żwacza oraz stabilizuje pH przewodu pokarmowego, wspierając prawidłowe trawienie i zdrowie zwierząt.",
			Dosage:      "300g dziennie",
		},
		{
			Name:        "Biomix Inokulant (Zakiszacz)",
			Category:    "premixy",
			Description: "Kompleks mikroorganizmów, który zawiera kultury bakterii kwasu propionowego oraz mlekowego. Kwasy te hamują procesy gnilne, ograniczają grzanie się kiszonki i znacznie wydłużają ich trwałość.",
			Dosage:      "Według etykiety",
		},
		{
			Name:        "Biomix Calf",
			Category:    "premixy",
			Description: "Mieszanka dla cieląt i młodzieży. Gwarantuje prawidłowy rozwój układu rozrodczego i hormonalnego. Wspomaga właściwy rozwój gruczołów płciowych. U zacielonych jałówek wspiera utrzymanie płodu.",
			Dosage:      "Zależne od wieku",
		},
	}

	for _, p := range products {
		catID, ok := categoryIDs[p.Category]
		if !ok {
			log.Printf("Category not found for product %s: %s", p.Name, p.Category)
			continue
		}

		// Use description as content if content is empty
		content := p.Content
		if content == "" {
			content = fmt.Sprintf("<p>%s</p>", p.Description)
		}

		_, err := DB.Exec(`INSERT INTO products (name, description, dosage, data, content, image_path, category_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.Name, p.Description, p.Dosage, p.Data, content, p.ImagePath, catID)
		if err != nil {
			log.Printf("Failed to insert product %s: %v", p.Name, err)
		}
	}

	return nil
}
