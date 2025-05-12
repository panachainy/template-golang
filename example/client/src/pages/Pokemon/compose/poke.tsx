import { useQuery } from '@tanstack/react-query'

interface Pokemon {
  name: string
  url: string
}

interface PokemonResponse {
  count: number
  next: string | null
  previous: string | null
  results: Pokemon[]
}

const useFetchPokemon = () =>
  useQuery({
    queryKey: ['pokemon'],
    queryFn: async (): Promise<PokemonResponse> => {
      const response = await fetch('https://pokeapi.co/api/v2/pokemon')
      return await response.json()
    },
  })

function PokemonData() {
  const pokemonQuery = useFetchPokemon()

  if (pokemonQuery.isLoading) return <p>Loading...</p>

  return (
    <div>
      Pokemon
      {pokemonQuery.data && <p>{JSON.stringify(pokemonQuery.data)}</p>}
      {pokemonQuery.isError && <p>{JSON.stringify(pokemonQuery.error)}</p>}
    </div>
  )
}

export default PokemonData
