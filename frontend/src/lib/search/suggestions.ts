import { type MetadataElements } from "$lib/api";

export function metadataSuggest(
	metadataElements: MetadataElements,
): (value: string, selected: number) => [string[], number] {
	const elements = Object.values(metadataElements)
		.flat()
		.map((element) => element.name);

	let prevSuggestions = elements;
	return function suggest(value: string, selected: number): [string[], number] {
		const newSuggestions = filterMatches(elements, value);

		let newSelected = value === "" ? -1 : 0;
		if (selected > 0) {
			newSelected = newSuggestions.findIndex((choice) => choice === prevSuggestions[selected]);
		}

		prevSuggestions = newSuggestions;
		return [newSuggestions, newSelected];
	};
}

function filterMatches(elements: string[], value: string): string[] {
	return elements
		.filter((elem) => {
			return elem.toLowerCase().includes(value.toLowerCase());
		})
		.sort((l, r) => {
			const lFirst = l.startsWith(value);
			const rFirst = r.startsWith(value);
			if (lFirst || rFirst) {
				// if either of them matches prefix, prefer that one
				return +rFirst - +lFirst;
			}

			// ... split the elements into words
			const lAny = l.split(" ").some((elem) => elem.startsWith(value));
			const rAny = r.split(" ").some((elem) => elem.startsWith(value));

			// and otherwise prefer the one where any of the words match
			return +rAny - +lAny;
		});
}
