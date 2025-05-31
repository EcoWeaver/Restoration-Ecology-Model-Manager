<script lang="ts">
	import type { Metadata, MetadataElements } from "$lib/api";

	let {
		elements,
		identifier = $bindable(),
		metadata = $bindable(),
	}: {
		elements: MetadataElements;
		identifier: string;
		metadata: Metadata[];
	} = $props();

	$effect(() => {
		metadata = elements.mandatory.map((kind) => ({
			kind: kind.name,
			value: "",
		}));
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
		<span class="table-cell pl-2"><input class="outline" bind:value={identifier} /></span>
	</div>

	{#each metadata as m}
		<div class="table-row">
			<span class="table-cell pr-2">{m.kind}</span>
			<span class="table-cell pl-2"><input class="outline" bind:value={m.value} /></span>
		</div>
	{/each}
</div>
