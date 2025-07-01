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

export function suggestChoices(choices: string[]) {
	let prevChoices = choices;

	return async function suggest(value: string, selected: number): Promise<[string[], number]> {
		if (value === "") {
			return [choices, selected];
		}

		const newChoices = filterMatches(choices, value, " ", false);

		let newSelected = 0;
		if (selected !== 0 && typeof prevChoices[selected] === "string") {
			const index = newChoices.findIndex((choice) => choice === prevChoices[selected]);
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

	const response = await fetch(
		`https://api.gbif.org/v1/species/suggest?rank=SPECIES&q=${encodeURIComponent(value)}`,
	);
	const suggestions: SpeciesSuggestion[] = await response.json();
	const newSuggestions = suggestions
		.filter((suggestion) => suggestion.status === "ACCEPTED")
		.map((suggestion) => suggestion.canonicalName)
		.filter((suggestion) => typeof suggestion === "string")
		.filter((suggestion) => suggestion !== value);

	return [newSuggestions, Math.min(selected, Math.max(0, newSuggestions.length - 1))];
}
