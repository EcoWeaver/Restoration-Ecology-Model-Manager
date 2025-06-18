export function suggestChoices(choices: string[]) {
	return async function suggest(value: string): Promise<string[]> {
		if (value === "") {
			return [];
		}

		return choices.filter(
			(choice) => choice.toLowerCase().startsWith(value.toLowerCase()) && choice !== value,
		);
	};
}

type SpeciesSuggestion = {
	status: string;
	canonicalName?: string;
};

export async function suggestSpecies(value: string): Promise<string[]> {
	if (value === "") {
		return [];
	}

	let response = await fetch(
		`https://api.gbif.org/v1/species/suggest?rank=SPECIES&q=${encodeURIComponent(value)}`,
	);
	let suggestions: SpeciesSuggestion[] = await response.json();

	return suggestions
		.filter((suggestion) => suggestion.status == "ACCEPTED")
		.map((suggestion) => suggestion.canonicalName)
		.filter((suggestion) => suggestion)
		.filter((suggestion) => suggestion !== value);
}
