<script lang="ts">
	import { apiToFlow } from "$lib/model.js";
	import GraphEdit from "$lib/GraphEdit.svelte";

	let { data } = $props();
	let { model, version, metadataElements } = data;
	let { identifier, metadata } = $derived(model);

	// svelte-ignore state_referenced_locally
	let { nodes, edges } = apiToFlow(model.graph, true);
</script>

<GraphEdit {version} bind:identifier bind:metadata {metadataElements} {nodes} {edges} />
