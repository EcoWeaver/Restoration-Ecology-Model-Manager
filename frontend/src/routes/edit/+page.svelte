<script lang="ts">
	import type { PageProps } from "./$types";
	import GraphEdit from "$lib/GraphEdit.svelte";
	import { type Metadata } from "$lib/api";
	import { SvelteFlowProvider } from "@xyflow/svelte";

	let { data }: PageProps = $props();
	const { metadataElements } = data;
	let metadata = $state<Metadata[]>([]);
</script>

<svelte:head>
	<title>create new model</title>
</svelte:head>

<SvelteFlowProvider>
	<GraphEdit bind:metadata {metadataElements} nodes={[]} edges={[]} />
</SvelteFlowProvider>
