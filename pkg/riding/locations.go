package riding

import "math/rand"

var locations = map[string][]string{
	"Belfast": {
		"North",
		"South",
		"East",
		"West",
	},
	"Lisburn": {
		"Lambeg",
		"Finaghy",
		"Lisburn City Centre",
	},
	"Newtownabbey": {
		"Glengormley",
		"Whiteabbey",
		"Jordanstown",
	},
}

func GetRandomLocation() string {
	cities := make([]string, 0, len(locations))
	for city := range locations {
		cities = append(cities, city)
	}

	areas := locations[cities[rand.Intn(len(cities))]]
	return areas[rand.Intn(len(areas))]
}
