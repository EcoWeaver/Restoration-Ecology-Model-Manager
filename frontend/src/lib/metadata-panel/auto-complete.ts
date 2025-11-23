import type { OLS4 } from "$lib/api";

export function filterMatches(
	elements: string[],
	value: string,
	sep: string,
	includeExact: boolean,
) {
	return elements
		.filter((elem) => {
			return elem.toLowerCase().includes(value.toLowerCase()) && (includeExact || elem !== value);
		})
		.sort((l, r) => {
			const lFirst = l.startsWith(value);
			const rFirst = r.startsWith(value);
			if (lFirst || rFirst) {
				// if either of them matches prefix, prefer that one
				return +rFirst - +lFirst;
			}

			// ... split the elements into words
			const lAny = l.split(sep).some((elem) => elem.startsWith(value));
			const rAny = r.split(sep).some((elem) => elem.startsWith(value));

			// and otherwise prefer the one where any of the words match
			return +rAny - +lAny;
		});
}

export async function suggestSpecies(value: string): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const response = await fetch(
		`https://api.gbif.org/v1/species/suggest?rank=SPECIES&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		status: string;
		canonicalName?: string;
	}[] = await response.json();
	const filtered = suggestions
		.filter((suggestion) => suggestion.status === "ACCEPTED")
		.map((suggestion) => suggestion.canonicalName)
		.filter((suggestion) => suggestion !== undefined)
		.filter((suggestion) => suggestion !== value);

	return filtered;
}

export async function suggestOLS4(value: string, ols4: OLS4): Promise<string[]> {
	if (value === "") {
		return [];
	}

	const response = await fetch(
		`https://www.ebi.ac.uk/ols4/api/search?fieldList=label&queryFields=label&ontology=${ols4.onto}&childrenOf=${ols4.iri}&q=${encodeURIComponent(value)}`,
	);
	const suggestions: {
		response: {
			docs: {
				label: string;
			}[];
		};
	} = await response.json();
	const filtered = suggestions.response.docs
		.map((suggestion) => suggestion.label)
		.filter((suggestion) => suggestion !== value);

	return filtered;
}
