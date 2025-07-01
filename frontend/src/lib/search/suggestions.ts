import { metadataElements } from "$lib/api";
import { filterMatches } from "$lib/metadata-panel/auto-complete";

const allElements = metadataElements().then((elements) =>
	Object.values(elements)
		.flat()
		.map((element) => element.name),
);

let prevSuggestions = allElements;
export async function suggest(value: string, selected: number): Promise<[string[], number]> {
	const elements = await allElements;
	const newSuggestions = filterMatches(elements, value, "-", true);

	let newSelected = value === "" ? -1 : 0;
	const prev = await prevSuggestions;
	if (selected > 0 && typeof prev[selected] === "string") {
		newSelected = newSuggestions.findIndex((choice) => choice === prev[selected]);
	}

	prevSuggestions = Promise.resolve(newSuggestions);
	return [newSuggestions, newSelected];
}
