import type { Onto } from "$lib/api";

export async function suggestSpecies(value: string): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const res = await fetch(
		`https://api.gbif.org/v1/species/suggest?status=ACCEPTED&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		canonicalName: string;
	}[] = await res.json();

	return suggestions.map((suggestion) => suggestion.canonicalName);
}

export async function suggestOLS4(value: string, onto: Onto): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const res = await fetch(
		`https://www.ebi.ac.uk/ols4/api/suggest?ontology=${onto.onto}&childrenOf=${onto.iri}&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		response: {
			docs: {
				autosuggest: string;
			}[];
		};
	} = await res.json();

	return suggestions.response.docs.map((suggestion) => suggestion.autosuggest);
}
