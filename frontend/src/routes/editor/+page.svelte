<script lang="ts">
	import { Background, BackgroundVariant, Panel, SvelteFlow } from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";
	import NodeComponent, { type Node } from "$lib/NodeComponent.svelte";
	import EdgeComponent, { type Edge } from "$lib/EdgeComponent.svelte";
	import MetadataPanel from "$lib/MetadataPanel.svelte";
	import {
		putModel,
		metadataElements,
		type Node as ApiNode,
		type Relation as ApiRelation,
		type Metadata,
		type GraphRequest,
	} from "$lib/api";

	let id = $state(0);

	let nodes: Node[] = $state.raw([]);
	let edges: Edge[] = $state.raw([]);

	let metadata = $state<Metadata[]>([]);

	function save() {
		const outNodes: ApiNode[] = nodes.map((n) => ({
			...n.data,
			position: {
				x: Math.round(n.position.x),
				y: Math.round(n.position.y),
			},
		}));
		const outRelations: ApiRelation[] = edges.map((e) => {
			const from = nodes.findIndex((n) => n.id === e.source);
			const to = nodes.findIndex((n) => n.id === e.target)!!;
			return {
				...e.data!!,
				from,
				to,
			};
		});

		const req: GraphRequest = {
			metadata,
			graph: { nodes: outNodes, relations: outRelations },
		};

		putModel(req);
	}

	function add() {
		const i = id++;
		const node: Node = {
			id: i.toString(),
			type: "custom",
			position: { x: 100, y: 100 },
			data: { body: "new", color: ["#ff0000", "#00ff00", "#0000ff"][i] },
		};
		nodes = [...nodes, node];
	}
</script>

<div class="h-dvh w-dvw">
	<SvelteFlow
		bind:nodes
		bind:edges
		proOptions={{ hideAttribution: true }}
		nodeTypes={{ custom: NodeComponent }}
		edgeTypes={{ custom: EdgeComponent }}
		onconnect={(c) => {
			const edge: Edge = {
				id: `${c.source}-${c.target}`,
				type: "custom",
				source: c.source,
				target: c.target,
				data: {
					kind: "correlated-with",
				},
			};
			edges = [...edges, edge];
		}}
	>
		<Panel position="top-left">
			<h1 class="text-2xl">Restoration Ecology Model Manager</h1>
			<button class="hover:bg-neutral-200" onclick={save}>Save</button>
			<button class="hover:bg-neutral-200" onclick={add}>Add</button>
		</Panel>
		<Panel position="top-right">
			{#await metadataElements then elements}
				<MetadataPanel {elements} bind:metadata />
			{/await}
		</Panel>
		<Background variant={BackgroundVariant.Dots} />
	</SvelteFlow>
</div>
