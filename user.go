package main

type User struct {
	PokemonMap map[string]Pokemon
}

type Pokemon struct {
	data []byte
}

func (u *User) AddPokemon(name string, data []byte) {
	u.PokemonMap[name] = Pokemon{
		data: data,
	}
}

func (u *User) GetPokemon(name string) (pokemonData []byte,pokemonCatched bool) {
	// Check if the Pokemon exists in the map
	pokemon, ok := u.PokemonMap[name]
	if !ok {
		return nil, false
	}
	// Return the data field of the Pokemon
	return pokemon.data, true
}