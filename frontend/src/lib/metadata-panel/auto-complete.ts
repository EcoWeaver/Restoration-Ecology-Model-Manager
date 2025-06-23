export function suggestChoices(choices: string[]) {
	let prevChoices = choices;

	return async function suggest(value: string, selected: number): Promise<[string[], number]> {
		if (value === "") {
			return [choices, selected];
		}

		let newChoices = choices.filter(
			(choice) => choice.toLowerCase().startsWith(value.toLowerCase()) && choice !== value,
		);

		let newSelected = 0;
		if (selected !== 0 && typeof prevChoices[selected] === "string") {
			let index = newChoices.findIndex((choice) => choice === prevChoices[selected]);
			newSelected = index === -1 ? 0 : index;
		}

		prevChoices = newChoices;
		return [newChoices, newSelected];
	};
}

type SpeciesSuggestion = {
	status: string;
	canonicalName?: string;
};

export async function suggestSpecies(value: string, selected: number): Promise<[string[], number]> {
	if (value === "") {
		return [[], selected];
	}

	let response = await fetch(
		`https://api.gbif.org/v1/species/suggest?rank=SPECIES&q=${encodeURIComponent(value)}`,
	);
	let suggestions: SpeciesSuggestion[] = await response.json();
	let newSuggestions = suggestions
		.filter((suggestion) => suggestion.status == "ACCEPTED")
		.map((suggestion) => suggestion.canonicalName)
		.filter((suggestion) => typeof suggestion === "string")
		.filter((suggestion) => suggestion !== value);

	return [newSuggestions, Math.min(selected, Math.max(0, newSuggestions.length - 1))];
}
