<script lang="ts">
	import { type Metadata, type MetadataElements, type MetadataType } from "$lib/api";
	import YearInput from "./metadata-panel/YearInput.svelte";
	import DateInput from "./metadata-panel/DateInput.svelte";
	import OrcidInput from "./metadata-panel/OrcidInput.svelte";

	let {
		elements,
		identifier = $bindable(),
		metadata = $bindable(),
	}: {
		elements: MetadataElements;
		identifier: string;
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

	let meta = $state(
		elements.mandatory.map((kind) => ({
			mandatory: true,
			kind: kind.name,
			value: "",
		})),
	);

	export function addMetadata() {
		meta.push({ mandatory: false, kind: "", value: "" });
	}

	$effect(() => {
		metadata = meta.map(({ kind, value }) => ({ kind, value }));
	});
</script>

<h2 class="text-lg">Metadata</h2>

<div class="table border-collapse p-2">
	<div class="mb-2 table-row">
		<span class="table-cell p-2 pl-0 font-bold">kind</span>
		<span class="table-cell p-2 font-bold">value</span>
	</div>

	<div class="table-row">
		<span class="table-cell pr-2">identifier</span>
		<span class="table-cell pl-2">
			<input class="h-6 w-full outline" placeholder="doi" bind:value={identifier} />
		</span>
	</div>

	{#each meta as m, i (m)}
		{@const { type, choices } = flatMetadata.get(m.kind) || {}}

		<div class="table-row">
			{#if m.mandatory}
				<span class="table-cell pr-2">{m.kind}</span>
			{:else}
				<select bind:value={m.kind} class="table-cell field-sizing-content pr-2">
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

			<span class="table-cell pl-2">
				{#if m.kind === ""}
					<input class="h-6 w-full outline" disabled placeholder="please select a metadata kind" />
				{:else if type == "choice"}
					<select
						required
						class="picker-icon:text-black h-6 w-full outline outline-black"
						bind:value={m.value}
					>
						<option class="hidden" value="" disabled selected>-- select an option --</option>
						{#each choices! as choice}
							<option value={choice}>{choice}</option>
						{/each}
					</select>
				{:else if type === "date"}
					<DateInput bind:value={m.value} />
				{:else if type === "orcid"}
					<OrcidInput bind:value={m.value} />
				{:else if type == "url"}
					<input placeholder={type} class="h-6 w-full outline" bind:value={m.value} type="url" />
				{:else if type === "year"}
					<YearInput {type} bind:value={m.value} />
				{:else}
					<input placeholder={type || m.kind} class="h-6 w-full outline" bind:value={m.value} />
				{/if}
			</span>
			{#if !m.mandatory}
				<button
					onclick={() => meta.splice(i, 1)}
					class="cursor-pointer pl-2 font-bold text-red-500"
				>
					x
				</button>
			{:else}
				<button class="invisible pl-2 font-bold">x</button>
			{/if}
		</div>
	{/each}
</div>
