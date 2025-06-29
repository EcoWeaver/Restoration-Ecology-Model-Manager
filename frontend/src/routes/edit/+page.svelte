<script lang="ts">
	import type { PageProps } from "./$types";

	import GraphEdit from "$lib/GraphEdit.svelte";
	import { type Metadata } from "$lib/api";

	let { data }: PageProps = $props();
	let metadata = $state<Metadata[]>([]);
</script>

<svelte:head>
	<title>create new model</title>
</svelte:head>

<GraphEdit bind:metadata metadataElements={data.metadataElements} nodes={[]} edges={[]} />
