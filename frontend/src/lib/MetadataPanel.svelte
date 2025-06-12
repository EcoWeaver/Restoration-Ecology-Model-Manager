<script lang="ts">
	import { X } from "@lucide/svelte";
	import { type Metadata, type MetadataElements, type MetadataType } from "$lib/api";
	import YearInput from "./metadata-panel/YearInput.svelte";
	import DateInput from "./metadata-panel/DateInput.svelte";
	import OrcidInput from "./metadata-panel/OrcidInput.svelte";
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

	let flatMetadata = new Map<string, FlatMeta>();
	for (let { name, type, choices: c } of Object.values(elements).flat()) {
		let choices = c?.map(({ label }) => label);
		flatMetadata.set(name, { type, choices });
	}

	let meta: (Metadata & { mandatory: boolean })[] = $state([]);

	let initialMetadata = $state.snapshot(metadata);
	for (let kind of elements.mandatory) {
		let index = initialMetadata.findIndex((m) => m.kind === kind.name);
		if (index >= 0) {
			let value = initialMetadata[index].value;
			initialMetadata.splice(index, 1);
			meta.push({ mandatory: true, kind: kind.name, value });
		} else {
			meta.push({ mandatory: true, kind: kind.name, value: "" });
		}
	}
	for (let { kind, value } of initialMetadata) {
		meta.push({ mandatory: false, kind, value });
	}

	export function addMetadata() {
		meta.push({ mandatory: false, kind: "", value: "" });
	}

	$effect(() => {
		metadata = meta.map(({ kind, value }) => ({ kind, value }));
	});
</script>

<h2 class="px-2 text-xl">Metadata</h2>

<div class="table">
	<div class="mb-2 table-row">
		<span class="table-cell px-2 py-2 font-bold">kind</span>
		<span class="table-cell py-2 pl-2 font-bold">value</span>
	</div>

	<div class="table-row">
		<span class="table-cell px-2 py-1">identifier</span>
		<span class="table-cell py-1 pl-2">
			{#if identifierEditable}
				<Input placeholder="doi" bind:value={identifier} />
			{:else}
				<a
					class="text-blue-700 hover:underline dark:text-blue-500"
					href={`https://doi.org/${identifier}`}
					target="_blank"
					rel="noopener noreferrer"
				>
					{identifier}
				</a>
			{/if}
		</span>
	</div>

	{#each meta as m, i (m)}
		{@const { type, choices } = flatMetadata.get(m.kind) || {}}

		<div class="table-row">
			{#if m.mandatory}
				<span class="table-cell px-2 py-1">{m.kind}</span>
			{:else}
				<select
					bind:value={m.kind}
					class="table-cell field-sizing-content rounded px-2 py-1 focus-within:outline-2 focus-within:outline-indigo-400"
				>
					<option class="hidden" value="" disabled selected>-- select a kind --</option>
					{#each Object.entries(elements) as [type, values]}
						<optgroup label={type}>
							{#each values as { name: kind }}
								<option value={kind}>{kind}</option>
							{/each}
						</optgroup>
					{/each}
				</select>
			{/if}

			<span class="table-cell py-1 pl-2">
				{#if m.kind === ""}
					<Input value="" disabled placeholder="please select a metadata kind" />
				{:else if type === "choice"}
					<select
						required
						class="picker-icon:text-black h-8 w-full cursor-pointer rounded border border-neutral-600 px-2 outline-black"
						bind:value={m.value}
					>
						<option class="hidden" value="" disabled selected>-- select an option --</option>
						{#each choices! as choice}
							<option class="cursor-pointer" value={choice}>{choice}</option>
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
				{:else}
					<Input placeholder={type || m.kind} bind:value={m.value} />
				{/if}
			</span>
			{#if !m.mandatory}
				<div class="relative table-cell">
					<button
						class="absolute top-2 left-[2px] inline h-6 cursor-pointer rounded-4xl focus-within:outline-indigo-400 focus:outline-2"
						onclick={() => meta.splice(i, 1)}
					>
						<X class="stroke-red-500 " />
					</button>
				</div>
			{/if}
		</div>
	{/each}
</div>
