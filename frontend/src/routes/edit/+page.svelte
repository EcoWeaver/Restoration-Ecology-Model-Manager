<script lang="ts">
	import type { PageProps } from "./$types";

	import { goto } from "$app/navigation";
	import { type Node, type Edge, saveModel } from "$lib/model";
	import Graph from "$lib/Graph.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import {
		putModel,
		type Node as ApiNode,
		type Relation as ApiRelation,
		type Metadata,
		type ModelRequest,
	} from "$lib/api";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";

	let { data }: PageProps = $props();
	let { metadataElements } = data;

	let id = $state(0);

	let nodes: Node[] = $state.raw([]);
	let edges: Edge[] = $state.raw([]);

	let error = $state<string | undefined>(undefined);
	let identifier = $state("");
	let metadata = $state<Metadata[]>([]);

	async function save() {
		try {
			const model = await saveModel(identifier, metadata, nodes, edges);
			await goto(`/view/${model.identifier}?version=${model.version}`);
		} catch (e: any) {
			error = e;
		}
	}

	function add() {
		const i = id++;
		const node: Node = {
			id: i.toString(),
			type: "custom",
			position: { x: 100, y: 100 },
			data: { body: "new", color: ["#ff0000", "#00ff00", "#0000ff"][i % 3], editable: true },
		};
		nodes = [...nodes, node];
	}

	let metadataPanel: ReturnType<typeof MetadataPanel>;
</script>

{#if error}
	<div
		class="relative rounded border border-red-400 bg-red-100 px-4 py-3 text-red-700"
		role="alert"
	>
		<span class="block sm:inline">{error}</span>
		<button
			onclick={() => (error = undefined)}
			class="absolute top-0 right-0 bottom-0 cursor-pointer px-5 py-3 font-extrabold"
		>
			X
		</button>
	</div>
{/if}

<div class="h-dvh w-dvw">
	<Graph
		bind:nodes
		bind:edges
		interactive={true}
		onconnect={(c) => {
			const edge: Edge = {
				id: `${c.source}-${c.target}`,
				type: "custom",
				source: c.source,
				target: c.target,
				data: {
					kind: "correlated-with",
					editable: true,
				},
			};
			edges = [...edges, edge];
		}}
	>
		{#snippet info()}
			<h1 class="text-2xl">Restoration Ecology Model Manager</h1>
			<button class="hover:bg-neutral-200" onclick={save}>Save</button>
			<button class="hover:bg-neutral-200" onclick={add}>Add</button>
		{/snippet}

		{#snippet panel()}
			<PanelComponent onmetadata={() => metadataPanel.addMetadata()}>
				<MetadataPanel
					bind:this={metadataPanel}
					elements={metadataElements}
					bind:identifier
					bind:metadata
				/>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
