<script lang="ts">
	import { apiToFlow } from "$lib/model.js";
	import GraphEdit from "$lib/GraphEdit.svelte";
	import { SvelteFlowProvider } from "@xyflow/svelte";

	let { data } = $props();
	const { model, metadataElements } = data;
	let { uuid, metadata } = $derived(model);

	const { nodes, edges } = apiToFlow(model.graph, true);
</script>

<svelte:head>
	<title>edit {metadata.find(({ kind }) => kind === "title")!.value}</title>
</svelte:head>

<SvelteFlowProvider>
	<GraphEdit parent={uuid} bind:metadata {metadataElements} {nodes} {edges} />
</SvelteFlowProvider>
