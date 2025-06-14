<script lang="ts">
	import Graph from "$lib/Graph.svelte";
	import { saveModel } from "$lib/model.js";
	import type { Edge, Node } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import { goto } from "$app/navigation";
	import type { Metadata, MetadataElements } from "./api";

	let {
		version = undefined,
		identifier = $bindable(),
		metadata = $bindable(),
		metadataElements,
		nodes: initialNodes,
		edges: initialEdges,
	}: {
		version?: string;
		identifier: string;
		metadata: Metadata[];
		metadataElements: MetadataElements;
		nodes: Node[];
		edges: Edge[];
	} = $props();

	let nodes: Node[] = $state.raw(initialNodes);
	let edges: Edge[] = $state.raw(initialEdges);

	let error = $state<string | undefined>(undefined);
	async function save() {
		try {
			const model = await saveModel(identifier, metadata, nodes, edges, version);
			await goto(`/view/${model.identifier}?version=${model.version}`);
		} catch (e: any) {
			error = e;
		}
	}

	let id = $state(0);
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
		{#snippet headerLeft()}
			<button class="hover:bg-neutral-200 dark:hover:bg-neutral-700" onclick={add}>Add</button>
		{/snippet}

		{#snippet headerRight()}
			<button class="hover:bg-neutral-200 dark:hover:bg-neutral-700" onclick={save}>Save</button>
		{/snippet}

		{#snippet panel()}
			<PanelComponent onmetadata={() => metadataPanel.addMetadata()}>
				<MetadataPanel
					bind:this={metadataPanel}
					elements={metadataElements}
					{identifier}
					identifierEditable={false}
					bind:metadata
				/>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
