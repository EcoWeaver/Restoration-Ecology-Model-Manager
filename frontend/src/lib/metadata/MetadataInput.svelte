<script lang="ts">
	import { useElements, type Metadata, type MetadataElement } from "#lib/api.js";
	import AutoComplete from "./AutoComplete.svelte";
	import DateInput from "./DateInput.svelte";
	import DoiInput from "./DoiInput.svelte";
	import Input from "./Input.svelte";
	import OrcidInput from "./OrcidInput.svelte";
	import UuidInput from "./UuidInput.svelte";
	import YearInput from "./YearInput.svelte";

	let { meta = $bindable() }: { meta: Metadata } = $props();

	const { getMetadataElements } = useElements();
	const flatMetadata = new Map<string, Omit<MetadataElement, "name">>();
	for (const { name, type, choices, onto } of Object.values(getMetadataElements()).flat()) {
		flatMetadata.set(name, { type, choices, onto });
	}

	const { type, choices, onto } = $derived(flatMetadata.get(meta.kind) || {});
</script>

{#if meta.kind === ""}
	<Input value="" disabled placeholder="-- select a kind --" />
{:else if type === "bioportal"}
	<AutoComplete
		placeholder={`subclass of '${onto!.label}' in ontology '${onto!.onto}'`}
		bind:meta
	/>
{:else if type === "choice"}
	<select
		bind:value={meta.value}
		title={meta.value}
		class="field-sizing-content h-8 w-full cursor-pointer rounded border border-neutral-600 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
		required
	>
		<option class="hidden" value="" disabled selected>-- select an option --</option>
		{#each choices! as choice (choice)}
			<option class="cursor-pointer dark:bg-neutral-700" value={choice}>
				{choice}
			</option>
		{/each}
	</select>
{:else if type === "date"}
	<DateInput bind:value={meta.value} />
{:else if type === "doi"}
	<DoiInput bind:value={meta.value} />
{:else if type === "ols4"}
	<AutoComplete
		placeholder={`subclass of '${onto!.label}' in ontology '${onto!.onto}'`}
		bind:meta
	/>
{:else if type === "orcid"}
	<OrcidInput bind:value={meta.value} />
{:else if type === "species"}
	<AutoComplete placeholder={type} bind:meta />
{:else if type === "url"}
	<Input placeholder={type} bind:value={meta.value} type="url" />
{:else if type === "uuid"}
	<UuidInput bind:value={meta.value} />
{:else if type === "year"}
	<YearInput {type} bind:value={meta.value} />
{:else}
	<textarea
		class="field-sizing-content min-h-8 w-full max-w-lg rounded border border-neutral-600 py-1 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
		rows="1"
		bind:value={meta.value}></textarea>
{/if}
