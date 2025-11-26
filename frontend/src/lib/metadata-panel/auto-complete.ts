import type { OLS4 } from "$lib/api";

export async function suggestSpecies(value: string): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const response = await fetch(
		`https://api.gbif.org/v1/species/suggest?rank=SPECIES&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		status: string;
		canonicalName: string;
	}[] = await response.json();

	return suggestions
		.filter((suggestion) => suggestion.status === "ACCEPTED")
		.map((suggestion) => suggestion.canonicalName);
}

export async function suggestOLS4(value: string, ols4: OLS4): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const response = await fetch(
		`https://www.ebi.ac.uk/ols4/api/suggest?ontology=${ols4.onto}&childrenOf=${ols4.iri}&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		response: {
			docs: {
				autosuggest: string;
			}[];
		};
	} = await response.json();

	return suggestions.response.docs.map((suggestion) => suggestion.autosuggest);
}
