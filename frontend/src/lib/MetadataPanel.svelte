<script lang="ts">
	import type { Metadata, MetadataElements } from '$lib/api';

	let { elements, metadata = $bindable() }: { elements: MetadataElements; metadata: Metadata[] } =
		$props();

	$effect(() => {
		metadata = elements.mandatory.map((kind) => ({
			kind,
			value: '',
			ontology: ''
		}));
	});
</script>

<h2 class="text-lg">Metadata</h2>
<table>
	<thead>
		<tr>
			<th>Kind</th>
			<th>Value</th>
			<th>Ontology</th>
		</tr>
	</thead>
	<tbody>
		{#each metadata as m}
			<tr>
				<td>{m.kind}</td>
				<td><input class="border" bind:value={m.value} /></td>
				<td><input class="border" bind:value={m.ontology} /></td>
			</tr>
		{/each}
	</tbody>
</table>
