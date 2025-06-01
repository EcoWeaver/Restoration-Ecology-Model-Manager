<script lang="ts">
	import type { Metadata, MetadataElements, MetadataType } from "$lib/api";

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

	metadata = elements.mandatory.map((kind) => ({
		kind: kind.name,
		value: "",
	}));
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
			<input class="h-6 w-full outline" bind:value={identifier} />
		</span>
	</div>

	{#each metadata as m}
		{@const { type, choices } = flatMetadata.get(m.kind)!}

		<div class="table-row">
			<span class="table-cell pr-2">{m.kind}</span>
			<span class="table-cell pl-2">
				{#if type == "choice"}
					<select required class="h-6 w-full outline" bind:value={m.value}>
						{#each choices! as choice}
							<option value={choice}>{choice}</option>
						{/each}
					</select>
				{:else if type === "year"}
					<input
						placeholder={type}
						class="h-6 w-full outline"
						bind:value={() => parseInt(m.value), (v) => (m.value = v?.toString())}
						type="number"
					/>
				{:else if type == "url"}
					<input placeholder={type} class="h-6 w-full outline" bind:value={m.value} type="url" />
				{:else}
					<input placeholder={type || m.kind} class="h-6 w-full outline" bind:value={m.value} />
				{/if}
			</span>
		</div>
	{/each}
</div>
