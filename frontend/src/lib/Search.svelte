<script lang="ts">
	import MetadataSearch from "./MetadataSearch.svelte";
	import IdentifierSearch from "./IdentifierSearch.svelte";
	import type { ChipElement } from "./search/Chip.svelte";

	let {
		params = [],
		identifier = "",
		selected,
	}: { params?: ChipElement[]; identifier?: string; selected: "meta" | "ident" } = $props();

	type Thing = { value: "meta" | "ident"; tooltip: string };

	let things: Thing[] = [
		{ value: "meta", tooltip: "search by metadata" },
		{ value: "ident", tooltip: "search by identifier" },
	];
</script>

<div class="flex p-4 pr-0" role="radiogroup">
	{#each things as { value, tooltip }}
		<label
			title={tooltip}
			class="flex h-10 cursor-pointer place-items-center bg-indigo-100 px-2 first:rounded-l last:rounded-r hover:bg-indigo-200 has-checked:bg-indigo-300 hover:has-checked:bg-indigo-400"
		>
			<input class="hidden" type="radio" name="search" {value} bind:group={selected} />
			<span>{value}</span>
		</label>
	{/each}
</div>

{#snippet button({ onclick }: { onclick: () => any })}
	<button
		{onclick}
		class="m-4 ml-0 flex h-10 cursor-pointer place-items-center rounded bg-indigo-300 px-4 hover:bg-indigo-400"
	>
		<span class="font-medium">search</span>
	</button>
{/snippet}

{#if selected == "meta"}
	<MetadataSearch {params} {button} />
{:else if selected == "ident"}
	<IdentifierSearch value={identifier} {button} />
{/if}
