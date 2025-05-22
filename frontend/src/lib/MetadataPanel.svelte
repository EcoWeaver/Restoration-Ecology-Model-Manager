<script lang="ts">
	import type { Metadata, MetadataElements } from "$lib/api";

	let {
		elements,
		error,
		identifier = $bindable(),
		metadata = $bindable(),
	}: {
		elements: MetadataElements;
		error: boolean;
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
<div>Identifier: <input class="border" class:border-red-500={error} bind:value={identifier} /></div>

<table>
	<thead>
		<tr>
			<th>Kind</th>
			<th>Value</th>
		</tr>
	</thead>
	<tbody>
		{#each metadata as m}
			<tr>
				<td>{m.kind}</td>
				<td><input class="border" bind:value={m.value} /></td>
			</tr>
		{/each}
	</tbody>
</table>
