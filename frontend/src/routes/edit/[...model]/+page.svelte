<script lang="ts">
	import { apiToFlow } from "$lib/model.js";
	import GraphEdit from "$lib/GraphEdit.svelte";

	let { data } = $props();
	const { model, version, metadataElements } = data;
	let { identifier, metadata } = $derived(model);

	const { nodes, edges } = apiToFlow(model.graph, true);
</script>

<svelte:head>
	<title>edit {model.identifier}@v{model.version}</title>
</svelte:head>

<GraphEdit
	{version}
	bind:identifier
	identifierEditable={false}
	bind:metadata
	{metadataElements}
	{nodes}
	{edges}
/>
