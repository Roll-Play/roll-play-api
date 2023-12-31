package entities

import (
	"fmt"
	"math/rand"
)

type NPC struct {
	Gender      string `json:"gender"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Race        string `json:"race"`
	Appearance  string `json:"appearance"`
	Age         string `json:"age"`
	Build       string `json:"build"`
	Description string `json:"description"`
}

type NPCAttribute int

const (
	Gender NPCAttribute = iota
	FullName
	Race
	Appearance
	Age
	Build
)

func (npc *NPC) getRandomAttribute(attribute NPCAttribute) {
	switch attribute {
	case Gender:
		npc.Gender = Genders[rand.Intn(len(Genders))]
		return
	case FullName:
		var name, surname string
		if npc.Gender == "Male" {
			name = MaleNames[rand.Intn(len(MaleNames))]
			surname = MaleNames[rand.Intn(len(MaleNames))]
		} else if npc.Gender == "Female" {
			name = FemaleNames[rand.Intn(len(FemaleNames))]
			surname = FemaleNames[rand.Intn(len(FemaleNames))]
		} else {
			names := [2][]string{
				MaleNames[:],
				FemaleNames[:],
			}

			nameList := names[rand.Intn(2)]

			name = nameList[rand.Intn(len(nameList))]
			surname = nameList[rand.Intn(len(nameList))]
		}

		npc.Name = name
		npc.Surname = surname
		return
	case Race:
		npc.Race = "Human"
		return
	case Appearance:
		npc.Appearance = AppearanceAdjectives[rand.Intn(len(AppearanceAdjectives))]
		return
	case Age:
		npc.Age = AgeAdjectives[rand.Intn(len(AgeAdjectives))]
		return
	case Build:
		npc.Build = BuildAdjectives[rand.Intn(len(BuildAdjectives))]
		return
	}
}

func (npc *NPC) description() {
	npc.Description = fmt.Sprintf(
		"%s %s is a %s, %s and %s %s",
		npc.Name,
		npc.Surname,
		npc.Appearance,
		npc.Build,
		npc.Age,
		npc.Race,
	)
}

func NewNPC() *NPC {
	npc := new(NPC)

	for i := NPCAttribute(0); i <= 5; i++ {
		npc.getRandomAttribute(i)
	}

	npc.description()

	return npc
}

var Genders = [...]string{
	"Male",
	"Female",
	"Non-binary",
}

var AppearanceAdjectives = [...]string{
	"Attractive",
	"Beautiful",
	"Handsome",
	"Stunning",
	"Intolerant",
	"Gorgeous",
	"Pretty",
	"Ugly",
	"Caring",
	"Cute",
	"Lovely",
	"Flirtatious",
	"Elegant",
	"Aggressive",
	"Charming",
	"Timid",
	"Moody",
	"Radiant",
	"Sophisticated",
	"Gullible",
	"Alluring",
	"Grumpy",
	"Rude",
	"Lazy",
	"Stylish",
	"Classy",
}

var AgeAdjectives = [...]string{
	"Child",
	"Teenage",
	"Young adult",
	"Adult",
	"Senior citizen",
	"Elderly",
	"Octogenarian",
	"Centenarian",
	"Mature",
	"Spry",
	"Venerable",
	"Timeless",
	"Aged",
	"Long-lived",
	"Ancient",
}

var BuildAdjectives = [...]string{
	"Well-Built",
	"Plump",
	"Thin",
	"Fat",
	"Slim",
	"Petite",
	"Athletic",
	"Stocky",
	"Lanky",
	"Stout",
	"Curvy",
	"Slender",
	"Muscular",
	"Chubby",
	"Skinny",
}
var HeightAdjectives = [...]string{
	"Short",
	"Medium-height",
	"Tall",
	"Petite",
	"Lanky",
	"Stunted",
	"Stumpy",
}

var FemaleNames = [...]string{
	"Abigayl",
	"Aebria",
	"Aeobreia",
	"Breia",
	"Aedria",
	"Aodreia",
	"Dreia",
	"Aeliya",
	"Aliya",
	"Aella",
	"Aemilya",
	"Aemma",
	"Aemy",
	"Amy",
	"Ami",
	"Aeria",
	"Arya",
	"Aeva",
	"Aevelyn",
	"Evylann",
	"Alaexa",
	"Alyxa",
	"Alina",
	"Aelina",
	"Aelinea",
	"Allisann",
	"Allysann",
	"Alyce",
	"Alys",
	"Alysea",
	"Alyssia",
	"Aelyssa",
	"Amelya",
	"Maelya",
	"Andreya",
	"Aendrea",
	"Arianna",
	"Aryanna",
	"Arielle",
	"Aryell",
	"Ariella",
	"Ashlena",
	"Aurora",
	"Avaery",
	"Avyrie",
	"Bella",
	"Baella",
	"Brooklinea",
	"Bryanna",
	"Brynna",
	"Brinna",
	"Caemila",
	"Chloe",
	"Chloeia",
	"Claira",
	"Clayre",
	"Clayra",
	"Delyla",
	"Dalyla",
	"Elisybeth",
	"Aelisabeth",
	"Ellia",
	"Ellya",
	"Elyana",
	"Eliana",
	"Eva",
	"Falyne",
	"Genaesis",
	"Genaesys",
	"Gianna",
	"Jianna",
	"Janna",
	"Graece",
	"Grassa",
	"Haenna",
	"Hanna",
	"Halya",
	"Harperia",
	"Peria",
	"Hazyl",
	"Hazel",
	"Jasmyne",
	"Jasmine",
	"Jocelyne",
	"Joceline",
	"Celine",
	"Kaelia",
	"Kaelya",
	"Kathryne",
	"Kathrine",
	"Kayla",
	"Kaila",
	"Kymber",
	"Kimbera",
	"Layla",
	"Laylanna",
	"Leia",
	"Leya",
	"Leah",
	"Lilia",
	"Lylia",
	"Luna",
	"Maedisa",
	"Maelania",
	"Melania",
	"Maya",
	"Mya",
	"Myla",
	"Milae",
	"Naomi",
	"Naome",
	"Natalya",
	"Talya",
	"Nathylie",
	"Nataliae",
	"Thalia",
	"Nicola",
	"Nikola",
	"Nycola",
	"Olivya",
	"Alivya",
	"Penelope",
	"Paenelope",
	"Pynelope",
	"Rianna",
	"Ryanna",
	"Ruby",
	"Ryla",
	"Samaentha",
	"Samytha",
	"Sara",
	"Sarah",
	"Savannia",
	"Scarletta",
	"Sharlotta",
	"Caerlotta",
	"Sophya",
	"Stella",
	"Stylla",
	"Valentyna",
	"Valerya",
	"Valeria",
	"Valia",
	"Valea",
	"Victorya",
	"Vilettia",
	"Ximena",
	"Imaena",
	"Ysabel",
	"Zoe",
	"Zoeia",
	"Zoea",
	"Zoesia",
}

var MaleNames = [...]string{
	"Aaryn",
	"Aaro",
	"Aarus",
	"Abramus",
	"Abrahm",
	"Abyl",
	"Abelus",
	"Adannius",
	"Adanno",
	"Aedam",
	"Adym",
	"Adamus",
	"Aedrian",
	"Aedrio",
	"Aedyn",
	"Aidyn",
	"Aelijah",
	"Elyjah",
	"Aendro",
	"Androe",
	"Aenry",
	"Hynroe",
	"Hynrus",
	"Aethan",
	"Aethyn",
	"Aevan",
	"Evyn",
	"Evanus",
	"Alecks",
	"Alyx",
	"Alexandyr",
	"Xandyr",
	"Alyn",
	"Alaen",
	"Andrus",
	"Aendrus",
	"Anglo",
	"Aenglo",
	"Anglus",
	"Antony",
	"Antonyr",
	"Astyn",
	"Astinus",
	"Axelus",
	"Axyl",
	"Benjamyn",
	"Benjamyr",
	"Braidyn",
	"Brydus",
	"Braddeus",
	"Brandyn",
	"Braendyn",
	"Bryus",
	"Bryne",
	"Bryn",
	"Branus",
	"Caeleb",
	"Caelyb",
	"Caerlos",
	"Carlus",
	"Cameryn",
	"Camerus",
	"Cartus",
	"Caertero",
	"Charlus",
	"Chaerles",
	"Chyrles",
	"Christophyr",
	"Christo",
	"Chrystian",
	"Chrystan",
	"Connorus",
	"Connyr",
	"Daemian",
	"Damyan",
	"Daenyel",
	"Danyel",
	"Davyd",
	"Daevo",
	"Dominac",
	"Dylaen",
	"Dylus",
	"Elius",
	"Aeli",
	"Elyas",
	"Helius",
	"Helian",
	"Emilyan",
	"Emilanus",
	"Emmanus",
	"Emynwell",
	"Ericus",
	"Eryc",
	"Eryck",
	"Ezekius",
	"Zeckus",
	"Ezekio",
	"Ezrus",
	"Yzra",
	"Gabrael",
	"Gaebriel",
	"Gael",
	"Gayl",
	"Gayel",
	"Gaeus",
	"Gavyn",
	"Gaevyn",
	"Goshwa",
	"Joshoe",
	"Graysus",
	"Graysen",
	"Gwann",
	"Ewan",
	"Gwyllam",
	"Gwyllem",
	"Haddeus",
	"Hudsyn",
	"Haesoe",
	"Haesys",
	"Haesus",
	"Handus",
	"Handyr",
	"Hantus",
	"Huntyr",
	"Haroldus",
	"Haryld",
	"Horgus",
	"Horus",
	"Horys",
	"Horyce",
	"Hosea",
	"Hosius",
	"Iaen",
	"Yan",
	"Ianus",
	"Ivaen",
	"Yvan",
	"Jaecoby",
	"Jaecob",
	"Jaeden",
	"Jaedyn",
	"Jaeremiah",
	"Jeremus",
	"Jasyn",
	"Jaesen",
	"Jaxon",
	"Jaxyn",
	"Jaxus",
	"Johnus",
	"Jonus",
	"Jonaeth",
	"Jonathyn",
	"Jordus",
	"Jordyn",
	"Josaeth",
	"Josephus",
	"Josaeus",
	"Josayah",
	"Jovanus",
	"Giovan",
	"Julyan",
	"Julyo",
	"Jyck",
	"Jaeck",
	"Jacus",
	"Kaevin",
	"Kevyn",
	"Vinkus",
	"Laevi",
	"Levy",
	"Levius",
	"Landyn",
	"Laendus",
	"Leo",
	"Leonus",
	"Leonaerdo",
	"Leonyrdo",
	"Lynardus",
	"Lincon",
	"Lyncon",
	"Linconus",
	"Logaen",
	"Logus",
	"Louis",
	"Lucius",
	"Lucae",
	"Lucaen",
	"Lucaes",
	"Lucoe",
	"Lucus",
	"Lyam",
	"Maeson",
	"Masyn",
	"Maetho",
	"Mathoe",
	"Matteus",
	"Matto",
	"Maxus",
	"Maximus",
	"Maximo",
	"Maxymer",
	"Mychael",
	"Mygwell",
	"Miglus",
	"Mythro",
	"Mithrus",
	"Naemo",
	"Naethyn",
	"Nathanus",
	"Naethynel",
	"Nicholaes",
	"Nycholas",
	"Nicholys",
	"Nicolus",
	"Nolyn",
	"Nolanus",
	"Olivyr",
	"Alivyr",
	"Olivus",
	"Oscarus",
	"Oscoe",
	"Raen",
	"Ryn",
	"Robertus",
	"Robett",
	"Bertus",
	"Romyn",
	"Romanus",
	"Ryderus",
	"Ridyr",
	"Samwell",
	"Saemuel",
	"Santegus",
	"Santaegus",
	"Sybasten",
	"Bastyen",
	"Tago",
	"Aemo",
	"Tagus",
	"Theodorus",
	"Theodus",
	"Thaeodore",
	"Thomys",
	"Thomas",
	"Tommus",
	"Tylus",
	"Tilyr",
	"Uwyn",
	"Oewyn",
	"Victor",
	"Victyr",
	"Victorus",
	"Vincynt",
	"Vyncent",
	"Vincentus",
	"Wyttus",
	"Wyaett",
	"Xavius",
	"Havius",
	"Xavyer",
	"Yago",
	"Tyago",
	"Tyego",
	"Ysaac",
	"Aisaac",
	"Ysaiah",
	"Aisiah",
	"Siahus",
	"Zacharus",
	"Zachar",
	"Zachaery",
}
