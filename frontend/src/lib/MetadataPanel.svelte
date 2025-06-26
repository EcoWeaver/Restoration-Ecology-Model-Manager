<script lang="ts">
	import { X } from "@lucide/svelte";
	import { type Metadata, type MetadataElements, type MetadataType } from "$lib/api";
	import YearInput from "./metadata-panel/YearInput.svelte";
	import DateInput from "./metadata-panel/DateInput.svelte";
	import OrcidInput from "./metadata-panel/OrcidInput.svelte";
	import AutoComplete from "./metadata-panel/AutoComplete.svelte";
	import { suggestChoices, suggestSpecies } from "./metadata-panel/auto-complete";
	import Input from "./metadata-panel/Input.svelte";

	let {
		elements,
		identifier = $bindable(),
		identifierEditable = true,
		metadata = $bindable(),
	}: {
		elements: MetadataElements;
		identifier: string;
		identifierEditable?: boolean;
		metadata: Metadata[];
	} = $props();

	type FlatMeta = {
		type?: MetadataType;
		choices?: string[];
	};

	const flatMetadata = new Map<string, FlatMeta>();
	for (const { name, type, choices: c } of Object.values(elements).flat()) {
		const choices = c?.map(({ label }) => label);
		flatMetadata.set(name, { type, choices });
	}

	const meta: (Metadata & { mandatory: boolean })[] = $state([]);

	const initialMetadata = $state.snapshot(metadata);
	for (const kind of elements.mandatory) {
		const index = initialMetadata.findIndex((m) => m.kind === kind.name);
		if (index >= 0) {
			const value = initialMetadata[index].value;
			initialMetadata.splice(index, 1);
			meta.push({ mandatory: true, kind: kind.name, value });
		} else {
			meta.push({ mandatory: true, kind: kind.name, value: "" });
		}
	}
	for (const { kind, value } of initialMetadata) {
		meta.push({ mandatory: false, kind, value });
	}

	export function addMetadata() {
		meta.push({ mandatory: false, kind: "", value: "" });
	}

	$effect(() => {
		metadata = meta.map(({ kind, value }) => ({ kind, value }));
	});
</script>

<h2 class="mb-1 px-2 text-xl font-bold">Metadata</h2>

<div class="px table">
	<div class="mb-2 table-row">
		<span class="table-cell px-2 py-2 font-bold">kind</span>
		<span class="table-cell py-2 pl-2 font-bold">value</span>
	</div>

	<div class="table-row">
		<span class="table-cell px-2 py-1">identifier</span>
		<span class:pl-2={identifierEditable} class="table-cell py-1">
			{#if identifierEditable}
				<Input placeholder="doi" bind:value={identifier} />
			{:else}
				<a
					class="ml-1 block w-full pl-1 text-blue-700 outline-indigo-400 hover:underline focus-visible:outline-2 dark:text-blue-400"
					href={`https://doi.org/${identifier}`}
					target="_blank"
					rel="noopener noreferrer"
				>
					{identifier}
				</a>
			{/if}
		</span>
	</div>

	{#each meta as m, idx (idx)}
		{@const { type, choices } = flatMetadata.get(m.kind) || {}}

		<div class="table-row">
			{#if m.mandatory}
				<div class="table-cell align-top">
					<div class="px-2 py-1">{m.kind}</div>
				</div>
			{:else}
				<div class="table-cell align-top">
					<select
						bind:value={m.kind}
						title={m.kind}
						class="field-sizing-content rounded px-2 py-1 outline-indigo-400 focus-visible:outline-2"
					>
						<option class="hidden" value="" disabled selected>-- select a kind --</option>
						{#each Object.entries(elements) as [group, values] (group)}
							<optgroup class="dark:bg-neutral-700" label={group}>
								{#each values as { name: kind } (kind)}
									<option class="dark:bg-neutral-700" value={kind}>{kind}</option>
								{/each}
							</optgroup>
						{/each}
					</select>
				</div>
			{/if}

			<span class="table-cell py-0.75 pl-2">
				{#if m.kind === ""}
					<Input value="" disabled placeholder="please select a metadata kind" />
				{:else if type === "choice"}
					<select
						bind:value={m.value}
						title={m.value}
						class="picker-icon:text-black h-8 w-full cursor-pointer rounded border border-neutral-600 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
						required
					>
						<option class="hidden" value="" disabled selected>-- select an option --</option>
						{#each choices! as choice (choice)}
							<option class="cursor-pointer dark:bg-neutral-700" value={choice}>{choice}</option>
						{/each}
					</select>
				{:else if type === "date"}
					<DateInput bind:value={m.value} />
				{:else if type === "orcid"}
					<OrcidInput bind:value={m.value} />
				{:else if type === "url"}
					<Input placeholder={type} bind:value={m.value} type="url" />
				{:else if type === "year"}
					<YearInput {type} bind:value={m.value} />
				{:else if type === "species"}
					<AutoComplete placeholder="species" suggest={suggestSpecies} bind:value={m.value} />
				{:else if choices}
					<AutoComplete
						placeholder={type || m.kind}
						suggest={suggestChoices(choices)}
						bind:value={m.value}
					/>
				{:else}
					<textarea
						class="field-sizing-content min-h-8 w-full max-w-128 rounded border border-neutral-600 py-1 pr-2 pl-3 outline-indigo-400 focus-visible:outline-2"
						rows="1"
						placeholder={m.kind}
						bind:value={m.value}
					></textarea>
				{/if}
			</span>

			{#if !m.mandatory}
				<div class="relative table-cell">
					<button
						class="absolute top-2 left-[4px] inline h-6 cursor-pointer rounded-4xl outline-red-500 focus-visible:outline-2"
						onclick={() => meta.splice(idx, 1)}
					>
						<X class="stroke-red-500" />
					</button>
				</div>
			{/if}
		</div>
	{/each}
</div>
