import { type MetadataElements } from "$lib/api";
import { filterMatches } from "$lib/metadata-panel/auto-complete";

export function metadataSuggest(metadataElements: MetadataElements) {
	const elements = Object.values(metadataElements)
		.flat()
		.map((element) => element.name);

	let prevSuggestions = elements;
	return function suggest(value: string, selected: number): [string[], number] {
		const newSuggestions = filterMatches(elements, value, "-", true);

		let newSelected = value === "" ? -1 : 0;
		if (selected > 0 && typeof prevSuggestions[selected] === "string") {
			newSelected = newSuggestions.findIndex((choice) => choice === prevSuggestions[selected]);
		}

		prevSuggestions = newSuggestions;
		return [newSuggestions, newSelected];
	};
}
