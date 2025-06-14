<script lang="ts">
	import type { PageProps } from "./$types";

	import GraphEdit from "$lib/GraphEdit.svelte";
	import { type Metadata } from "$lib/api";

	let { data }: PageProps = $props();

	let identifier = $state("");
	let metadata = $state<Metadata[]>([]);
</script>

<GraphEdit
	bind:identifier
	bind:metadata
	metadataElements={data.metadataElements}
	nodes={[]}
	edges={[]}
/>
